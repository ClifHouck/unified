package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestGetConfigValuesFromViper(t *testing.T) {
	// Save original values
	origAPIKey := apiKey
	origHostname := hostname
	origKeepAliveInterval := keepAliveInterval
	origInsecureSkipVerify := insecureSkipVerify
	origLog := log

	// Setup test environment
	log = logrus.New()
	log.SetLevel(logrus.PanicLevel) // Suppress log output during tests

	// Reset viper state
	viper.Reset()

	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, ".unified.yaml")
	configContent := `host: "test-host"
apiKey: "test-api-key"
keepAliveInterval: 60s
insecure: false
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Set up viper to read the config file
	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Call the function being tested
	getConfigValuesFromViper()

	// Verify that the global variables were set correctly
	if apiKey != "test-api-key" {
		t.Errorf("Expected apiKey to be 'test-api-key', got '%s'", apiKey)
	}
	if hostname != "test-host" {
		t.Errorf("Expected hostname to be 'test-host', got '%s'", hostname)
	}
	if keepAliveInterval != 60*time.Second {
		t.Errorf("Expected keepAliveInterval to be 60s, got %v", keepAliveInterval)
	}
	if insecureSkipVerify != false {
		t.Errorf("Expected insecureSkipVerify to be false, got %v", insecureSkipVerify)
	}

	// Restore original values
	apiKey = origAPIKey
	hostname = origHostname
	keepAliveInterval = origKeepAliveInterval
	insecureSkipVerify = origInsecureSkipVerify
	log = origLog
	viper.Reset()
}

func TestConfigFileRespected(t *testing.T) {
	// Save original values
	origAPIKey := apiKey
	origHostname := hostname
	origKeepAliveInterval := keepAliveInterval
	origInsecureSkipVerify := insecureSkipVerify
	origCfgFile := cfgFile
	origLog := log

	// Setup test environment
	log = logrus.New()
	log.SetLevel(logrus.PanicLevel) // Suppress log output during tests

	// Reset viper state and global variables
	viper.Reset()
	hostname = "unifi"
	keepAliveInterval = 30 * time.Second
	insecureSkipVerify = true

	// Create a temporary config file with custom values
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, ".unified.yaml")
	configContent := `host: "192.168.1.1"
apiKey: "custom-api-key"
keepAliveInterval: 45s
insecure: false
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Set the config file flag
	cfgFile = configFile

	// Bind environment variable for API key
	err = viper.BindEnv("UNIFI_API_KEY")
	if err != nil {
		t.Fatalf("Failed to bind env var: %v", err)
	}

	// Configure viper
	viper.SetConfigType("yaml")

	// Read the config file using tryReadConfig
	if !tryReadConfig(configFile) {
		t.Fatal("tryReadConfig failed to read the config file")
	}

	// Call getConfigValuesFromViper to load values into global variables
	getConfigValuesFromViper()

	// Verify that the values from the config file are used
	if hostname != "192.168.1.1" {
		t.Errorf("Expected hostname to be '192.168.1.1' from config file, got '%s'", hostname)
	}
	if keepAliveInterval != 45*time.Second {
		t.Errorf("Expected keepAliveInterval to be 45s from config file, got %v", keepAliveInterval)
	}
	if insecureSkipVerify != false {
		t.Errorf("Expected insecureSkipVerify to be false from config file, got %v", insecureSkipVerify)
	}
	if apiKey != "custom-api-key" {
		t.Errorf("Expected apiKey to be 'custom-api-key' from config file, got '%s'", apiKey)
	}

	// Restore original values
	apiKey = origAPIKey
	hostname = origHostname
	keepAliveInterval = origKeepAliveInterval
	insecureSkipVerify = origInsecureSkipVerify
	cfgFile = origCfgFile
	log = origLog
	viper.Reset()
}
