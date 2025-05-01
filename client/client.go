package client

import "bytes"
import "context"
import "crypto/tls"
import "encoding/json"
import "fmt"
import "io"
import "net/http"
import "net/url"
import "time"

import "github.com/coder/websocket"
import "github.com/sirupsen/logrus"

import "github.com/ClifHouck/unified/types"

const URL_TEMPLATE string = "%s://%s/proxy/%s/integration/v1/%s"

type ApiEndpoint struct {
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

// TODO: Maybe move this to an api module?
var networkAPI = map[string]*ApiEndpoint{
	// Application related
	"Info": &ApiEndpoint{
		UrlFragment: "info",
		Method:      http.MethodGet,
		Description: "Get application information",
		Application: "network",
	},

	// Site related
	"Sites": &ApiEndpoint{
		UrlFragment: "sites",
		Method:      http.MethodGet,
		Description: "List local sites managed by this Network application",
		Application: "network",
	},

	/// Client related
	"Clients": &ApiEndpoint{
		UrlFragment:  "sites/%s/clients",
		Method:       http.MethodGet,
		Description:  "List clients of a site",
		Application:  "network",
		NumUrlArgs:   1,
		NumQueryArgs: 1,
	},
	"ClientDetails": &ApiEndpoint{
		UrlFragment: "sites/%s/clients/%s",
		Method:      http.MethodGet,
		Description: "Get client details",
		Application: "network",
		NumUrlArgs:  2,
	},
	"ClientExecuteAction": &ApiEndpoint{
		UrlFragment: "sites/%s/devices/%s/actions",
		Method:      http.MethodPost,
		Description: "Execute an action on a device",
		Application: "network",
		NumUrlArgs:  2,
	},

	// Devices related
	"Devices": &ApiEndpoint{
		UrlFragment: "sites/%s/devices",
		Method:      http.MethodGet,
		Description: "List adopted devices of a site",
		Application: "network",
		NumUrlArgs:  1,
	},
	"DeviceDetails": &ApiEndpoint{
		UrlFragment: "sites/%s/devices/%s",
		Method:      http.MethodGet,
		Description: "Get device details",
		Application: "network",
		NumUrlArgs:  2,
	},
	"DeviceStatistics": &ApiEndpoint{
		UrlFragment: "sites/%s/devices/%s/statistics/latest",
		Method:      http.MethodGet,
		Description: "Get latest device statistics",
		Application: "network",
		NumUrlArgs:  2,
	},
	"DeviceExecuteAction": &ApiEndpoint{
		UrlFragment:    "sites/%s/devices/%s/actions",
		Method:         http.MethodPost,
		Description:    "Execute an action on a device",
		Application:    "network",
		NumUrlArgs:     2,
		HasRequestBody: true,
	},
	"DevicePortExecuteAction": &ApiEndpoint{
		UrlFragment:    "sites/%s/devices/%s/actions",
		Method:         http.MethodPost,
		Description:    "Execute an action on a device",
		Application:    "network",
		NumUrlArgs:     3,
		HasRequestBody: true,
	},

	// Voucher related
	"Vouchers": &ApiEndpoint{
		UrlFragment:  "sites/%s/hotspot/vouchers",
		Method:       http.MethodGet,
		Description:  "List vouchers of a site",
		Application:  "network",
		NumUrlArgs:   1,
		NumQueryArgs: 3,
	},
	"VoucherDetails": &ApiEndpoint{
		UrlFragment: "sites/%s/vouchers/%s",
		Method:      http.MethodGet,
		Description: "Get voucher details",
		Application: "network",
		NumUrlArgs:  2,
	},
	"VoucherGenerate": &ApiEndpoint{
		UrlFragment:    "sites/%s/hotspot/vouchers",
		Method:         http.MethodPost,
		ExpectedStatus: http.StatusCreated,
		Description:    "Generate vouchers",
		Application:    "network",
		NumUrlArgs:     1,
		HasRequestBody: true,
	},
	"VoucherDelete": &ApiEndpoint{
		UrlFragment: "sites/%s/hotspot/vouchers/%s",
		Method:      http.MethodDelete,
		Description: "Delete vouchers",
		Application: "network",
		NumUrlArgs:  2,
	},
	"VoucherDeleteByFilters": &ApiEndpoint{
		UrlFragment:  "sites/%s/hotspot/vouchers",
		Method:       http.MethodDelete,
		Description:  "Delete vouchers by filter",
		Application:  "network",
		NumUrlArgs:   1,
		NumQueryArgs: 1,
	},
}

var protectAPI = map[string]*ApiEndpoint{
	"Info": &ApiEndpoint{
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
		NumUrlArgs:  1,
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
	ctx    context.Context
	config *Config
	client *http.Client

	log *logrus.Logger

	Network types.NetworkV1
	Protect types.ProtectV1
}

type networkV1Client struct {
	client *Client
}

type protectV1Client struct {
	client *Client
}

func NewClient(ctx context.Context, config *Config, log *logrus.Logger) *Client {
	client := &Client{
		ctx:    ctx,
		config: config,
		log:    log,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.InsecureSkipVerify,
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

	// TODO: Some sort of sanity checking on number of query args...

	if len(req.UrlArguments) > 0 {
		renderedFragment = fmt.Sprintf(req.Endpoint.UrlFragment, req.UrlArguments...)
	}

	protocol := "https"
	if req.Endpoint.Protocol != "" {
		protocol = req.Endpoint.Protocol
	}

	url := fmt.Sprintf(URL_TEMPLATE, protocol, c.config.Hostname, req.Endpoint.Application, renderedFragment)

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

type requestArgs struct {
	Endpoint     *ApiEndpoint
	UrlArguments []any
	RequestBody  io.Reader
	Query        *url.Values
}

func (c *Client) doRequest(req *requestArgs) ([]byte, error) {
	renderedUrl := c.renderUrl(req)

	if req.Endpoint.HasRequestBody && req.RequestBody == http.NoBody {
		c.log.WithFields(logrus.Fields{
			"urlFragment": req.Endpoint.UrlFragment,
		}).Fatal("Request should have a body but http.NoBody was passed")
	}

	request, err := http.NewRequest(req.Endpoint.Method, renderedUrl, req.RequestBody)
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
		return nil, fmt.Errorf("got unexpected http code %d when requesting '%s'", resp.StatusCode, renderedUrl)
	}

	c.log.WithFields(logrus.Fields{
		"url":    renderedUrl,
		"status": resp.StatusCode,
	}).Debug("URL request success")

	return body, nil
}

func (pc *protectV1Client) Info() (*types.ProtectInfo, error) {
	body, err := pc.client.doRequest(&requestArgs{Endpoint: protectAPI["Info"]})
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
	body, err := c.doRequest(&requestArgs{Endpoint: protectAPI["protect/cameras"]})
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
	body, err := c.doRequest(&requestArgs{
		Endpoint:     protectAPI["protect/camera/id"],
		UrlArguments: []any{cameraID},
	})
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
			c.log.WithFields(logrus.Fields{
				"url":   url,
				"error": err.Error(),
			}).Error("WebSocket.Ping returned error")
			return
		}
		c.log.WithFields(logrus.Fields{
			"url":          url,
			"next_ping_at": next,
		}).Debug("WebSocket.Ping (Keep-Alive) Success")
	}
}

func (c *Client) SubscribeProtectEvents(ctx context.Context) (<-chan *ProtectEventMessage, error) {
	url := c.renderUrl(&requestArgs{
		Endpoint: protectAPI["protect/subscribe/events"],
	})
	conn, _, err := websocket.Dial(ctx,
		url,
		&websocket.DialOptions{
			HTTPClient: c.client,
			HTTPHeader: *c.webSocketHeaders()})
	if err != nil {
		return nil, err
	}

	c.log.WithFields(logrus.Fields{
		"url": url,
	}).Info("WebSocket.Dial() success")

	go c.webSocketKeepAlive(ctx, conn, url)

	eventChan := make(chan *ProtectEventMessage)

	go func() {
		for {
			messageType, data, err := conn.Read(ctx)
			if err != nil {
				c.log.WithFields(logrus.Fields{
					"url":   url,
					"error": err.Error(),
				}).Error("WebSocket Read returned error")
				close(eventChan)
				return
			}

			var protectEvent types.ProtectEvent
			err = json.Unmarshal(data, &protectEvent)
			if err != nil {
				c.log.WithFields(logrus.Fields{
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
	url := c.renderUrl(&requestArgs{
		Endpoint: protectAPI["protect/subscribe/devices"],
	})
	conn, _, err := websocket.Dial(ctx,
		url,
		&websocket.DialOptions{
			HTTPClient: c.client,
			HTTPHeader: *c.webSocketHeaders()})
	if err != nil {
		return nil, err
	}

	c.log.WithFields(logrus.Fields{
		"url": url,
	}).Info("WebSocket Dial Success")

	go c.webSocketKeepAlive(ctx, conn, url)

	eventChan := make(chan *ProtectDeviceEventMessage)

	go func() {
		for {
			messageType, data, err := conn.Read(ctx)
			if err != nil {
				c.log.WithFields(logrus.Fields{
					"url":   url,
					"error": err.Error(),
				}).Error("json.Unmarshal returned error")
				close(eventChan)
				return
			}

			var protectDeviceUpdate types.ProtectDeviceEvent
			err = json.Unmarshal(data, &protectDeviceUpdate)
			if err != nil {
				c.log.WithFields(logrus.Fields{
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

func (nc *networkV1Client) Info() (*types.NetworkInfo, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint: networkAPI["Info"],
	})
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

func (nc *networkV1Client) Sites(filter types.Filter, pageArgs *types.PageArguments) ([]*types.Site, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint: networkAPI["Sites"],
	})
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

func (nc *networkV1Client) Clients(siteID types.SiteID, filter types.Filter, pageArgs *types.PageArguments) ([]*types.Client, error) {
	// FIXME: Deal with pagination and filter!
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["Clients"],
		UrlArguments: []any{siteID},
	})
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

func (nc *networkV1Client) ClientDetails(siteID types.SiteID, clientID types.ClientID) (*types.Client, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["ClientDetails"],
		UrlArguments: []any{siteID, clientID},
	})
	if err != nil {
		return nil, err
	}

	var client *types.Client

	err = json.Unmarshal(body, &client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (nc *networkV1Client) ClientExecuteAction(siteID types.SiteID, clientID types.ClientID, action *types.ClientActionRequest) error {
	jsonBody, err := json.Marshal(action)
	nc.client.log.WithFields(logrus.Fields{
		"method": "ClientExecuteAction",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(jsonBody)
	_, err = nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["ClientExecuteAction"],
		UrlArguments: []any{siteID, clientID},
		RequestBody:  bodyReader,
	})

	return err
}

func (nc *networkV1Client) Devices(siteID types.SiteID, pageArgs *types.PageArguments) ([]*types.DeviceListEntry, error) {
	// FIXME: We have to send url query args
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["Devices"],
		UrlArguments: []any{siteID},
	})
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

func (nc *networkV1Client) DeviceDetails(siteID types.SiteID, deviceID types.DeviceID) (*types.Device, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["DeviceDetails"],
		UrlArguments: []any{siteID, deviceID},
	})
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

func (nc *networkV1Client) DeviceStatistics(siteID types.SiteID, deviceID types.DeviceID) (*types.DeviceStatistics, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["DeviceStatistics"],
		UrlArguments: []any{siteID, deviceID},
	})
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

func (nc *networkV1Client) DeviceExecuteAction(siteID types.SiteID, deviceID types.DeviceID, action *types.DeviceActionRequest) error {
	jsonBody, err := json.Marshal(action)
	nc.client.log.WithFields(logrus.Fields{
		"method": "DeviceExecuteAction",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(jsonBody)
	_, err = nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["DeviceExecuteAction"],
		UrlArguments: []any{siteID, deviceID},
		RequestBody:  bodyReader,
	})

	return err
}

func (nc *networkV1Client) DevicePortExecuteAction(siteID types.SiteID, deviceID types.DeviceID, port types.PortIdx, action *types.DevicePortActionRequest) error {
	jsonBody, err := json.Marshal(action)
	nc.client.log.WithFields(logrus.Fields{
		"method": "DevicePortExecuteAction",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(jsonBody)
	_, err = nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["DevicePortExecuteAction"],
		UrlArguments: []any{siteID, deviceID},
		RequestBody:  bodyReader,
	})

	return err
}

func (nc *networkV1Client) Vouchers(siteID types.SiteID, filter types.Filter, pageArgs *types.PageArguments) ([]*types.Voucher, error) {
	// FIXME: We have to send url query args: ie filter and page args
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["Vouchers"],
		UrlArguments: []any{siteID},
	})
	if err != nil {
		return nil, err
	}

	var voucherListPage *types.VoucherListPage

	err = json.Unmarshal(body, &voucherListPage)
	if err != nil {
		return nil, err
	}

	// FIXME: Deal with pagination!

	return voucherListPage.Data, nil
}

func (nc *networkV1Client) VoucherDetails(siteID types.SiteID, voucherID types.VoucherID) (*types.Voucher, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["VoucherDetails"],
		UrlArguments: []any{siteID, voucherID},
	})
	if err != nil {
		return nil, err
	}

	var voucher *types.Voucher

	err = json.Unmarshal(body, &voucher)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (nc *networkV1Client) VoucherGenerate(siteID types.SiteID, request *types.VoucherGenerateRequest) ([]*types.Voucher, error) {
	jsonBody, err := json.Marshal(request)
	nc.client.log.WithFields(logrus.Fields{
		"method": "VoucherGenerate",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["VoucherGenerate"],
		UrlArguments: []any{siteID},
		RequestBody:  bodyReader,
	})
	if err != nil {
		return nil, err
	}

	var voucherGenerateResponse *types.VoucherGenerateResponse

	err = json.Unmarshal(body, &voucherGenerateResponse)
	if err != nil {
		return nil, err
	}

	return voucherGenerateResponse.Vouchers, nil
}

func (nc *networkV1Client) VoucherDelete(siteID types.SiteID, voucherID types.VoucherID) (*types.VoucherDeleteResponse, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["VoucherDelete"],
		UrlArguments: []any{siteID, voucherID},
	})
	if err != nil {
		return nil, err
	}

	var voucherDeleteResponse *types.VoucherDeleteResponse
	err = json.Unmarshal(body, &voucherDeleteResponse)
	if err != nil {
		return nil, err
	}

	return voucherDeleteResponse, nil
}

func (nc *networkV1Client) VoucherDeleteByFilter(siteID types.SiteID, filter types.Filter) (*types.VoucherDeleteResponse, error) {
	// FIXME: Send filter!!
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["VoucherDeleteByFilter"],
		UrlArguments: []any{siteID},
	})
	if err != nil {
		return nil, err
	}

	var voucherDeleteResponse *types.VoucherDeleteResponse

	err = json.Unmarshal(body, &voucherDeleteResponse)
	if err != nil {
		return nil, err
	}

	return voucherDeleteResponse, nil
}
