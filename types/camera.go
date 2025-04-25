package types

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
