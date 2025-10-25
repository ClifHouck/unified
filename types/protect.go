package types

import (
	"fmt"
	"image"
	"strconv"
)

type ProtectV1 interface {
	// About application
	Info() (*ProtectInfo, error)

	// Viewer Information & Management
	Viewers() ([]*Viewer, error)
	ViewerDetails(ViewerID) (*Viewer, error)
	ViewerSettings(ViewerID, *ViewerSettingsRequest) (*Viewer, error)

	// Live View Management
	LiveViews() ([]*LiveView, error)
	LiveViewPatch(LiveViewID, *LiveView) (*LiveView, error)
	LiveViewDetails(LiveViewID) (*LiveView, error)
	LiveViewCreate(*LiveView) (*LiveView, error)
	// TODO: But where is DELETE?

	// Websocket updates
	SubscribeDeviceEvents() (<-chan *ProtectDeviceEvent, error)
	SubscribeProtectEvents() (<-chan *ProtectEvent, error)

	// Camera Information & Management
	Cameras() ([]*Camera, error)
	CameraDetails(CameraID) (*Camera, error)
	CameraPatch(CameraID, *CameraPatchRequest) (*Camera, error)

	CameraCreateRTSPSStream(CameraID, *CameraCreateRTSPSStreamRequest) (*CameraCreateRTSPSStreamResponse, error)
	CameraDeleteRTSPSStream(CameraID, *CameraDeleteRTSPSStreamRequest) error
	CameraGetRTSPSStream(CameraID) (*CameraGetRTSPSStreamResponse, error)

	CameraGetSnapshot(CameraID, bool) (image.Image, error)

	CameraDisableMicPermanently(CameraID) (*Camera, error)
	CameraTalkbackSession(CameraID) (*CameraTalkbackSessionResponse, error)

	// Camera PTZ Control & Management
	CameraPTZPatrolStart(CameraID, CameraPatrolSlotNumber) error
	CameraPTZPatrolStop(CameraID) error
	CameraPTZGotoPresetPosition(CameraID, CameraPresetPositionSlotNumber) error

	// Lights
	Lights() ([]*Light, error)
	LightDetails(LightID) (*Light, error)
	LightPatch(LightID, *LightPatchRequest) (*Light, error)

	// NVRs
	NVRs() (*NVR, error)

	// Chimes
	Chimes() ([]*Chime, error)
	ChimeDetails(ChimeID) (*Chime, error)
	ChimePatch(ChimeID, *ChimePatchRequest) (*Chime, error)

	// Sensors
	Sensors() ([]*Sensor, error)
	SensorDetails(SensorID) (*Sensor, error)
	SensorPatch(SensorID, *SensorPatchRequest) (*Sensor, error)

	// Device Asset File Management
	Files(FileType) ([]*File, error)
	FileUpload(FileType, string, []byte) error

	// Alarm Manager
	AlarmManagerWebhook(AlarmTriggerID) error
}

// CameraID is a UniFI protect Camera ID. Interestingly *not* a UUID.
type CameraID string

type ViewerID string

type LiveViewID string

type LightID string

type ChimeID string

type SensorID string

type AlarmTriggerID string

type ProtectInfo struct {
	ApplicationVersion string `json:"applicationVersion"`
}

type Camera struct {
	ID           string `json:"id"`
	ModelKey     string `json:"modelKey"`
	State        string `json:"state"`
	Name         string `json:"name"`
	IsMicEnabled bool   `json:"isMicEnabled"`
	OsdSettings  struct {
		IsNameEnabled  bool `json:"isNameEnabled"`
		IsDateEnabled  bool `json:"isDateEnabled"`
		IsLogoEnabled  bool `json:"isLogoEnabled"`
		IsDebugEnabled bool `json:"isDebugEnabled"`
	} `json:"osdSettings"`
	LedSettings struct {
		IsEnabled bool `json:"isEnabled"`
	} `json:"ledSettings"`
	LcdMessage struct {
		Type    string `json:"type"`
		ResetAt int    `json:"resetAt"`
		Text    string `json:"text"`
	} `json:"lcdMessage"`
	MicVolume        int    `json:"micVolume"`
	ActivePatrolSlot int    `json:"activePatrolSlot"`
	VideoMode        string `json:"videoMode"`
	HdrType          string `json:"hdrType"`
	FeatureFlags     struct {
		SupportFullHdSnapshot bool     `json:"supportFullHdSnapshot"`
		HasHdr                bool     `json:"hasHdr"`
		SmartDetectTypes      []string `json:"smartDetectTypes"`
		SmartDetectAudioTypes []string `json:"smartDetectAudioTypes"`
		VideoModes            []string `json:"videoModes"`
		HasMic                bool     `json:"hasMic"`
		HasLedStatus          bool     `json:"hasLedStatus"`
		HasSpeaker            bool     `json:"hasSpeaker"`
	} `json:"featureFlags"`
	SmartDetectSettings struct {
		ObjectTypes []string `json:"objectTypes"`
		AudioTypes  []string `json:"audioTypes"`
	} `json:"smartDetectSettings"`
}

type lcdMessage struct {
	Type    string `json:"type,omitempty"`
	ResetAt int    `json:"resetAt,omitempty"`
	Text    string `json:"text,omitempty"`
}

type CameraPatchRequest struct {
	Name        string `json:"name,omitempty"`
	OsdSettings struct {
		IsNameEnabled  bool `json:"isNameEnabled,omitempty"`
		IsDateEnabled  bool `json:"isDateEnabled,omitempty"`
		IsLogoEnabled  bool `json:"isLogoEnabled,omitempty"`
		IsDebugEnabled bool `json:"isDebugEnabled,omitempty"`
	} `json:"osdSettings,omitzero"`
	LedSettings struct {
		IsEnabled bool `json:"isEnabled,omitempty"`
	} `json:"ledSettings,omitzero"`
	LcdMessage          lcdMessage `json:"lcdMessage,omitzero"`
	MicVolume           int        `json:"micVolume,omitzero"`
	VideoMode           string     `json:"videoMode,omitempty"`
	HdrType             string     `json:"hdrType,omitempty"`
	SmartDetectSettings struct {
		ObjectTypes []string `json:"objectTypes,omitempty"`
		AudioTypes  []string `json:"audioTypes,omitempty"`
	} `json:"smartDetectSettings,omitzero"`
}

type Viewer struct {
	ID          string `json:"id"`
	ModelKey    string `json:"modelKey"`
	State       string `json:"state"`
	Name        string `json:"name"`
	Liveview    string `json:"liveview"`
	StreamLimit int    `json:"streamLimit"`
}

type ViewerSettingsRequest struct {
	Name     string `json:"name"`
	Liveview string `json:"liveview"`
}

type LiveView struct {
	ID        string `json:"id,omitempty"`
	ModelKey  string `json:"modelKey,omitempty"`
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault,omitempty"`
	IsGlobal  bool   `json:"isGlobal"`
	Owner     string `json:"owner,omitempty"`
	Layout    int    `json:"layout"`
	Slots     []struct {
		Cameras       []string `json:"cameras"`
		CycleMode     string   `json:"cycleMode"`
		CycleInterval int      `json:"cycleInterval"`
	} `json:"slots"`
}

type CameraCreateRTSPSStreamRequest struct {
	Qualities []string `json:"qualities"`
}

type CameraDeleteRTSPSStreamRequest struct {
	Qualities []string `json:"qualities"`
}

type CameraCreateRTSPSStreamResponse struct {
	cameraStreamQualities
}

type CameraGetRTSPSStreamResponse struct {
	cameraStreamQualities
}

type cameraStreamQualities struct {
	High    string `json:"high,omitempty"`
	Medium  string `json:"medium,omitempty"`
	Low     string `json:"low,omitempty"`
	Package string `json:"package,omitempty"`
}

type ProtectErrorMessage struct {
	Error  string `json:"error"`
	Name   string `json:"name"`
	Entity string `json:"entity"`
	Issues []struct {
		InstancePath string `json:"instancePath"`
		Message      string `json:"message"`
		Keyword      string `json:"keyword"`
	} `json:"issues"`
	Body interface{} `json:"body"`
}

type CameraTalkbackSessionResponse struct {
	URL           string `json:"url"`
	Codec         string `json:"codec"`
	SamplingRate  int    `json:"samplingRate"`
	BitsPerSample int    `json:"bitsPerSample"`
}

type Light struct {
	ID                string `json:"id"`
	ModelKey          string `json:"modelKey"`
	State             string `json:"state"`
	Name              string `json:"name"`
	LightModeSettings struct {
		Mode     string `json:"mode"`
		EnableAt string `json:"enableAt"`
	} `json:"lightModeSettings"`
	LightDeviceSettings struct {
		IsIndicatorEnabled bool `json:"isIndicatorEnabled"`
		PirDuration        int  `json:"pirDuration"`
		PirSensitivity     int  `json:"pirSensitivity"`
		LedLevel           int  `json:"ledLevel"`
	} `json:"lightDeviceSettings"`
	IsDark              bool   `json:"isDark"`
	IsLightOn           bool   `json:"isLightOn"`
	IsLightForceEnabled bool   `json:"isLightForceEnabled"`
	LastMotion          int    `json:"lastMotion"`
	IsPirMotionDetected bool   `json:"isPirMotionDetected"`
	Camera              string `json:"camera"`
}

type LightPatchRequest struct {
	Name                string `json:"name,omitempty"`
	IsLightForceEnabled bool   `json:"isLightForceEnabled,omitempty"`
	LightModeSettings   struct {
		Mode     string `json:"mode,omitempty"`
		EnableAt string `json:"enableAt,omitempty"`
	} `json:"lightModeSettings,omitzero"`
	LightDeviceSettings struct {
		IsIndicatorEnabled bool `json:"isIndicatorEnabled,omitempty"`
		PirDuration        int  `json:"pirDuration,omitempty"`
		PirSensitivity     int  `json:"pirSensitivity,omitempty"`
		LedLevel           int  `json:"ledLevel,omitempty"`
	} `json:"lightDeviceSettings,omitzero"`
}

type NVR struct {
	ID               string `json:"id"`
	ModelKey         string `json:"modelKey"`
	Name             string `json:"name"`
	DoorbellSettings struct {
		DefaultMessageText           string   `json:"defaultMessageText"`
		DefaultMessageResetTimeoutMs int      `json:"defaultMessageResetTimeoutMs"`
		CustomMessages               []string `json:"customMessages"`
		CustomImages                 []struct {
			Preview string `json:"preview"`
			Sprite  string `json:"sprite"`
		} `json:"customImages"`
	} `json:"doorbellSettings"`
}

type Chime struct {
	ID           string   `json:"id"`
	ModelKey     string   `json:"modelKey"`
	State        string   `json:"state"`
	Name         string   `json:"name"`
	CameraIDs    []string `json:"cameraIds"`
	RingSettings []struct {
		CameraID    string `json:"cameraId"`
		RepeatTimes int    `json:"repeatTimes"`
		RingtoneID  string `json:"ringtoneId"`
		Volume      int    `json:"volume"`
	} `json:"ringSettings"`
}

type ChimePatchRequest struct {
	Name         string   `json:"name,omitempty"`
	CameraIDs    []string `json:"cameraIds,omitempty"`
	RingSettings []struct {
		CameraID    string `json:"cameraId,omitempty"`
		RepeatTimes int    `json:"repeatTimes,omitempty"`
		RingtoneID  string `json:"ringtoneId,omitempty"`
		Volume      int    `json:"volume,omitempty"`
	} `json:"ringSettings,omitzero"`
}

type SensorSettings struct {
	IsEnabled     bool    `json:"isEnabled"`
	Margin        float64 `json:"margin"`
	LowThreshold  int     `json:"lowThreshold"`
	HighThreshold int     `json:"highThreshold"`
}

type Sensor struct {
	ID            string `json:"id"`
	ModelKey      string `json:"modelKey"`
	State         string `json:"state"`
	Name          string `json:"name"`
	MountType     string `json:"mountType"`
	BatteryStatus struct {
		Percentage int  `json:"percentage"`
		IsLow      bool `json:"isLow"`
	} `json:"batteryStatus"`
	Stats struct {
		Light struct {
			Value  int    `json:"value"`
			Status string `json:"status"`
		} `json:"light"`
		Humidity struct {
			Value  int    `json:"value"`
			Status string `json:"status"`
		} `json:"humidity"`
		Temperature struct {
			Value  float64 `json:"value"`
			Status string  `json:"status"`
		} `json:"temperature"`
	} `json:"stats"`
	LightSettings       SensorSettings `json:"lightSettings"`
	HumiditySettings    SensorSettings `json:"humiditySettings"`
	TemperatureSettings SensorSettings `json:"temperatureSettings"`
	IsOpened            bool           `json:"isOpened"`
	OpenStatusChangedAt int            `json:"openStatusChangedAt"`
	IsMotionDetected    bool           `json:"isMotionDetected"`
	MotionDetectedAt    int            `json:"motionDetectedAt"`
	MotionSettings      struct {
		IsEnabled   bool `json:"isEnabled"`
		Sensitivity int  `json:"sensitivity"`
	} `json:"motionSettings"`
	AlarmTriggeredAt int `json:"alarmTriggeredAt"`
	AlarmSettings    struct {
		IsEnabled bool `json:"isEnabled"`
	} `json:"alarmSettings"`
	LeakDetectedAt      int `json:"leakDetectedAt"`
	TamperingDetectedAt int `json:"tamperingDetectedAt"`
}

type SensorPatchRequest struct {
	Name                string         `json:"name,omitempty"`
	LightSettings       SensorSettings `json:"lightSettings,omitzero"`
	HumiditySettings    SensorSettings `json:"humiditySettings,omitzero"`
	TemperatureSettings SensorSettings `json:"temperatureSettings,omitzero"`
	MotionSettings      struct {
		IsEnabled   bool `json:"isEnabled,omitempty"`
		Sensitivity int  `json:"sensitivity,omitempty"`
	} `json:"motionSettings,omitzero"`
	AlarmSettings struct {
		IsEnabled bool `json:"isEnabled,omitempty"`
	} `json:"alarmSettings,omitzero"`
}

type FileType int

const (
	FileTypeAnimations = iota
)

var fileTypeToString = map[FileType]string{
	FileTypeAnimations: "animations",
}

func (fta FileType) String() string {
	return fileTypeToString[fta]
}

type File struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	OriginalName string `json:"originalName"`
	Path         string `json:"path"`
}

type SlotNumber int

const (
	slotLow  SlotNumber = 0
	slotHigh SlotNumber = 4
)

func (s SlotNumber) Valid() bool {
	return s >= slotLow && s < slotHigh
}

func (s SlotNumber) String() string {
	return strconv.Itoa(int(s))
}

type CameraPatrolSlotNumber struct {
	SlotNumber
}

type CameraPresetPositionSlotNumber struct {
	SlotNumber
}

type SlotRangeError struct {
	Slot int
}

func (sre SlotRangeError) Error() string {
	return fmt.Sprintf("Slot must be between %d and %d inclusive, got %d", slotLow, slotHigh, sre.Slot)
}
