package types

type ProtectV1 interface {
	// About application
	Info() (*ProtectInfo, error)

	// Websocket updates
	SubscribeDeviceEvents() (<-chan *ProtectDeviceEvent, error)
	SubscribeProtectEvents() (<-chan *ProtectEvent, error)

	// Camera Information & Management
	Cameras() ([]*Camera, error)
	CameraDetails(CameraID) (*Camera, error)

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

	// TODO: Rest of protect API!
}

// CameraID is a UniFI protect Camera ID. Interestingly *not* a UUID.
type CameraID string

type ViewerID string

type LiveViewID string

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
