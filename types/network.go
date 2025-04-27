package types

import (
	"time"
)

type NetworkInfo struct {
	ApplicationVersion string `json:"applicationVersion"`
}

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

type DeviceStatistics struct {
	UptimeSec            int64     `json:"uptimeSec"`
	LastHeartbeatAt      time.Time `json:"lastHeartbeatAt"`
	NextHeartbeatAt      time.Time `json:"nextHeartbeatAt"`
	LoadAverage1Min      float64   `json:"loadAverage1Min"`
	LoadAverage5Min      float64   `json:"loadAverage5Min"`
	LoadAverage15Min     float64   `json:"loadAverage15Min"`
	CPUUtilizationPct    float64   `json:"cpuUtilizationPct"`
	MemoryUtilizationPct float64   `json:"memoryUtilizationPct"`
	Uplink               struct {
		TxRateBps int64 `json:"txRateBps"`
		RxRateBps int64 `json:"rxRateBps"`
	} `json:"uplink"`
	Interfaces struct {
		Radios []struct {
			// TODO: This differed from UniFi API docs. It's not a string.
			FrequencyGHz float64 `json:"frequencyGHz"`
			TxRetriesPct float64 `json:"txRetriesPct"`
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

type ClientListPage struct {
	Data []*Client `json:"data"`
	Page
}

type Client struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ConnectedAt time.Time `json:"connectedAt"`
	IPAddress   string    `json:"ipAddress"`
}

type DeviceActionRequest struct {
	Action string `json:"action"`
}
