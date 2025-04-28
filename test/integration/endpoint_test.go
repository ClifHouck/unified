package integration

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	idSet.SiteID = strings.Split(string(output), "\n")[0]

	cmd = exec.Command("../../build/unified", "network", "devices", "list", idSet.SiteID, "--id-only")
	output, err = cmd.Output()
	assert.NoError(t, err)
	idSet.DeviceID = strings.Split(string(output), "\n")[0]

	cmd = exec.Command("../../build/unified", "network", "clients", "list", idSet.SiteID, "--id-only")
	output, err = cmd.Output()
	assert.NoError(t, err)
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
		&TestCase{
			Name:    "Test 'network info'",
			Command: []string{"network", "info"},
		},
		&TestCase{
			Name:    "Test 'network sites list'",
			Command: []string{"network", "sites", "list"},
		},
		&TestCase{
			Name:      "Test 'network devices list'",
			Command:   []string{"network", "devices", "list"},
			NeedsSite: true,
		},
		&TestCase{
			Name:        "Test 'network devices detail'",
			Command:     []string{"network", "devices", "details"},
			NeedsSite:   true,
			NeedsDevice: true,
		},
		&TestCase{
			Name:        "Test 'network devices stats'",
			Command:     []string{"network", "devices", "stats"},
			NeedsSite:   true,
			NeedsDevice: true,
		},
		&TestCase{
			Name:      "Test 'network clients list'",
			Command:   []string{"network", "clients", "list"},
			NeedsSite: true,
		},
		&TestCase{
			Name:        "Test 'network clients details'",
			Command:     []string{"network", "clients", "details"},
			NeedsSite:   true,
			NeedsClient: true,
		},
	}

	binaryName := "../../build/unified"

	for _, tc := range networkTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.NeedsSite {
				tc.Command = append(tc.Command, idSet.SiteID)
			}
			if tc.NeedsDevice {
				tc.Command = append(tc.Command, idSet.DeviceID)
			}

			fullCmd := []string{binaryName}
			fullCmd = append(fullCmd, tc.Command...)
			fmt.Println("Running Command: '" + strings.Join(fullCmd, " ") + "'")

			cmd := exec.Command(binaryName, tc.Command...)
			output, err := cmd.Output()
			assert.NoError(t, err)
			fmt.Print(string(output))
		})
	}
}
