package types

import "fmt"

import "encoding/json"

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
		return fmt.Errorf("type '%s'", pe.Type)
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

var ALL_PROTECT_EVENTS = []interface{}{
	RingEvent{},
	SensorExtremeValuesEvent{},
	SensorWaterLeakEvent{},
	SensorTamperEvent{},
	SensorBatteryLowEvent{},
	SensorAlarmEvent{},
	SensorOpenedEvent{},
	SensorClosedEvent{},
	LightMotionEvent{},
	CameraMotionEvent{},
	CameraSmartDetectAudioEvent{},
	CameraSmartDetectZoneEvent{},
	CameraSmartDetectLineEvent{},
	CameraSmartDetectLoiterEvent{},
}

type ProtectDeviceEvent struct {
	Type string `json:"type"`

	// Polymorphic object that maps to an Event
	Item     interface{}     `json:"-"`
	ItemType string          `json:"-"`
	RawItem  json.RawMessage `json:"item"`
}

type ProtectDeviceEventItem struct {
	ID       string `json:"id"`
	ModelKey string `json:"modelKey"`
	Name     string `json:"name"`
	State    string `json:"state"`
}

type ProtectAddCameraEvent struct {
	ProtectDeviceEventItem
	// FIXME: Add the reset of device events!
}

var ALL_PROTECT_DEVICE_EVENTS = []interface{}{
	ProtectAddCameraEvent{},
}
