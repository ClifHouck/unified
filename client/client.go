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

const URL_TEMPLATE string = "%s://%s/proxy/%s/integration/v1/%s"

type Config struct {
	// The hostname of the unifi control plane.
	Hostname string
	// API key issued by unifi control plane. Must be included in requests
	// for authorization.
	ApiKey string
	// Controls the interval between keep-alive pings for websocket
	// connections.
	WebSocketKeepAliveInterval time.Duration
	// Controls the configuration of http.Client TLS verification behavior.
	InsecureSkipVerify bool
}

func NewDefaultConfig(apiKey string) *Config {
	return &Config{
		Hostname:                   "unifi",
		ApiKey:                     apiKey,
		WebSocketKeepAliveInterval: time.Second * 30,
		// Unfortunately, unifi doesn't seem to self-sign for `unifi`, nor
		// `192.168.1.1` for that matter.
		InsecureSkipVerify: true,
	}
}

// Returns true if config is valid, false otherwise along with a list of
// reasons verification failed.
func (c *Config) IsValid() (valid bool, reasons []string) {
	if c.ApiKey == "" {
		reasons = append(reasons, "ApiKey must not be empty")
	}

	if c.WebSocketKeepAliveInterval < time.Second {
		reasons = append(reasons, "WebSocketKeepAliveInterval is too short. Must be longer than one second.")
	}

	if c.WebSocketKeepAliveInterval > time.Minute*10 {
		reasons = append(reasons, "WebSocketKeepAliveInterval is too long. Must be shorter than ten minutes.")
	}

	valid = len(reasons) == 0
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
					InsecureSkipVerify: config.InsecureSkipVerify, //nolint:gosec,G402 // TODO: Figure out how to always enable TLS verification!
				},
			},
		},
	}
	client.Network = &networkV1Client{client: client}
	client.Protect = &protectV1Client{client: client}
	return client
}

func (c *Client) headers() *http.Header {
	headers := &http.Header{}
	headers.Add("X-Api-Key", c.config.ApiKey)
	headers.Add("Accept", "application/json")
	headers.Add("Content-Type", "application/json")
	return headers
}

func (c *Client) webSocketHeaders() *http.Header {
	headers := &http.Header{}
	headers.Add("X-Api-Key", c.config.ApiKey)
	return headers
}

type apiEndpoint struct {
	UrlFragment    string
	Method         string
	ExpectedStatus int
	Description    string
	Application    string
	Protocol       string
	NumUrlArgs     int
	NumQueryArgs   int
	HasRequestBody bool
}

type requestArgs struct {
	Endpoint     *apiEndpoint
	UrlArguments []any
	RequestBody  io.Reader
	Query        *url.Values
}

func (c *Client) renderUrl(req *requestArgs) string {
	renderedFragment := req.Endpoint.UrlFragment

	if req.Endpoint.NumUrlArgs != len(req.UrlArguments) {
		c.log.WithFields(logrus.Fields{
			"expected_args": req.Endpoint.NumUrlArgs,
			"actual_args":   len(req.UrlArguments),
			"urlFragment":   req.Endpoint.UrlFragment,
		}).Fatal("Number of url arguments does not match number of arguments " +
			"required by the API endpoint")
	}

	if len(req.UrlArguments) > 0 {
		renderedFragment = fmt.Sprintf(req.Endpoint.UrlFragment, req.UrlArguments...)
	}

	protocol := "https"
	if req.Endpoint.Protocol != "" {
		protocol = req.Endpoint.Protocol
	}

	url := fmt.Sprintf(URL_TEMPLATE, protocol, c.config.Hostname, req.Endpoint.Application, renderedFragment)

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

func (c *Client) doRequest(req *requestArgs) ([]byte, error) {
	renderedUrl := c.renderUrl(req)

	if req.Endpoint.HasRequestBody && req.RequestBody == http.NoBody {
		c.log.WithFields(logrus.Fields{
			"urlFragment": req.Endpoint.UrlFragment,
		}).Fatal("Request should have a body but http.NoBody was passed")
	}

	request, err := http.NewRequestWithContext(c.ctx, req.Endpoint.Method, renderedUrl, req.RequestBody)
	if err != nil {
		return nil, err
	}

	request.Header = *c.headers()

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			c.log.WithFields(logrus.Fields{
				"url":    renderedUrl,
				"status": resp.StatusCode,
			}).Errorf("Error closing response body: %s", err.Error())
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
		var unifiError types.Error
		err = json.Unmarshal(body, &unifiError)
		if err == nil && unifiError.StatusCode != 0 {
			c.log.WithFields(logrus.Fields{
				"code":    unifiError.StatusCode,
				"name":    unifiError.StatusName,
				"message": unifiError.Message,
			}).Error("UniFi application returned an error")
		}

		return nil, fmt.Errorf("got unexpected http code %d when requesting '%s'", resp.StatusCode, renderedUrl)
	}

	c.log.WithFields(logrus.Fields{
		"url":    renderedUrl,
		"status": resp.StatusCode,
	}).Debug("URL request success")

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
