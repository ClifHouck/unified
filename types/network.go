package types

import (
	"time"
)

type Page struct {
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	Count      int `json:"count"`
	TotalCount int `json:"totalCount"`
}

type DeviceListPage struct {
	Data []*DeviceListEntry `json:"data"`
	Page
}

type DeviceListEntry struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Model      string   `json:"model"`
	MacAddress string   `json:"macAddress"`
	IPAddress  string   `json:"ipAddress"`
	State      string   `json:"state"`
	Features   []string `json:"features"`
	Interfaces []string `json:"interfaces"`
}

type Device struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Model             string    `json:"model"`
	Supported         bool      `json:"supported"`
	MacAddress        string    `json:"macAddress"`
	IPAddress         string    `json:"ipAddress"`
	State             string    `json:"state"`
	FirmwareVersion   string    `json:"firmwareVersion"`
	FirmwareUpdatable bool      `json:"firmwareUpdatable"`
	AdoptedAt         time.Time `json:"adoptedAt"`
	ProvisionedAt     time.Time `json:"provisionedAt"`
	ConfigurationID   string    `json:"configurationId"`
	Uplink            struct {
		DeviceID string `json:"deviceId"`
	} `json:"uplink"`
	Features struct {
		Switching struct {
		} `json:"switching"`
		AccessPoint struct {
		} `json:"accessPoint"`
	} `json:"features"`
	Interfaces struct {
		Ports []struct {
			Idx          int    `json:"idx"`
			State        string `json:"state"`
			Connector    string `json:"connector"`
			MaxSpeedMbps int    `json:"maxSpeedMbps"`
			SpeedMbps    int    `json:"speedMbps"`
		} `json:"ports"`
		Radios []struct {
			WlanStandard string `json:"wlanStandard"`
			// TODO: This differed from UniFi API docs. It's not a string.
			FrequencyGHz    float64 `json:"frequencyGHz"`
			ChannelWidthMHz int     `json:"channelWidthMHz"`
			Channel         int     `json:"channel"`
		} `json:"radios"`
	} `json:"interfaces"`
}

type SiteListPage struct {
	Data []*Site `json:"data"`
	Page
}

type Site struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
