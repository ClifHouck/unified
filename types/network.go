package types

import (
	"time"
)

type NetworkV1 interface {
	// About application
	Info() (*NetworkInfo, error)

	// Sites
	Sites(Filter, *PageArguments) ([]*Site, *Page, error)

	// Clients
	Clients(SiteID, Filter, *PageArguments) ([]*Client, *Page, error)
	ClientDetails(SiteID, ClientID) (*Client, error)
	ClientExecuteAction(SiteID, ClientID, *ClientActionRequest) error

	// Devices
	Devices(SiteID, *PageArguments) ([]*DeviceListEntry, *Page, error)
	DeviceDetails(SiteID, DeviceID) (*Device, error)
	DeviceStatistics(SiteID, DeviceID) (*DeviceStatistics, error)
	DeviceExecuteAction(SiteID, DeviceID, *DeviceActionRequest) error
	DevicePortExecuteAction(SiteID, DeviceID, PortIdx, *DevicePortActionRequest) error

	// Vouchers
	Vouchers(SiteID, Filter, *PageArguments) ([]*Voucher, *Page, error)
	VoucherDetails(SiteID, VoucherID) (*Voucher, error)
	VoucherGenerate(SiteID, *VoucherGenerateRequest) ([]*Voucher, error)
	VoucherDelete(SiteID, VoucherID) (*VoucherDeleteResponse, error)
	VoucherDeleteByFilter(SiteID, Filter) (*VoucherDeleteResponse, error)
}

// TODO: All IDs appear to be UUIDs. Add some methods to verify IDs.
type UnifiID string

type SiteID UnifiID

type DeviceID UnifiID

type ClientID UnifiID

type VoucherID UnifiID

type PortIdx uint32

// TODO: Filter has a bunch of syntax and rules associated with it
// would like to parse and verify if possible.
type Filter string

type PageArguments struct {
	Offset uint32
	Limit  uint32
}

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
	Type        string    `json:"type"`
	// TODO: Not sure what 'access' is yet based on API docs.
	// Access      string    `json:"access"`
}

// TODO: Add action verification or enum
type ClientActionRequest struct {
	Action string `json:"action"`
}

// TODO: Add action verification or enum
type DeviceActionRequest struct {
	Action string `json:"action"`
}

// TODO: Add action verification or enum
type DevicePortActionRequest struct {
	Action string `json:"action"`
}

type VoucherListPage struct {
	Data []*Voucher `json:"data"`
	Page
}

type Voucher struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `json:"createdAt"`
	Name                 string    `json:"name"`
	Code                 string    `json:"code"`
	AuthorizedGuestLimit int       `json:"authorizedGuestLimit"`
	AuthorizedGuestCount int       `json:"authorizedGuestCount"`
	ActivatedAt          time.Time `json:"activatedAt"`
	ExpiresAt            time.Time `json:"expiresAt"`
	Expired              bool      `json:"expired"`
	TimeLimitMinutes     int       `json:"timeLimitMinutes"`
	DataUsageLimitMBytes int       `json:"dataUsageLimitMBytes"`
	RxRateLimitKbps      int       `json:"rxRateLimitKbps"`
	TxRateLimitKbps      int       `json:"txRateLimitKbps"`
}

type VoucherGenerateRequest struct {
	Count                int    `json:"count"`
	Name                 string `json:"name"`
	AuthorizedGuestLimit int    `json:"authorizedGuestLimit"`
	TimeLimitMinutes     int    `json:"timeLimitMinutes"`
	DataUsageLimitMBytes int    `json:"dataUsageLimitMBytes"`
	RxRateLimitKbps      int    `json:"rxRateLimitKbps"`
	TxRateLimitKbps      int    `json:"txRateLimitKbps"`
}

type VoucherGenerateResponse struct {
	Vouchers []*Voucher `json:"vouchers"`
}

type VoucherDeleteResponse struct {
	VouchersDeleted int `json:"vouchersDeleted"`
}
