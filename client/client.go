package client

import "github.com/ClifHouck/unified/types"

import "github.com/coder/websocket"
import log "github.com/sirupsen/logrus"

import "context"
import "errors"
import "fmt"
import "io"
import "encoding/json"
import "net/http"
import "crypto/tls"
import "time"

const URL_TEMPLATE string = "%s://%s/proxy/%s/integration/v1/%s"

type ApiEndpoint struct {
	UrlFragment string
	Method      string
	Description string
	Application string
	Protocol    string
}

var API = map[string]*ApiEndpoint{
	"protect/meta/info": &ApiEndpoint{
		UrlFragment: "meta/info",
		Method:      http.MethodGet,
		Description: "Get application information",
		Application: "protect",
	},
	"protect/subscribe/events": &ApiEndpoint{
		UrlFragment: "subscribe/events",
		Method:      http.MethodGet,
		Description: "Get Protect event messages",
		Application: "protect",
		Protocol:    "wss",
	},
	"protect/subscribe/devices": &ApiEndpoint{
		UrlFragment: "subscribe/devices",
		Method:      http.MethodGet,
		Description: "Get Protect device updates",
		Application: "protect",
		Protocol:    "wss",
	},
	"protect/cameras": &ApiEndpoint{
		UrlFragment: "cameras",
		Method:      http.MethodGet,
		Description: "Get all cameras",
		Application: "protect",
	},
	"protect/camera/id": &ApiEndpoint{
		UrlFragment: "cameras/%s",
		Method:      http.MethodGet,
		Description: "Get camera details",
		Application: "protect",
	},
}

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
	config *Config
	client *http.Client
}

// FIXME: Client-wide context
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.InsecureSkipVerify,
				},
			},
		},
	}
}

func (c *Client) headers() *http.Header {
	headers := &http.Header{}
	headers.Add("X-API-KEY", c.config.ApiKey)
	headers.Add("Accept", "application/json")
	return headers
}

func (c *Client) webSocketHeaders() *http.Header {
	headers := &http.Header{}
	headers.Add("X-API-KEY", c.config.ApiKey)
	return headers
}

func (c *Client) renderUrl(endpoint *ApiEndpoint, urlArgs []any) string {
	renderedFragment := fmt.Sprintf(endpoint.UrlFragment, urlArgs...)

	protocol := "https"
	if endpoint.Protocol != "" {
		protocol = endpoint.Protocol
	}

	return fmt.Sprintf(URL_TEMPLATE, protocol, c.config.Hostname, endpoint.Application, renderedFragment)
}

func (c *Client) doRequest(endpoint *ApiEndpoint, expectedStatus int) ([]byte, error) {
	return c.doRequestArgs(endpoint, expectedStatus, []any{})
}

func (c *Client) doRequestArgs(endpoint *ApiEndpoint, expectedStatus int, urlArgs []any) ([]byte, error) {
	renderedUrl := c.renderUrl(endpoint, urlArgs)

	req, err := http.NewRequest(endpoint.Method, renderedUrl, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header = *c.headers()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode != expectedStatus {
		return nil, errors.New(fmt.Sprintf("Got bad status %d when requesting '%s'", resp.StatusCode, renderedUrl))
	}

	log.WithFields(log.Fields{
		"url":    renderedUrl,
		"status": resp.StatusCode,
	}).Info("URL request success")

	return body, nil
}

func (c *Client) ProtectInfo() (*types.ProtectInfo, error) {
	body, err := c.doRequest(API["protect/meta/info"], http.StatusOK)
	if err != nil {
		return nil, err
	}

	var info types.ProtectInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *Client) GetCameras() ([]*types.Camera, error) {
	body, err := c.doRequest(API["protect/cameras"], http.StatusOK)
	if err != nil {
		return nil, err
	}

	var cameras []*types.Camera

	err = json.Unmarshal(body, &cameras)
	if err != nil {
		return nil, err
	}

	return cameras, nil
}

func (c *Client) GetCameraDetails(cameraID string) (*types.Camera, error) {
	body, err := c.doRequestArgs(API["protect/camera/id"], http.StatusOK, []any{cameraID})
	if err != nil {
		return nil, err
	}

	var camera *types.Camera

	err = json.Unmarshal(body, &camera)
	if err != nil {
		return nil, err
	}

	return camera, nil
}

type WebSocketMessage struct {
	MessageType websocket.MessageType
	Error       error
}

type ProtectEventMessage struct {
	WebSocketMessage
	Event types.ProtectEvent
}

// Periodically pings the websocket connection to keep it alive.
// coder/websocket is concurrency-safe for writes so this may be used with
// any websocket connection.
func (c *Client) webSocketKeepAlive(ctx context.Context, conn *websocket.Conn, url string) {
	tickChan := time.Tick(c.config.WebSocketKeepAliveInterval)
	for next := range tickChan {
		err := conn.Ping(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"url":   url,
				"error": err.Error(),
			}).Error("WebSocket.Ping returned error")
			return
		}
		log.WithFields(log.Fields{
			"url":          url,
			"next_ping_at": next,
		}).Debug("WebSocket.Ping (Keep-Alive) Success")
	}
}

func (c *Client) SubscribeProtectEvents(ctx context.Context) (<-chan *ProtectEventMessage, error) {
	url := c.renderUrl(API["protect/subscribe/events"], []any{})
	conn, _, err := websocket.Dial(ctx,
		url,
		&websocket.DialOptions{
			HTTPClient: c.client,
			HTTPHeader: *c.webSocketHeaders()})
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"url": url,
	}).Info("WebSocket.Dial() success")

	go c.webSocketKeepAlive(ctx, conn, url)

	eventChan := make(chan *ProtectEventMessage)

	go func() {
		for {
			messageType, data, err := conn.Read(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"url":   url,
					"error": err.Error(),
				}).Error("WebSocket Read returned error")
				close(eventChan)
				return
			}

			var protectEvent types.ProtectEvent
			err = json.Unmarshal(data, &protectEvent)
			if err != nil {
				log.WithFields(log.Fields{
					"url":   url,
					"error": err.Error(),
				}).Error("json.Unmarshal returned error")
				close(eventChan)
				return
			}

			eventChan <- &ProtectEventMessage{
				WebSocketMessage: WebSocketMessage{
					MessageType: messageType,
					Error:       err,
				},
				Event: protectEvent,
			}
		}
	}()

	return eventChan, nil
}

type ProtectDeviceMessage struct {
	WebSocketMessage
	Event types.ProtectDeviceUpdate
}

func (c *Client) SubscribeProtectDeviceUpdates(ctx context.Context) (<-chan *ProtectDeviceMessage, error) {
	url := c.renderUrl(API["protect/subscribe/devices"], []any{})
	conn, _, err := websocket.Dial(ctx,
		url,
		&websocket.DialOptions{
			HTTPClient: c.client,
			HTTPHeader: *c.webSocketHeaders()})
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"url": url,
	}).Info("WebSocket Dial Success")

	go c.webSocketKeepAlive(ctx, conn, url)

	eventChan := make(chan *ProtectDeviceMessage)

	go func() {
		for {
			messageType, data, err := conn.Read(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"url":   url,
					"error": err.Error(),
				}).Error("json.Unmarshal returned error")
				close(eventChan)
				return
			}

			var protectDeviceUpdate types.ProtectDeviceUpdate
			err = json.Unmarshal(data, &protectDeviceUpdate)
			if err != nil {
				log.WithFields(log.Fields{
					"url":   url,
					"error": err.Error(),
				}).Error("json.Unmarshal returned error")
				close(eventChan)
				return
			}

			eventChan <- &ProtectDeviceMessage{
				WebSocketMessage: WebSocketMessage{
					MessageType: messageType,
					Error:       err,
				},
				Event: protectDeviceUpdate,
			}
		}
	}()

	return eventChan, nil
}
