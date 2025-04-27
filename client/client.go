package client

import "github.com/ClifHouck/unified/types"

import "github.com/coder/websocket"
import log "github.com/sirupsen/logrus"

import "bytes"
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

// TODO: Maybe move this to an api module?
var API = map[string]*ApiEndpoint{
	"network/info": &ApiEndpoint{
		UrlFragment: "info",
		Method:      http.MethodGet,
		Description: "Get application information",
		Application: "network",
	},
	"network/sites/list": &ApiEndpoint{
		UrlFragment: "sites",
		Method:      http.MethodGet,
		Description: "List local sites managed by this Network application",
		Application: "network",
	},
	"network/devices/list": &ApiEndpoint{
		UrlFragment: "sites/%s/devices",
		Method:      http.MethodGet,
		Description: "List adopted devices of a site",
		Application: "network",
	},
	"network/devices/id": &ApiEndpoint{
		UrlFragment: "sites/%s/devices/%s",
		Method:      http.MethodGet,
		Description: "Get device details",
		Application: "network",
	},
	"network/devices/id/stats": &ApiEndpoint{
		UrlFragment: "sites/%s/devices/%s/statistics/latest",
		Method:      http.MethodGet,
		Description: "Get latest device statistics",
		Application: "network",
	},
	"network/devices/id/actions": &ApiEndpoint{
		UrlFragment: "sites/%s/devices/%s/actions",
		Method:      http.MethodPost,
		Description: "Execute an action on a device",
		Application: "network",
	},
	"network/clients/list": &ApiEndpoint{
		UrlFragment: "sites/%s/clients",
		Method:      http.MethodGet,
		Description: "List clients of a site",
		Application: "network",
	},
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
	headers.Add("Content-Type", "application/json")
	return headers
}

func (c *Client) webSocketHeaders() *http.Header {
	headers := &http.Header{}
	headers.Add("X-API-KEY", c.config.ApiKey)
	return headers
}

func (c *Client) renderUrl(endpoint *ApiEndpoint, urlArgs []any) string {
	renderedFragment := endpoint.UrlFragment
	if len(urlArgs) > 0 {
		renderedFragment = fmt.Sprintf(endpoint.UrlFragment, urlArgs...)
	}

	protocol := "https"
	if endpoint.Protocol != "" {
		protocol = endpoint.Protocol
	}

	url := fmt.Sprintf(URL_TEMPLATE, protocol, c.config.Hostname, endpoint.Application, renderedFragment)
	log.WithFields(log.Fields{
		"url": url,
	}).Trace("Rendered url")
	return url
}

func (c *Client) doRequest(endpoint *ApiEndpoint, expectedStatus int) ([]byte, error) {
	return c.doRequestArgs(endpoint, expectedStatus, []any{})
}

func (c *Client) doRequestArgs(endpoint *ApiEndpoint, expectedStatus int, urlArgs []any) ([]byte, error) {
	return c.doRequestArgsAndBody(endpoint, expectedStatus, urlArgs, http.NoBody)
}

func (c *Client) doRequestArgsAndBody(endpoint *ApiEndpoint, expectedStatus int, urlArgs []any, requestBody io.Reader) ([]byte, error) {
	renderedUrl := c.renderUrl(endpoint, urlArgs)

	req, err := http.NewRequest(endpoint.Method, renderedUrl, requestBody)
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
		return nil, errors.New(fmt.Sprintf("Got unexpected http code %d when requesting '%s'", resp.StatusCode, renderedUrl))
	}

	log.WithFields(log.Fields{
		"url":    renderedUrl,
		"status": resp.StatusCode,
	}).Debug("URL request success")

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

type ProtectDeviceEventMessage struct {
	WebSocketMessage
	Event types.ProtectDeviceEvent
}

func (c *Client) SubscribeProtectDeviceUpdates(ctx context.Context) (<-chan *ProtectDeviceEventMessage, error) {
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

	eventChan := make(chan *ProtectDeviceEventMessage)

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

			var protectDeviceUpdate types.ProtectDeviceEvent
			err = json.Unmarshal(data, &protectDeviceUpdate)
			if err != nil {
				log.WithFields(log.Fields{
					"url":   url,
					"error": err.Error(),
				}).Error("json.Unmarshal returned error")
				close(eventChan)
				return
			}

			eventChan <- &ProtectDeviceEventMessage{
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

func (c *Client) NetworkInfo() (*types.NetworkInfo, error) {
	body, err := c.doRequest(API["network/info"], http.StatusOK)
	if err != nil {
		return nil, err
	}

	var info types.NetworkInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *Client) ListAllDevices(siteID string) ([]*types.DeviceListEntry, error) {
	// FIXME: We have to send url query args
	body, err := c.doRequestArgs(API["network/devices/list"], http.StatusOK, []any{siteID})
	if err != nil {
		return nil, err
	}

	var deviceListPage *types.DeviceListPage

	err = json.Unmarshal(body, &deviceListPage)
	if err != nil {
		return nil, err
	}

	// FIXME: Deal with pagination!

	return deviceListPage.Data, nil

}

func (c *Client) GetDeviceDetails(siteID string, deviceID string) (*types.Device, error) {
	body, err := c.doRequestArgs(API["network/devices/id"], http.StatusOK, []any{siteID, deviceID})
	if err != nil {
		return nil, err
	}

	var device *types.Device

	err = json.Unmarshal(body, &device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (c *Client) GetDeviceStats(siteID string, deviceID string) (*types.DeviceStatistics, error) {
	body, err := c.doRequestArgs(API["network/devices/id/stats"], http.StatusOK, []any{siteID, deviceID})
	if err != nil {
		return nil, err
	}

	var stats *types.DeviceStatistics

	err = json.Unmarshal(body, &stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (c *Client) ExecuteDeviceAction(siteID string, deviceID string, action *types.DeviceActionRequest) error {
	jsonBody, err := json.Marshal(action)
	log.WithFields(log.Fields{
		"method": "ExecuteDeviceAction",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(jsonBody)
	_, err = c.doRequestArgsAndBody(API["network/devices/id/actions"], http.StatusOK,
		[]any{siteID, deviceID},
		bodyReader)

	return err
}

func (c *Client) ListAllSites() ([]*types.Site, error) {
	body, err := c.doRequest(API["network/sites/list"], http.StatusOK)
	if err != nil {
		return nil, err
	}

	var siteListPage *types.SiteListPage

	err = json.Unmarshal(body, &siteListPage)
	if err != nil {
		return nil, err
	}

	// FIXME: Deal with pagination!

	return siteListPage.Data, nil
}

func (c *Client) ListAllClients(siteID string) ([]*types.Client, error) {
	body, err := c.doRequestArgs(API["network/clients/list"], http.StatusOK, []any{siteID})
	if err != nil {
		return nil, err
	}

	var clientListPage *types.ClientListPage

	err = json.Unmarshal(body, &clientListPage)
	if err != nil {
		return nil, err
	}

	// FIXME: Deal with pagination!

	return clientListPage.Data, nil
}
