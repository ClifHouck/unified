package types

import (
	"encoding/json"
	"fmt"
)

type ProtectEvent struct {
	Type string `json:"type"`
	// Polymorphic object that maps to an Event
	Item     interface{}     `json:"-"`
	ItemType string          `json:"-"`
	RawItem  json.RawMessage `json:"item"`
}

func (pe *ProtectEvent) UnmarshalJSON(data []byte) error {
	type event ProtectEvent

	err := json.Unmarshal(data, (*event)(pe))
	if err != nil {
		return err
	}

	var item ProtectEventItem
	err = json.Unmarshal(pe.RawItem, &item)
	if err != nil {
		return err
	}

	switch item.Type {
	case "ring":
		pe.Item = &RingEvent{}
	case "sensorExtremeValues":
		pe.Item = &SensorExtremeValuesEvent{}
	case "sensorWaterLeak":
		pe.Item = &SensorWaterLeakEvent{}
	case "sensorTamper":
		pe.Item = &SensorTamperEvent{}
	case "sensorBatteryLow":
		pe.Item = &SensorBatteryLowEvent{}
	case "sensorAlarm":
		pe.Item = &SensorAlarmEvent{}
	case "sensorOpened":
		pe.Item = &SensorOpenedEvent{}
	case "sensorClosed":
		pe.Item = &SensorClosedEvent{}
	case "sensorMotion":
		pe.Item = &SensorMotionEvent{}
	case "lightMotion":
		pe.Item = &LightMotionEvent{}
	case "motion":
		pe.Item = &CameraMotionEvent{}
	case "smartAudioDetect":
		pe.Item = &CameraSmartDetectAudioEvent{}
	case "smartDetectZone":
		pe.Item = &CameraSmartDetectZoneEvent{}
	case "smartDetectLine":
		pe.Item = &CameraSmartDetectLineEvent{}
	case "smartDetectLoiterZone":
		pe.Item = &CameraSmartDetectLoiterEvent{}
	default:
		return fmt.Errorf("ProtectEvent unrecognized type '%s'", pe.Type)
	}

	err = json.Unmarshal(pe.RawItem, pe.Item)
	if err != nil {
		return err
	}

	pe.ItemType = item.Type

	return nil
}

type ProtectEventItem struct {
	ID       string `json:"id"`
	ModelKey string `json:"modelKey"`
	Type     string `json:"type"`
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	Device   string `json:"device"`
}

type TextObject struct {
	Text string `json:"text"`
}

type IntNumberObject struct {
	Number int `json:"number"`
}

type FloatNumberObject struct {
	Number float64 `json:"number"`
}

type RingEvent struct {
	ProtectEventItem
}

type SensorExtremeValuesEvent struct {
	ProtectEventItem
	Metadata struct {
		SensorType  TextObject `json:"sensorType"`
		SensorValue TextObject `json:"sensorValue"`
		Status      TextObject `json:"status"`
	} `json:"metadata"`
}

type SensorWaterLeakEvent struct {
	ProtectEventItem
	Metadata struct {
		SensorMountType TextObject `json:"sensorMountType"`
	} `json:"metadata"`
}

type SensorTamperEvent struct {
	ProtectEventItem
}

type SensorBatteryLowEvent struct {
	ProtectEventItem
	Metadata struct {
		SensorBatteryPercentage FloatNumberObject `json:"sensorBatteryPercentage"`
	} `json:"metadata"`
}

type SensorAlarmEvent struct {
	ProtectEventItem
}

type SensorOpenedEvent struct {
	ProtectEventItem
	Metadata struct {
		SensorMountType TextObject `json:"sensorMountType"`
	} `json:"metadata"`
}

type SensorClosedEvent struct {
	ProtectEventItem
	Metadata struct {
		SensorMountType TextObject `json:"sensorMountType"`
	} `json:"metadata"`
}

type SensorMotionEvent struct {
	ProtectEventItem
	Metadata struct {
		SensorMountType TextObject `json:"sensorMountType"`
	} `json:"metadata"`
}

type LightMotionEvent struct {
	ProtectEventItem
	Metadata struct {
		SensorMountType TextObject `json:"sensorMountType"`
	} `json:"metadata"`
}

type CameraMotionEvent struct {
	ProtectEventItem
	Metadata struct {
		SensorMountType TextObject `json:"sensorMountType"`
	} `json:"metadata"`
}

type CameraSmartDetectAudioEvent struct {
	ProtectEventItem
	SmartDetectTypes []string `json:"smartDetectTypes"`
}

type CameraSmartDetectZoneEvent struct {
	ProtectEventItem
	SmartDetectTypes []string `json:"smartDetectTypes"`
}

type CameraSmartDetectLineEvent struct {
	ProtectEventItem
	SmartDetectTypes []string `json:"smartDetectTypes"`
}

type CameraSmartDetectLoiterEvent struct {
	ProtectEventItem
	SmartDetectTypes []string `json:"smartDetectTypes"`
}

var AllProtectEvents = []interface{}{
	RingEvent{},
	SensorExtremeValuesEvent{},
	SensorWaterLeakEvent{},
	SensorTamperEvent{},
	SensorBatteryLowEvent{},
	SensorAlarmEvent{},
	SensorOpenedEvent{},
	SensorClosedEvent{},
	SensorMotionEvent{},
	LightMotionEvent{},
	CameraMotionEvent{},
	CameraSmartDetectAudioEvent{},
	CameraSmartDetectZoneEvent{},
	CameraSmartDetectLineEvent{},
	CameraSmartDetectLoiterEvent{},
}

type ProtectDeviceEvent struct {
	Type     string `json:"type"`
	ModelKey string `json:"modelKey"`

	// Polymorphic object that maps to an Event
	Item     interface{}     `json:"-"`
	ItemType string          `json:"-"`
	RawItem  json.RawMessage `json:"item"`
}

func (pde *ProtectDeviceEvent) UnmarshalJSON(data []byte) error {
	type event ProtectDeviceEvent

	err := json.Unmarshal(data, (*event)(pde))
	if err != nil {
		return err
	}

	var item ProtectDeviceEventItem
	err = json.Unmarshal(pde.RawItem, &item)
	if err != nil {
		return err
	}

	switch item.ModelKey {
	case "nvr":
		pde.Item = &ProtectNVREvent{}
	case "camera":
		pde.Item = &ProtectCameraEvent{}
	case "chime":
		pde.Item = &ProtectChimeEvent{}
	case "light":
		pde.Item = &ProtectLightEvent{}
	case "viewer":
		pde.Item = &ProtectViewerEvent{}
	case "speaker":
		pde.Item = &ProtectSpeakerEvent{}
	case "bridge":
		pde.Item = &ProtectBridgeEvent{}
	case "doorlock":
		pde.Item = &ProtectDoorlockEvent{}
	case "sensor":
		pde.Item = &ProtectSensorEvent{}
	case "aiProcessor":
		pde.Item = &ProtectAIProcessorEvent{}
	case "aiPort":
		pde.Item = &ProtectAIPortEvent{}
	case "linkStation":
		pde.Item = &ProtectLinkStationEvent{}
	default:
		return fmt.Errorf("ProtectDeviceEvent unrecognized type '%s'", pde.ModelKey)
	}

	err = json.Unmarshal(pde.RawItem, pde.Item)
	if err != nil {
		return err
	}

	pde.ModelKey = item.ModelKey

	return nil
}

type ProtectDeviceEventItem struct {
	ID       string `json:"id"`
	ModelKey string `json:"modelKey"`
	Name     string `json:"name"`
	State    string `json:"state"`
}

type ProtectCameraEvent Camera

type ProtectNVREvent NVR

type ProtectChimeEvent Chime

type ProtectLightEvent Light

type ProtectViewerEvent Viewer

type ProtectSpeakerEvent struct {
	ProtectDeviceEventItem
}

type ProtectBridgeEvent struct {
	ProtectDeviceEventItem
}

type ProtectDoorlockEvent struct {
	ProtectDeviceEventItem
}

type ProtectSensorEvent Sensor

type ProtectAIProcessorEvent struct {
	ProtectDeviceEventItem
}

type ProtectAIPortEvent struct {
	ProtectDeviceEventItem
}

type ProtectLinkStationEvent struct {
	ProtectDeviceEventItem
}

var AllProtectDeviceEvents = []interface{}{
	ProtectCameraEvent{},
	ProtectNVREvent{},
	ProtectChimeEvent{},
	ProtectLightEvent{},
	ProtectViewerEvent{},
	ProtectSpeakerEvent{},
	ProtectBridgeEvent{},
	ProtectDoorlockEvent{},
	ProtectSensorEvent{},
	ProtectAIProcessorEvent{},
	ProtectAIPortEvent{},
	ProtectLinkStationEvent{},
}
