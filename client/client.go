package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/sirupsen/logrus"

	"github.com/ClifHouck/unified/types"
)

type Config struct {
	// The hostname of the unifi control plane.
	Hostname string
	// API key issued by unifi control plane. Must be included in requests
	// for authorization.
	APIKey string
	// Controls the interval between keep-alive pings for websocket
	// connections.
	WebSocketKeepAliveInterval time.Duration
	// Controls the configuration of http.Client TLS verification behavior.
	InsecureSkipVerify bool
}

func NewDefaultConfig(apiKey string) *Config {
	return &Config{
		Hostname:                   "unifi",
		APIKey:                     apiKey,
		WebSocketKeepAliveInterval: time.Second * 30,
		// Unfortunately, unifi doesn't seem to self-sign for `unifi`, nor
		// `192.168.1.1` for that matter.
		InsecureSkipVerify: true,
	}
}

// IsValid returns true if config is valid, and false otherwise. Also returns a list of
// reasons verification failed.
func (c *Config) IsValid() (bool, []string) {
	reasons := []string{}

	if c.APIKey == "" {
		reasons = append(reasons, "APIKey must not be empty")
	}

	if c.WebSocketKeepAliveInterval < time.Second {
		reasons = append(
			reasons,
			"WebSocketKeepAliveInterval is too short. Must be longer than one second.",
		)
	}

	if c.WebSocketKeepAliveInterval > time.Minute*10 {
		reasons = append(
			reasons,
			"WebSocketKeepAliveInterval is too long. Must be shorter than ten minutes.",
		)
	}

	valid := len(reasons) == 0
	return valid, reasons
}

type Client struct {
	ctx    context.Context
	config *Config
	client *http.Client

	log *logrus.Logger

	Network types.NetworkV1
	Protect types.ProtectV1
}

func NewClient(ctx context.Context, config *Config, log *logrus.Logger) *Client {
	client := &Client{
		ctx:    ctx,
		config: config,
		log:    log,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.InsecureSkipVerify, //nolint:gosec // TODO: Figure out how to always enable TLS verification!
				},
			},
		},
	}
	client.Network = &networkV1Client{client: client}
	client.Protect = &protectV1Client{client: client}
	return client
}

func (c *Client) headers(contentType string) *http.Header {
	headers := &http.Header{}
	headers.Add("X-Api-Key", c.config.APIKey)
	headers.Add("Accept", "application/json")
	if contentType != "" {
		headers.Add("Content-Type", contentType)
	} else {
		headers.Add("Content-Type", "application/json")
	}
	return headers
}

func (c *Client) webSocketHeaders() *http.Header {
	headers := &http.Header{}
	headers.Add("X-Api-Key", c.config.APIKey)
	return headers
}

type apiEndpoint struct {
	Application    string
	ContentType    string
	Description    string
	ExpectedStatus int
	HasRequestBody bool
	Method         string
	NumQueryArgs   int
	NumURLArgs     int
	URLFragment    string
}

type requestArgs struct {
	Endpoint     *apiEndpoint
	URLArguments []any
	RequestBody  io.Reader
	Query        *url.Values
}

const urlTemplate string = "%s://%s/proxy/%s/integration/v1/%s"

func (c *Client) renderURL(req *requestArgs) string {
	renderedFragment := req.Endpoint.URLFragment

	if req.Endpoint.NumURLArgs != len(req.URLArguments) {
		c.log.WithFields(logrus.Fields{
			"expected_args": req.Endpoint.NumURLArgs,
			"actual_args":   len(req.URLArguments),
			"URLFragment":   req.Endpoint.URLFragment,
		}).Fatal("Number of url arguments does not match number of arguments " +
			"required by the API endpoint")
	}

	if len(req.URLArguments) > 0 {
		renderedFragment = fmt.Sprintf(req.Endpoint.URLFragment, req.URLArguments...)
	}

	url := fmt.Sprintf(
		urlTemplate,
		"https",
		c.config.Hostname,
		req.Endpoint.Application,
		renderedFragment,
	)

	// TODO: Some sort of sanity checking on number & keys of query args...
	if req.Query != nil {
		encodedQuery := req.Query.Encode()
		if len(encodedQuery) > 0 {
			url = url + "?" + encodedQuery
		}
	}

	c.log.WithFields(logrus.Fields{
		"url": url,
	}).Trace("Rendered url")
	return url
}

func (c *Client) decodeErrorResponse(body []byte) {
	var unifiError types.Error
	c.log.Trace(string(body))
	err := json.Unmarshal(body, &unifiError)

	switch {
	case err == nil && unifiError.StatusCode != 0:
		c.log.WithFields(logrus.Fields{
			"code":    unifiError.StatusCode,
			"name":    unifiError.StatusName,
			"message": unifiError.Message,
		}).Error("UniFi application returned an error")
	case err == nil && unifiError.StatusCode == 0:
		// This is probably an undocumented Protect Error
		var protectError types.ProtectErrorMessage
		protectErr := json.Unmarshal(body, &protectError)
		if protectErr != nil {
			c.log.Error("Could not unwrap error message as ProtectErrorMessage")
			break
		}
		c.log.WithFields(logrus.Fields{"code": protectError.Error,
			"name":   protectError.Name,
			"entity": protectError.Entity,
		}).Error("UniFi application returned a protect error")

		for _, issue := range protectError.Issues {
			c.log.WithFields(logrus.Fields{
				"instance_path": issue.InstancePath,
				"keyword":       issue.Keyword,
			}).Errorf("Issue with Request: %s", issue.Message)
		}
	default:
		c.log.Debug(string(body))
		c.log.Errorf("Could not decode UniFi error despite bad response code: %s", err.Error())
	}
}

func (c *Client) doRequest(req *requestArgs) ([]byte, error) {
	renderedURL := c.renderURL(req)

	if req.Endpoint.HasRequestBody && req.RequestBody == http.NoBody {
		c.log.WithFields(logrus.Fields{
			"URLFragment": req.Endpoint.URLFragment,
		}).Fatal("Request should have a body but http.NoBody was passed")
	}

	request, err := http.NewRequestWithContext(
		c.ctx,
		req.Endpoint.Method,
		renderedURL,
		req.RequestBody,
	)
	if err != nil {
		return nil, err
	}

	request.Header = *c.headers(req.Endpoint.ContentType)

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			c.log.WithFields(logrus.Fields{
				"url":    renderedURL,
				"status": resp.StatusCode,
			}).Errorf("Error closing response body: %s", closeErr.Error())
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	expectedStatus := req.Endpoint.ExpectedStatus
	if expectedStatus == 0 {
		expectedStatus = http.StatusOK
	}

	if resp.StatusCode != expectedStatus {
		c.decodeErrorResponse(body)

		return nil, fmt.Errorf(
			"got unexpected http code %d when requesting '%s'",
			resp.StatusCode,
			renderedURL,
		)
	}

	c.log.WithFields(logrus.Fields{
		"url":    renderedURL,
		"status": resp.StatusCode,
	}).Debug("https request success")

	return body, nil
}

// Periodically pings the websocket connection to keep it alive.
// coder/websocket is concurrency-safe for writes so this may be used with
// any websocket connection.
func (c *Client) webSocketKeepAlive(conn *websocket.Conn, url string) {
	tickChan := time.Tick(c.config.WebSocketKeepAliveInterval)
	for next := range tickChan {
		err := conn.Ping(c.ctx)
		if err != nil {
			c.log.WithFields(logrus.Fields{
				"url":   url,
				"error": err.Error(),
			}).Error("WebSocket.Ping returned error")
			return
		}
		c.log.WithFields(logrus.Fields{
			"url":          url,
			"next_ping_at": next,
		}).Trace("WebSocket.Ping (Keep-Alive) Success")
	}
}

func buildQuery(filter types.Filter, pageArgs *types.PageArguments) *url.Values {
	query := &url.Values{}
	filterStr := string(filter)
	if len(filterStr) > 0 {
		query.Add("filter", string(filter))
	}
	if pageArgs != nil {
		if pageArgs.Limit != 0 {
			query.Add("limit", strconv.FormatUint(uint64(pageArgs.Limit), 10))
		}
		if pageArgs.Offset != 0 {
			query.Add("offset", strconv.FormatUint(uint64(pageArgs.Offset), 10))
		}
	}
	return query
}
