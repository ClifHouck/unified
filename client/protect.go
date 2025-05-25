package client

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/coder/websocket"
	"github.com/sirupsen/logrus"

	"github.com/ClifHouck/unified/types"
)

// FIXME - This pattern is not great because these are tightly coupled with
// their individual methods. Refactor this to store the ApiEndpoint struct
// with the method somehow.
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
	"CameraPatch": {
		URLFragment:    "cameras/%s",
		Method:         http.MethodPatch,
		Description:    "Patch the settings for a specific camera",
		Application:    "protect",
		NumURLArgs:     1,
		HasRequestBody: true,
	},
	"CameraCreateRTSPSStream": {
		URLFragment:    "cameras/%s/rtsps-stream",
		Method:         http.MethodPost,
		Description:    "Returns RTSPS stream URLs for specified quality levels",
		Application:    "protect",
		NumURLArgs:     1,
		HasRequestBody: true,
	},
	"CameraDeleteRTSPSStream": {
		URLFragment:    "cameras/%s/rtsps-stream",
		Method:         http.MethodDelete,
		Description:    "Removes the RTSPS stream for a specified camera",
		Application:    "protect",
		NumURLArgs:     1,
		NumQueryArgs:   1,
		ExpectedStatus: http.StatusNoContent,
	},
	"CameraGetRTSPSStream": {
		URLFragment: "cameras/%s/rtsps-stream",
		Description: "Gets existing RTSPS streams for a specified camera",
		Application: "protect",
		NumURLArgs:  1,
	},
	"CameraGetSnapshot": {
		URLFragment:  "cameras/%s/snapshot",
		Method:       http.MethodGet,
		Description:  "Get camera details",
		Application:  "protect",
		NumURLArgs:   1,
		NumQueryArgs: 1,
	},
	"CameraDisableMicPermanently": {
		URLFragment: "cameras/%s/disable-mic-permanently",
		Method:      http.MethodPost,
		Description: "Disable the microphone for a specific camera",
		Application: "protect",
		NumURLArgs:  1,
	},
	"CameraTalkbackSession": {
		URLFragment: "cameras/%s/talkback-session",
		Method:      http.MethodPost,
		Description: "Get the talkback stream URL and audio config for a camera",
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
	"LiveViews": {
		URLFragment: "liveviews",
		Method:      http.MethodGet,
		Description: "Get all liveviews",
		Application: "protect",
	},
	"LiveViewDetails": {
		URLFragment: "liveviews/%s",
		Method:      http.MethodGet,
		Description: "Get liveview details",
		Application: "protect",
		NumURLArgs:  1,
	},
	"LiveViewCreate": {
		URLFragment:    "liveviews",
		Method:         http.MethodPost,
		Description:    "Create a new live view",
		Application:    "protect",
		HasRequestBody: true,
	},
	"LiveViewPatch": {
		URLFragment:    "liveviews/%s",
		Method:         http.MethodPatch,
		Description:    "Patch the configuration about a specific live view",
		Application:    "protect",
		HasRequestBody: true,
		NumURLArgs:     1,
	},
	"Lights": {
		URLFragment: "lights",
		Method:      http.MethodGet,
		Description: "Get all lights",
		Application: "protect",
	},
	"LightDetails": {
		URLFragment: "lights/%s",
		Method:      http.MethodGet,
		Description: "Get light details",
		Application: "protect",
		NumURLArgs:  1,
	},
	"LightPatch": {
		URLFragment:    "lights/%s",
		Method:         http.MethodPatch,
		Description:    "Patch the settings for a specific light",
		Application:    "protect",
		NumURLArgs:     1,
		HasRequestBody: true,
	},
	"NVRs": {
		URLFragment: "nvrs",
		Method:      http.MethodGet,
		Description: "Get detailed information about the NVR",
		Application: "protect",
	},
	"Chimes": {
		URLFragment: "chimes",
		Method:      http.MethodGet,
		Description: "Get all chimes",
		Application: "protect",
	},
	"ChimeDetails": {
		URLFragment: "chimes/%s",
		Method:      http.MethodGet,
		Description: "Get chime details",
		Application: "protect",
		NumURLArgs:  1,
	},
	"ChimePatch": {
		URLFragment:    "chimes/%s",
		Method:         http.MethodPatch,
		Description:    "Patch the settings for a specific chime",
		Application:    "protect",
		NumURLArgs:     1,
		HasRequestBody: true,
	},
	"Sensors": {
		URLFragment: "sensors",
		Method:      http.MethodGet,
		Description: "Get all sensors",
		Application: "protect",
	},
	"SensorDetails": {
		URLFragment: "sensors/%s",
		Method:      http.MethodGet,
		Description: "Get sensor details",
		Application: "protect",
		NumURLArgs:  1,
	},
	"SensorPatch": {
		URLFragment:    "sensors/%s",
		Method:         http.MethodPatch,
		Description:    "Patch the settings for a specific sensor",
		Application:    "protect",
		NumURLArgs:     1,
		HasRequestBody: true,
	},
	"Files": {
		URLFragment: "files/%s",
		Method:      http.MethodGet,
		Description: "Get device assset files by type",
		NumURLArgs:  1,
		Application: "protect",
	},
	"FileUpload": {
		URLFragment:    "files/%s",
		Method:         http.MethodPost,
		Description:    "Upload device asset file",
		NumURLArgs:     1,
		Application:    "protect",
		HasRequestBody: true,
	},
	"AlarmManagerWebhook": {
		URLFragment:    "alarm-manager/webhook/%s",
		Method:         http.MethodPost,
		Description:    "Send a webhook to the alarm manager to trigger configured alarms",
		NumURLArgs:     1,
		Application:    "protect",
		ExpectedStatus: http.StatusNoContent,
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

func (pc *protectV1Client) CameraPatch(
	cameraID types.CameraID,
	camera *types.CameraPatchRequest,
) (*types.Camera, error) {
	jsonBody, err := json.Marshal(camera)
	pc.client.log.WithFields(logrus.Fields{
		"method": "CameraPatch",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["CameraPatch"],
		URLArguments: []any{cameraID},
		RequestBody:  bodyReader,
	})
	if err != nil {
		return nil, err
	}

	var updatedCamera *types.Camera

	err = json.Unmarshal(body, &updatedCamera)
	if err != nil {
		return nil, err
	}

	return updatedCamera, nil
}

func (pc *protectV1Client) CameraCreateRTSPSStream(
	cameraID types.CameraID,
	req *types.CameraCreateRTSPSStreamRequest,
) (*types.CameraCreateRTSPSStreamResponse, error) {
	jsonBody, err := json.Marshal(req)
	pc.client.log.WithFields(logrus.Fields{
		"method": "CameraCreateRTSPSStream",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["CameraCreateRTSPSStream"],
		URLArguments: []any{cameraID},
		RequestBody:  bodyReader,
	})
	if err != nil {
		return nil, err
	}

	var resp *types.CameraCreateRTSPSStreamResponse

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (pc *protectV1Client) CameraDeleteRTSPSStream(
	cameraID types.CameraID,
	req *types.CameraDeleteRTSPSStreamRequest,
) error {
	query := &url.Values{}
	for _, qual := range req.Qualities {
		query.Add("qualities[]", qual)
	}

	_, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["CameraDeleteRTSPSStream"],
		URLArguments: []any{cameraID},
		Query:        query,
	})
	if err != nil {
		return err
	}

	return nil
}

func (pc *protectV1Client) CameraGetRTSPSStream(
	cameraID types.CameraID,
) (*types.CameraGetRTSPSStreamResponse, error) {
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["CameraGetRTSPSStream"],
		URLArguments: []any{cameraID},
	})
	if err != nil {
		return nil, err
	}

	var resp *types.CameraGetRTSPSStreamResponse

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (pc *protectV1Client) CameraGetSnapshot(
	cameraID types.CameraID,
	highQuality bool,
) (image.Image, error) {
	query := &url.Values{}
	quality := "false"
	if highQuality {
		quality = "true"
	}
	query.Add("highQuality", quality)

	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["CameraGetSnapshot"],
		URLArguments: []any{cameraID},
		Query:        query,
	})
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(body)
	image, err := jpeg.Decode(reader)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (pc *protectV1Client) CameraDisableMicPermanently(
	cameraID types.CameraID,
) (*types.Camera, error) {
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["CameraDisableMicPermanently"],
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

func (pc *protectV1Client) CameraTalkbackSession(
	cameraID types.CameraID,
) (*types.CameraTalkbackSessionResponse, error) {
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["CameraTalkbackSession"],
		URLArguments: []any{cameraID},
	})
	if err != nil {
		return nil, err
	}

	var cameraTalkbackResp *types.CameraTalkbackSessionResponse

	err = json.Unmarshal(body, &cameraTalkbackResp)
	if err != nil {
		return nil, err
	}

	return cameraTalkbackResp, nil
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

func (pc *protectV1Client) LiveViews() ([]*types.LiveView, error) {
	body, err := pc.client.doRequest(&requestArgs{Endpoint: protectAPI["LiveViews"]})
	if err != nil {
		return nil, err
	}

	var liveViews []*types.LiveView

	err = json.Unmarshal(body, &liveViews)
	if err != nil {
		return nil, err
	}

	return liveViews, nil
}

func (pc *protectV1Client) LiveViewDetails(liveViewID types.LiveViewID) (*types.LiveView, error) {
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["LiveViewDetails"],
		URLArguments: []any{liveViewID},
	})
	if err != nil {
		return nil, err
	}

	var liveView *types.LiveView

	err = json.Unmarshal(body, &liveView)
	if err != nil {
		return nil, err
	}

	return liveView, nil
}

func (pc *protectV1Client) LiveViewCreate(lv *types.LiveView) (*types.LiveView, error) {
	jsonBody, err := json.Marshal(lv)
	pc.client.log.WithFields(logrus.Fields{
		"method": "LiveViewCreate",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:    protectAPI["LiveViewCreate"],
		RequestBody: bodyReader,
	})
	if err != nil {
		return nil, err
	}

	var liveView *types.LiveView

	err = json.Unmarshal(body, &liveView)
	if err != nil {
		return nil, err
	}

	return liveView, nil
}

func (pc *protectV1Client) LiveViewPatch(
	liveViewID types.LiveViewID,
	lv *types.LiveView,
) (*types.LiveView, error) {
	jsonBody, err := json.Marshal(lv)
	pc.client.log.WithFields(logrus.Fields{
		"method": "LiveViewPatch",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["LiveViewPatch"],
		URLArguments: []any{liveViewID},
		RequestBody:  bodyReader,
	})
	if err != nil {
		return nil, err
	}

	var liveView *types.LiveView

	err = json.Unmarshal(body, &liveView)
	if err != nil {
		return nil, err
	}

	return liveView, nil
}

func (pc *protectV1Client) Lights() ([]*types.Light, error) {
	body, err := pc.client.doRequest(&requestArgs{Endpoint: protectAPI["Lights"]})
	if err != nil {
		return nil, err
	}

	var lights []*types.Light

	err = json.Unmarshal(body, &lights)
	if err != nil {
		return nil, err
	}

	return lights, nil
}

func (pc *protectV1Client) LightDetails(lightID types.LightID) (*types.Light, error) {
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["LightDetails"],
		URLArguments: []any{lightID},
	})
	if err != nil {
		return nil, err
	}

	var light *types.Light

	err = json.Unmarshal(body, &light)
	if err != nil {
		return nil, err
	}

	return light, nil
}

func (pc *protectV1Client) LightPatch(
	lightID types.LightID,
	light *types.LightPatchRequest,
) (*types.Light, error) {
	jsonBody, err := json.Marshal(light)
	pc.client.log.WithFields(logrus.Fields{
		"method": "LightPatch",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["LightPatch"],
		URLArguments: []any{lightID},
		RequestBody:  bodyReader,
	})
	if err != nil {
		return nil, err
	}

	var updatedLight *types.Light

	err = json.Unmarshal(body, &updatedLight)
	if err != nil {
		return nil, err
	}

	return updatedLight, nil
}

func (pc *protectV1Client) NVRs() (*types.NVR, error) {
	body, err := pc.client.doRequest(&requestArgs{Endpoint: protectAPI["NVRs"]})
	if err != nil {
		return nil, err
	}

	var nvr *types.NVR

	err = json.Unmarshal(body, &nvr)
	if err != nil {
		return nil, err
	}

	return nvr, nil
}

func (pc *protectV1Client) Chimes() ([]*types.Chime, error) {
	body, err := pc.client.doRequest(&requestArgs{Endpoint: protectAPI["Chimes"]})
	if err != nil {
		return nil, err
	}

	var chimes []*types.Chime

	err = json.Unmarshal(body, &chimes)
	if err != nil {
		return nil, err
	}

	return chimes, nil
}

func (pc *protectV1Client) ChimeDetails(chimeID types.ChimeID) (*types.Chime, error) {
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["ChimeDetails"],
		URLArguments: []any{chimeID},
	})
	if err != nil {
		return nil, err
	}

	var chime *types.Chime

	err = json.Unmarshal(body, &chime)
	if err != nil {
		return nil, err
	}

	return chime, nil
}

func (pc *protectV1Client) ChimePatch(
	chimeID types.ChimeID,
	chime *types.ChimePatchRequest,
) (*types.Chime, error) {
	jsonBody, err := json.Marshal(chime)
	pc.client.log.WithFields(logrus.Fields{
		"method": "ChimePatch",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["ChimePatch"],
		URLArguments: []any{chimeID},
		RequestBody:  bodyReader,
	})
	if err != nil {
		return nil, err
	}

	var updatedChime *types.Chime

	err = json.Unmarshal(body, &updatedChime)
	if err != nil {
		return nil, err
	}

	return updatedChime, nil
}

func (pc *protectV1Client) Sensors() ([]*types.Sensor, error) {
	body, err := pc.client.doRequest(&requestArgs{Endpoint: protectAPI["Sensors"]})
	if err != nil {
		return nil, err
	}

	var sensors []*types.Sensor

	err = json.Unmarshal(body, &sensors)
	if err != nil {
		return nil, err
	}

	return sensors, nil
}

func (pc *protectV1Client) SensorDetails(sensorID types.SensorID) (*types.Sensor, error) {
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["SensorDetails"],
		URLArguments: []any{sensorID},
	})
	if err != nil {
		return nil, err
	}

	var sensor *types.Sensor

	err = json.Unmarshal(body, &sensor)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

func (pc *protectV1Client) SensorPatch(
	sensorID types.SensorID,
	sensor *types.SensorPatchRequest,
) (*types.Sensor, error) {
	jsonBody, err := json.Marshal(sensor)
	pc.client.log.WithFields(logrus.Fields{
		"method": "SensorPatch",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	body, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["SensorPatch"],
		URLArguments: []any{sensorID},
		RequestBody:  bodyReader,
	})
	if err != nil {
		return nil, err
	}

	var updatedSensor *types.Sensor

	err = json.Unmarshal(body, &updatedSensor)
	if err != nil {
		return nil, err
	}

	return updatedSensor, nil
}

func (pc *protectV1Client) Files(fileType types.FileType) ([]*types.File, error) {
	body, err := pc.client.doRequest(
		&requestArgs{
			Endpoint:     protectAPI["Files"],
			URLArguments: []any{fileType.String()},
		})
	if err != nil {
		return nil, err
	}

	var files []*types.File

	err = json.Unmarshal(body, &files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (pc *protectV1Client) FileUpload(fileType types.FileType, filename string, contents []byte) error {
	buf := new(bytes.Buffer)
	mpBodyWriter := multipart.NewWriter(buf)

	formFile, err := createFormFileProtect(mpBodyWriter, "file", filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(formFile, bytes.NewReader(contents))
	if err != nil {
		return err
	}
	err = mpBodyWriter.Close()
	if err != nil {
		return err
	}

	endpoint := protectAPI["FileUpload"]
	endpoint.ContentType = mpBodyWriter.FormDataContentType()

	pc.client.log.WithFields(logrus.Fields{
		"filename":     filename,
		"filetype":     fileType.String(),
		"content-type": endpoint.ContentType,
	}).Trace("Uploading file...")

	_, err = pc.client.doRequest(
		&requestArgs{
			Endpoint:     endpoint,
			URLArguments: []any{fileType.String()},
			RequestBody:  buf,
		})
	return err
}

func (pc *protectV1Client) AlarmManagerWebhook(triggerID types.AlarmTriggerID) error {
	_, err := pc.client.doRequest(&requestArgs{
		Endpoint:     protectAPI["AlarmManagerWebhook"],
		URLArguments: []any{triggerID},
	})
	return err
}
