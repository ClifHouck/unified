package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ClifHouck/unified/types"
)

// TODO: Maybe move this to an api module?
var networkAPI = map[string]*apiEndpoint{
	// Application related
	"Info": &apiEndpoint{
		URLFragment: "info",
		Method:      http.MethodGet,
		Description: "Get application information",
		Application: "network",
	},

	// Site related
	"Sites": &apiEndpoint{
		URLFragment: "sites",
		Method:      http.MethodGet,
		Description: "List local sites managed by this Network application",
		Application: "network",
	},

	/// Client related
	"Clients": &apiEndpoint{
		URLFragment:  "sites/%s/clients",
		Method:       http.MethodGet,
		Description:  "List clients of a site",
		Application:  "network",
		NumURLArgs:   1,
		NumQueryArgs: 1,
	},
	"ClientDetails": &apiEndpoint{
		URLFragment: "sites/%s/clients/%s",
		Method:      http.MethodGet,
		Description: "Get client details",
		Application: "network",
		NumURLArgs:  2,
	},
	"ClientExecuteAction": &apiEndpoint{
		URLFragment: "sites/%s/devices/%s/actions",
		Method:      http.MethodPost,
		Description: "Execute an action on a device",
		Application: "network",
		NumURLArgs:  2,
	},

	// Devices related
	"Devices": &apiEndpoint{
		URLFragment: "sites/%s/devices",
		Method:      http.MethodGet,
		Description: "List adopted devices of a site",
		Application: "network",
		NumURLArgs:  1,
	},
	"DeviceDetails": &apiEndpoint{
		URLFragment: "sites/%s/devices/%s",
		Method:      http.MethodGet,
		Description: "Get device details",
		Application: "network",
		NumURLArgs:  2,
	},
	"DeviceStatistics": &apiEndpoint{
		URLFragment: "sites/%s/devices/%s/statistics/latest",
		Method:      http.MethodGet,
		Description: "Get latest device statistics",
		Application: "network",
		NumURLArgs:  2,
	},
	"DeviceExecuteAction": &apiEndpoint{
		URLFragment:    "sites/%s/devices/%s/actions",
		Method:         http.MethodPost,
		Description:    "Execute an action on a device",
		Application:    "network",
		NumURLArgs:     2,
		HasRequestBody: true,
	},
	"DevicePortExecuteAction": &apiEndpoint{
		URLFragment:    "sites/%s/devices/%s/actions",
		Method:         http.MethodPost,
		Description:    "Execute an action on a device",
		Application:    "network",
		NumURLArgs:     3,
		HasRequestBody: true,
	},

	// Voucher related
	"Vouchers": &apiEndpoint{
		URLFragment:  "sites/%s/hotspot/vouchers",
		Method:       http.MethodGet,
		Description:  "List vouchers of a site",
		Application:  "network",
		NumURLArgs:   1,
		NumQueryArgs: 3,
	},
	"VoucherDetails": &apiEndpoint{
		URLFragment: "sites/%s/hotspot/vouchers/%s",
		Method:      http.MethodGet,
		Description: "Get voucher details",
		Application: "network",
		NumURLArgs:  2,
	},
	"VoucherGenerate": &apiEndpoint{
		URLFragment:    "sites/%s/hotspot/vouchers",
		Method:         http.MethodPost,
		ExpectedStatus: http.StatusCreated,
		Description:    "Generate vouchers",
		Application:    "network",
		NumURLArgs:     1,
		HasRequestBody: true,
	},
	"VoucherDelete": &apiEndpoint{
		URLFragment: "sites/%s/hotspot/vouchers/%s",
		Method:      http.MethodDelete,
		Description: "Delete vouchers",
		Application: "network",
		NumURLArgs:  2,
	},
	"VoucherDeleteByFilter": &apiEndpoint{
		URLFragment:  "sites/%s/hotspot/vouchers",
		Method:       http.MethodDelete,
		Description:  "Delete vouchers by filter",
		Application:  "network",
		NumURLArgs:   1,
		NumQueryArgs: 1,
	},
}

type networkV1Client struct {
	client *Client
}

func (nc *networkV1Client) Info() (*types.NetworkInfo, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint: networkAPI["Info"],
	})
	if err != nil {
		return nil, err
	}

	var info types.NetworkInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (nc *networkV1Client) Sites(filter types.Filter, pageArgs *types.PageArguments) ([]*types.Site, *types.Page, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint: networkAPI["Sites"],
		Query:    buildQuery(filter, pageArgs),
	})
	if err != nil {
		return nil, nil, err
	}

	var siteListPage *types.SiteListPage

	err = json.Unmarshal(body, &siteListPage)
	if err != nil {
		return nil, nil, err
	}

	return siteListPage.Data,
		&types.Page{
			Offset:     siteListPage.Offset,
			Limit:      siteListPage.Limit,
			Count:      siteListPage.Count,
			TotalCount: siteListPage.TotalCount,
		}, nil
}

func (nc *networkV1Client) Clients(siteID types.SiteID, filter types.Filter, pageArgs *types.PageArguments) ([]*types.Client, *types.Page, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["Clients"],
		URLArguments: []any{siteID},
		Query:        buildQuery(filter, pageArgs),
	})
	if err != nil {
		return nil, nil, err
	}

	var clientListPage *types.ClientListPage

	err = json.Unmarshal(body, &clientListPage)
	if err != nil {
		return nil, nil, err
	}

	return clientListPage.Data,
		&types.Page{
			Offset:     clientListPage.Offset,
			Limit:      clientListPage.Limit,
			Count:      clientListPage.Count,
			TotalCount: clientListPage.TotalCount,
		}, nil
}

func (nc *networkV1Client) ClientDetails(siteID types.SiteID, clientID types.ClientID) (*types.Client, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["ClientDetails"],
		URLArguments: []any{siteID, clientID},
	})
	if err != nil {
		return nil, err
	}

	var client *types.Client

	err = json.Unmarshal(body, &client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (nc *networkV1Client) ClientExecuteAction(siteID types.SiteID, clientID types.ClientID, action *types.ClientActionRequest) error {
	jsonBody, err := json.Marshal(action)
	nc.client.log.WithFields(logrus.Fields{
		"method": "ClientExecuteAction",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(jsonBody)
	_, err = nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["ClientExecuteAction"],
		URLArguments: []any{siteID, clientID},
		RequestBody:  bodyReader,
	})

	return err
}

func (nc *networkV1Client) Devices(siteID types.SiteID, pageArgs *types.PageArguments) ([]*types.DeviceListEntry, *types.Page, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["Devices"],
		URLArguments: []any{siteID},
		Query:        buildQuery(types.Filter(""), pageArgs),
	})
	if err != nil {
		return nil, nil, err
	}

	var deviceListPage *types.DeviceListPage

	err = json.Unmarshal(body, &deviceListPage)
	if err != nil {
		return nil, nil, err
	}

	return deviceListPage.Data,
		&types.Page{
			Offset:     deviceListPage.Offset,
			Limit:      deviceListPage.Limit,
			Count:      deviceListPage.Count,
			TotalCount: deviceListPage.TotalCount,
		}, nil
}

func (nc *networkV1Client) DeviceDetails(siteID types.SiteID, deviceID types.DeviceID) (*types.Device, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["DeviceDetails"],
		URLArguments: []any{siteID, deviceID},
	})
	if err != nil {
		return nil, err
	}

	var device *types.Device

	err = json.Unmarshal(body, &device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (nc *networkV1Client) DeviceStatistics(siteID types.SiteID, deviceID types.DeviceID) (*types.DeviceStatistics, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["DeviceStatistics"],
		URLArguments: []any{siteID, deviceID},
	})
	if err != nil {
		return nil, err
	}

	var stats *types.DeviceStatistics

	err = json.Unmarshal(body, &stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (nc *networkV1Client) DeviceExecuteAction(siteID types.SiteID, deviceID types.DeviceID, action *types.DeviceActionRequest) error {
	jsonBody, err := json.Marshal(action)
	nc.client.log.WithFields(logrus.Fields{
		"method": "DeviceExecuteAction",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(jsonBody)
	_, err = nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["DeviceExecuteAction"],
		URLArguments: []any{siteID, deviceID},
		RequestBody:  bodyReader,
	})

	return err
}

func (nc *networkV1Client) DevicePortExecuteAction(siteID types.SiteID, deviceID types.DeviceID, port types.PortIdx, action *types.DevicePortActionRequest) error {
	jsonBody, err := json.Marshal(action)
	nc.client.log.WithFields(logrus.Fields{
		"method": "DevicePortExecuteAction",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(jsonBody)
	_, err = nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["DevicePortExecuteAction"],
		URLArguments: []any{siteID, deviceID},
		RequestBody:  bodyReader,
	})

	return err
}

func (nc *networkV1Client) Vouchers(siteID types.SiteID, filter types.Filter, pageArgs *types.PageArguments) ([]*types.Voucher, *types.Page, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["Vouchers"],
		URLArguments: []any{siteID},
		Query:        buildQuery(filter, pageArgs),
	})
	if err != nil {
		return nil, nil, err
	}

	var voucherListPage *types.VoucherListPage

	err = json.Unmarshal(body, &voucherListPage)
	if err != nil {
		return nil, nil, err
	}

	return voucherListPage.Data,
		&types.Page{
			Offset:     voucherListPage.Offset,
			Limit:      voucherListPage.Limit,
			Count:      voucherListPage.Count,
			TotalCount: voucherListPage.TotalCount,
		}, nil
}

func (nc *networkV1Client) VoucherDetails(siteID types.SiteID, voucherID types.VoucherID) (*types.Voucher, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["VoucherDetails"],
		URLArguments: []any{siteID, voucherID},
	})
	if err != nil {
		return nil, err
	}

	var voucher *types.Voucher

	err = json.Unmarshal(body, &voucher)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (nc *networkV1Client) VoucherGenerate(siteID types.SiteID, request *types.VoucherGenerateRequest) ([]*types.Voucher, error) {
	jsonBody, err := json.Marshal(request)
	nc.client.log.WithFields(logrus.Fields{
		"method": "VoucherGenerate",
		"body":   string(jsonBody),
	}).Trace("Request body")
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(jsonBody)
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["VoucherGenerate"],
		URLArguments: []any{siteID},
		RequestBody:  bodyReader,
	})
	if err != nil {
		return nil, err
	}

	var voucherGenerateResponse *types.VoucherGenerateResponse

	err = json.Unmarshal(body, &voucherGenerateResponse)
	if err != nil {
		return nil, err
	}

	return voucherGenerateResponse.Vouchers, nil
}

func (nc *networkV1Client) VoucherDelete(siteID types.SiteID, voucherID types.VoucherID) (*types.VoucherDeleteResponse, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["VoucherDelete"],
		URLArguments: []any{siteID, voucherID},
	})
	if err != nil {
		return nil, err
	}

	var voucherDeleteResponse *types.VoucherDeleteResponse
	err = json.Unmarshal(body, &voucherDeleteResponse)
	if err != nil {
		return nil, err
	}

	return voucherDeleteResponse, nil
}

func (nc *networkV1Client) VoucherDeleteByFilter(siteID types.SiteID, filter types.Filter) (*types.VoucherDeleteResponse, error) {
	body, err := nc.client.doRequest(&requestArgs{
		Endpoint:     networkAPI["VoucherDeleteByFilter"],
		URLArguments: []any{siteID},
		Query:        buildQuery(filter, nil),
	})
	if err != nil {
		return nil, err
	}

	var voucherDeleteResponse *types.VoucherDeleteResponse

	err = json.Unmarshal(body, &voucherDeleteResponse)
	if err != nil {
		return nil, err
	}

	return voucherDeleteResponse, nil
}
