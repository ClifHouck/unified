package integration_test

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func checkForUniFiProtectAPIHostSkip(t *testing.T) {
	if os.Getenv("UNIFIED_HAVE_UNIFI_PROTECT_API_HOST") == "" {
		t.Skip("Must set environment variable 'UNIFIED_HAVE_UNIFI_PROTECT_API_HOST' to " +
			"run this test. Requires an available UniFi Protect API at 'https://unifi' " +
			"or set the hostname at 'UNIFIED_UNIFI_API_HOSTNAME'")
	}
}

type TestProtectIDSet struct {
	CameraID   string
	LiveviewID string
	// ViewerID   string TODO I don't have any viewers...
}

func helperSeedProtectIDValues(t *testing.T) *TestProtectIDSet {
	var idSet TestProtectIDSet

	cmd := exec.Command(
		"../../build/unified",
		"protect",
		"cameras",
		"list",
		"--id-only",
	)
	output, err := cmd.Output()
	require.NoError(t, err)
	idSet.CameraID = strings.Split(string(output), "\n")[0]

	cmd = exec.Command(
		"../../build/unified",
		"protect",
		"liveviews",
		"list",
		"--id-only",
	)
	output, err = cmd.Output()
	require.NoError(t, err)
	idSet.LiveviewID = strings.Split(string(output), "\n")[0]

	return &idSet
}

// Try to test all Protect application GET commands.
func TestUnifiedCmdProtectGETCommands(t *testing.T) {
	checkForUniFiProtectAPIHostSkip(t)

	err := os.Mkdir("/tmp/unified/", 0750)
	if !os.IsExist(err) {
		require.NoError(t, err)
	}

	idSet := helperSeedProtectIDValues(t)

	type TestCase struct {
		Name         string
		Command      []string
		AfterCommand func(*testing.T) error
	}

	protectTestCases := []*TestCase{
		{
			Name:    "Test 'protect info'",
			Command: []string{"protect", "info"},
		},
		{
			Name:    "Test 'protect liveviews list'",
			Command: []string{"protect", "liveviews", "list"},
		},
		{
			Name:    "Test 'protect liveviews details'",
			Command: []string{"protect", "liveviews", "details", idSet.LiveviewID},
		},
		{
			Name:    "Test 'protect viewers list'",
			Command: []string{"protect", "viewers", "list"},
		},
		{
			Name:    "Test 'protect cameras list'",
			Command: []string{"protect", "cameras", "list"},
		},
		{
			Name:    "Test 'protect cameras details'",
			Command: []string{"protect", "cameras", "details", idSet.CameraID},
		},
		{
			Name:    "Test 'protect cameras snapshot'",
			Command: []string{"protect", "cameras", "snapshot", idSet.CameraID, "/tmp/unified/test_snapshot.jpg"},
			AfterCommand: func(t *testing.T) error {
				// Try loading the image
				data, acErr := os.ReadFile("/tmp/unified/test_snapshot.jpg")
				require.NoError(t, acErr)

				reader := bytes.NewReader(data)
				_, acErr = jpeg.Decode(reader)
				require.NoError(t, acErr)

				// Then remove it.
				err = os.Remove("/tmp/unified/test_snapshot.jpg")
				require.NoError(t, acErr)
				return nil
			},
		},
		{
			Name:    "Test 'protect cameras get-stream'",
			Command: []string{"protect", "cameras", "stream-get", idSet.CameraID},
		},
	}

	for _, tc := range protectTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Command = append(tc.Command, "--debug")

			fullCmd := []string{unifiedBinary}
			fullCmd = append(fullCmd, tc.Command...)
			fmt.Println("Running Command: '" + strings.Join(fullCmd, " ") + "'")

			cmd := exec.Command(unifiedBinary, tc.Command...)
			output, tcErr := cmd.Output()
			require.NoError(t, tcErr)
			fmt.Print(string(output))
			if tc.AfterCommand != nil {
				acErr := tc.AfterCommand(t)
				require.NoError(t, acErr)
			}
		})
	}
}
