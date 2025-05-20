package types

import (
	"image"
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

	// TODO: Rest of protect API!
}

// CameraID is a UniFI protect Camera ID. Interestingly *not* a UUID.
type CameraID string

type ViewerID string

type LiveViewID string

type LightID string

type ChimeID string

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
	CameraIds    []string `json:"cameraIds"`
	RingSettings []struct {
		CameraID    string `json:"cameraId"`
		RepeatTimes int    `json:"repeatTimes"`
		RingtoneID  string `json:"ringtoneId"`
		Volume      int    `json:"volume"`
	} `json:"ringSettings"`
}

type ChimePatchRequest struct {
	Name         string   `json:"name,omitempty"`
	CameraIds    []string `json:"cameraIds,omitempty"`
	RingSettings []struct {
		CameraID    string `json:"cameraId,omitempty"`
		RepeatTimes int    `json:"repeatTimes,omitempty"`
		RingtoneID  string `json:"ringtoneId,omitempty"`
		Volume      int    `json:"volume,omitempty"`
	} `json:"ringSettings,omitzero"`
}
