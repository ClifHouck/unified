package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/coder/websocket"
	"github.com/sirupsen/logrus"

	"github.com/ClifHouck/unified/types"
)

var protectAPI = map[string]*apiEndpoint{
	"Info": {
		URLFragment: "meta/info",
		Method:      http.MethodGet,
		Description: "Get application information",
		Application: "protect",
	},
	"SubscribeProtectEvents": {
		URLFragment: "subscribe/events",
		Method:      http.MethodGet,
		Description: "Get Protect event messages",
		Application: "protect",
	},
	"SubscribeDeviceEvents": {
		URLFragment: "subscribe/devices",
		Method:      http.MethodGet,
		Description: "Get Protect device updates",
		Application: "protect",
	},
	"Cameras": {
		URLFragment: "cameras",
		Method:      http.MethodGet,
		Description: "Get all cameras",
		Application: "protect",
	},
	"CameraDetails": {
		URLFragment: "cameras/%s",
		Method:      http.MethodGet,
		Description: "Get camera details",
		Application: "protect",
		NumURLArgs:  1,
	},
	"Viewers": {
		URLFragment: "viewers",
		Method:      http.MethodGet,
		Description: "Get all viewers",
		Application: "protect",
	},
	"ViewerDetails": {
		URLFragment: "viewers/%s",
		Method:      http.MethodGet,
		Description: "Get viewer details",
		Application: "protect",
		NumURLArgs:  1,
	},
	"ViewerSettings": {
		URLFragment:    "viewers/%s",
		Method:         http.MethodPatch,
		Description:    "Patch the settings for a specific viewer",
		Application:    "protect",
		NumURLArgs:     1,
		HasRequestBody: true,
	},
}

type protectV1Client struct {
	client *Client
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

func (pc *protectV1Client) Viewers() ([]*types.Viewer, error) {
	body, err := pc.client.doRequest(&requestArgs{Endpoint: protectAPI["Viewers"]})
	if err != nil {
		return nil, err
	}

	var viewers []*types.Viewer

	err = json.Unmarshal(body, &viewers)
	if err != nil {
		return nil, err
	}

	return viewers, nil
}

func (pc *protectV1Client) ViewerDetails(viewerID types.ViewerID) (*types.Viewer, error) {
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["viewerDetails"],
		URLArguments: []any{viewerID},
	})
	if err != nil {
		return nil, err
	}

	var viewer *types.Viewer

	err = json.Unmarshal(body, &viewer)
	if err != nil {
		return nil, err
	}

	return viewer, nil
}

func (pc *protectV1Client) ViewerSettings(
	viewerID types.ViewerID,
	settings *types.ViewerSettingsRequest,
) (*types.Viewer, error) {
	jsonBody, err := json.Marshal(settings)
	pc.client.log.WithFields(logrus.Fields{
		"method": "ViewerSettings",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["ViewerSettings"],
		URLArguments: []any{viewerID},
		RequestBody:  bodyReader,
	})
	if err != nil {
		return nil, err
	}

	var viewer *types.Viewer

	err = json.Unmarshal(body, &viewer)
	if err != nil {
		return nil, err
	}

	return viewer, nil
}

func (pc *protectV1Client) Cameras() ([]*types.Camera, error) {
	body, err := pc.client.doRequest(&requestArgs{Endpoint: protectAPI["Cameras"]})
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

func (pc *protectV1Client) CameraDetails(cameraID types.CameraID) (*types.Camera, error) {
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["CameraDetails"],
		URLArguments: []any{cameraID},
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

func (pc *protectV1Client) SubscribeProtectEvents() (<-chan *types.ProtectEvent, error) {
	url := pc.client.renderURL(&requestArgs{
		Endpoint: protectAPI["SubscribeProtectEvents"],
	})
	conn, _, err := websocket.Dial(pc.client.ctx,
		url,
		&websocket.DialOptions{
			HTTPClient: pc.client.client,
			HTTPHeader: *pc.client.webSocketHeaders()})
	if err != nil {
		return nil, err
	}

	pc.client.log.WithFields(logrus.Fields{
		"url": url,
	}).Info("WebSocket.Dial() success")

	go pc.client.webSocketKeepAlive(conn, url)

	eventChan := make(chan *types.ProtectEvent)

	go func() {
		for {
			// Make sure context is good.
			select {
			case <-pc.client.ctx.Done():
				pc.client.log.WithFields(logrus.Fields{
					"url": url,
				}).Trace("Context done.")
				return
			default:
			}

			messageType, data, readErr := conn.Read(pc.client.ctx)
			if readErr != nil {
				pc.client.log.WithFields(logrus.Fields{
					"url":   url,
					"error": readErr.Error(),
				}).Error("WebSocket Read returned error")
				close(eventChan)
				return
			}

			if messageType != websocket.MessageText {
				pc.client.log.WithFields(logrus.Fields{
					"url": url,
				}).Error("Got unhandled websocket message type!")
				close(eventChan)
				return
			}

			var protectEvent *types.ProtectEvent
			err = json.Unmarshal(data, &protectEvent)
			if err != nil {
				pc.client.log.WithFields(logrus.Fields{
					"url":   url,
					"error": err.Error(),
				}).Error("json.Unmarshal returned error")
				close(eventChan)
				return
			}

			eventChan <- protectEvent
		}
	}()

	return eventChan, nil
}

func (pc *protectV1Client) SubscribeDeviceEvents() (<-chan *types.ProtectDeviceEvent, error) {
	url := pc.client.renderURL(&requestArgs{
		Endpoint: protectAPI["SubscribeDeviceEvents"],
	})
	conn, _, err := websocket.Dial(pc.client.ctx,
		url,
		&websocket.DialOptions{
			HTTPClient: pc.client.client,
			HTTPHeader: *pc.client.webSocketHeaders()})
	if err != nil {
		return nil, err
	}

	pc.client.log.WithFields(logrus.Fields{
		"url": url,
	}).Info("WebSocket Dial Success")

	go pc.client.webSocketKeepAlive(conn, url)

	eventChan := make(chan *types.ProtectDeviceEvent)

	go func() {
		for {
			// Make sure context is good.
			select {
			case <-pc.client.ctx.Done():
				pc.client.log.WithFields(logrus.Fields{
					"url": url,
				}).Trace("Context done.")
				return
			default:
			}

			messageType, data, readErr := conn.Read(pc.client.ctx)
			if readErr != nil {
				pc.client.log.WithFields(logrus.Fields{
					"url":   url,
					"error": readErr.Error(),
				}).Error("json.Unmarshal returned error")
				close(eventChan)
				return
			}

			if messageType != websocket.MessageText {
				pc.client.log.WithFields(logrus.Fields{
					"url": url,
				}).Error("Got unhandled websocket message type!")
				close(eventChan)
				return
			}

			var protectDeviceUpdate *types.ProtectDeviceEvent
			unmarshalErr := json.Unmarshal(data, &protectDeviceUpdate)
			if unmarshalErr != nil {
				pc.client.log.WithFields(logrus.Fields{
					"url":   url,
					"error": err.Error(),
				}).Error("json.Unmarshal returned error")
				close(eventChan)
				return
			}

			eventChan <- protectDeviceUpdate
		}
	}()

	return eventChan, nil
}
