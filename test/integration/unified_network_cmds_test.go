package integration_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ClifHouck/unified/types"
)

const (
	unifiedBinary = "../../build/unified"
)

func checkForUniFiAPIHostSkip(t *testing.T) {
	if os.Getenv("UNIFIED_HAVE_UNIFI_API_HOST") == "" {
		t.Skip("Must set environment variable 'UNIFIED_HAVE_UNIFI_API_HOST' to " +
			"run this test. Requires an available UniFi API at 'https://unifi' " +
			"or set the hostname at 'UNIFIED_UNIFI_API_HOSTNAME'")
	}
}

type TestNetworkIDSet struct {
	SiteID   string
	DeviceID string
	ClientID string
}

func helperSeedIDValues(t *testing.T) *TestNetworkIDSet {
	var idSet TestNetworkIDSet

	cmd := exec.Command("../../build/unified", "network", "sites", "list", "--id-only")
	output, err := cmd.Output()
	require.NoError(t, err)
	idSet.SiteID = strings.Split(string(output), "\n")[0]

	cmd = exec.Command(
		"../../build/unified",
		"network",
		"devices",
		"list",
		idSet.SiteID,
		"--id-only",
	)
	output, err = cmd.Output()
	require.NoError(t, err)
	idSet.DeviceID = strings.Split(string(output), "\n")[0]

	cmd = exec.Command(
		"../../build/unified",
		"network",
		"clients",
		"list",
		idSet.SiteID,
		"--id-only",
	)
	output, err = cmd.Output()
	require.NoError(t, err)
	idSet.ClientID = strings.Split(string(output), "\n")[0]

	return &idSet
}

// Test all Network application GET commands.
// TODO: POST/PATCH/DELETE commands are a bit trickier to test...
func TestUnifiedCmdNetworkGETCommands(t *testing.T) {
	checkForUniFiAPIHostSkip(t)

	idSet := helperSeedIDValues(t)

	type TestCase struct {
		Name        string
		Command     []string
		NeedsSite   bool
		NeedsDevice bool
		NeedsClient bool
	}

	networkTestCases := []*TestCase{
		{
			Name:    "Test 'network info'",
			Command: []string{"network", "info"},
		},
		{
			Name:    "Test 'network sites list'",
			Command: []string{"network", "sites", "list"},
		},
		{
			Name:      "Test 'network devices list'",
			Command:   []string{"network", "devices", "list"},
			NeedsSite: true,
		},
		{
			Name:        "Test 'network devices detail'",
			Command:     []string{"network", "devices", "details"},
			NeedsSite:   true,
			NeedsDevice: true,
		},
		{
			Name:        "Test 'network devices stats'",
			Command:     []string{"network", "devices", "stats"},
			NeedsSite:   true,
			NeedsDevice: true,
		},
		{
			Name:      "Test 'network clients list'",
			Command:   []string{"network", "clients", "list"},
			NeedsSite: true,
		},
		{
			Name:        "Test 'network clients details'",
			Command:     []string{"network", "clients", "details"},
			NeedsSite:   true,
			NeedsClient: true,
		},
		// TODO: Support voucher calls
	}

	for _, tc := range networkTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.NeedsSite {
				tc.Command = append(tc.Command, idSet.SiteID)
			}
			if tc.NeedsDevice {
				tc.Command = append(tc.Command, idSet.DeviceID)
			}
			if tc.NeedsClient {
				tc.Command = append(tc.Command, idSet.ClientID)
			}

			tc.Command = append(tc.Command, "--debug")

			fullCmd := []string{unifiedBinary}
			fullCmd = append(fullCmd, tc.Command...)
			fmt.Println("Running Command: '" + strings.Join(fullCmd, " ") + "'")

			cmd := exec.Command(unifiedBinary, tc.Command...)
			output, err := cmd.Output()
			require.NoError(t, err)
			fmt.Print(string(output))
		})
	}
}

func TestVoucherCmdsNetwork(t *testing.T) {
	checkForUniFiAPIHostSkip(t)

	idSet := helperSeedIDValues(t)

	voucherName := "unified-integration-test-vouchers"
	voucherCount := "3"
	numVouchers := 3

	var vouchers []types.Voucher
	setup := false
	t.Run("Voucher setup", func(t *testing.T) {
		cmd := exec.Command(unifiedBinary, "network", "vouchers", "generate", idSet.SiteID,
			"--count", voucherCount,
			"--rx-limit", "2",
			"--tx-limit", "2",
			"--guest-limit", "1",
			"--time-limit", "60",
			"--name", voucherName,
			"--data-limit", "100")
		output, err := cmd.Output()
		fmt.Print(string(output))
		require.NoError(t, err)

		err = json.Unmarshal(output, &vouchers)
		require.NoError(t, err)

		assert.Len(t, vouchers, numVouchers)

		setup = true
	})
	if !setup {
		t.Fatalf("Voucher setup failed. Skipping the rest of this test.")
	}

	type TestCase struct {
		Name    string
		Command []string
		Case    func(*testing.T, []byte)
	}

	voucherTestCases := []*TestCase{
		{
			Name: "Test 'vouchers list'",
			Command: []string{
				"network",
				"vouchers",
				"list",
				idSet.SiteID,
				"--hide-page",
				"--filter",
				"name.eq('unified-integration-test-vouchers')",
			},
			Case: func(t *testing.T, output []byte) {
				var vouchers []*types.Voucher
				err := json.Unmarshal(output, &vouchers)
				require.NoError(t, err)
				assert.Len(t, vouchers, numVouchers)
				assert.Equal(t, voucherName, vouchers[0].Name)
			},
		},
		{
			Name:    "Test 'vouchers details'",
			Command: []string{"network", "vouchers", "details", idSet.SiteID, vouchers[0].ID},
			Case: func(t *testing.T, output []byte) {
				var voucher *types.Voucher
				err := json.Unmarshal(output, &voucher)
				require.NoError(t, err)
				assert.Equal(t, vouchers[0].ID, voucher.ID)
				assert.Equal(t, vouchers[0].Name, voucher.Name)
			},
		},
		{
			Name:    "Test 'vouchers delete'",
			Command: []string{"network", "vouchers", "delete", idSet.SiteID, vouchers[0].ID},
			Case: func(t *testing.T, output []byte) {
				var voucherDeleteResp *types.VoucherDeleteResponse
				err := json.Unmarshal(output, &voucherDeleteResp)
				require.NoError(t, err)
				assert.Equal(t, 1, voucherDeleteResp.VouchersDeleted)
			},
		},
		{
			Name: "Test 'vouchers delete-filter'",
			Command: []string{
				"network",
				"vouchers",
				"delete-filter",
				idSet.SiteID,
				"--filter",
				"name.eq('unified-integration-test-vouchers')",
			},
			Case: func(t *testing.T, output []byte) {
				var voucherDeleteResp *types.VoucherDeleteResponse
				err := json.Unmarshal(output, &voucherDeleteResp)
				require.NoError(t, err)
				assert.Equal(t, numVouchers-1, voucherDeleteResp.VouchersDeleted)
			},
		},
		{
			Name: "Test 'vouchers list' is now empty",
			Command: []string{
				"network",
				"vouchers",
				"list",
				idSet.SiteID,
				"--hide-page",
				"--filter",
				"name.eq('unified-integration-test-vouchers')",
			},
			Case: func(t *testing.T, output []byte) {
				err := json.Unmarshal(output, &vouchers)
				require.NoError(t, err)
				assert.Empty(t, vouchers)
			},
		},
	}

	for _, tc := range voucherTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			fullCmd := []string{unifiedBinary}
			tc.Command = append(tc.Command, "--debug")
			fullCmd = append(fullCmd, tc.Command...)
			fmt.Println("Running Command: '" + strings.Join(fullCmd, " ") + "'")

			cmd := exec.Command(unifiedBinary, tc.Command...)
			output, err := cmd.Output()
			require.NoError(t, err)
			fmt.Print(string(output))
			tc.Case(t, output)
		})
	}
}
