/*
  Copyright © 2025 Clif Houck <me@clifhouck.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ClifHouck/unified/client"
)

var log = logrus.New()

var rootCmd = &cobra.Command{
	Use:   "unified",
	Short: "Make UniFi Network or Protect API calls",
	Long:  `Allows a user to talk to UniFi application APIs.`,
}

var (
	cfgFile            string
	hostname           string
	apiKey             string
	keepAliveInterval  time.Duration
	insecureSkipVerify bool
)

func getClientConfig() *client.Config {
	config := &client.Config{
		Hostname:                   hostname,
		ApiKey:                     apiKey,
		WebSocketKeepAliveInterval: keepAliveInterval,
		InsecureSkipVerify:         insecureSkipVerify,
	}

	ok, reasons := config.IsValid()
	if !ok {
		log.Error("UniFi client configuration is invalid!")
		for _, reason := range reasons {
			log.Errorf("Reason: %s", reason)
		}
		log.Fatal("Configuration must be fixed in order to use this command.")
	}

	return config
}

func getClient() *client.Client {
	return client.NewClient(getClientConfig())
}

func Execute() {
	rootCmd.AddCommand(networkCmd)

	if apiKey == "" {
		apiKey = viper.GetString("unifi_api_key")
		log.Debug("UniFi API key set from viper.")
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	configureLog()

	// Config file.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.unified.yaml)")

	// Universal command flags
	rootCmd.PersistentFlags().StringVar(&hostname, "host", "unifi",
		"Hostname of UniFi API")
	rootCmd.PersistentFlags().DurationVar(&keepAliveInterval, "keep-alive-interval", time.Duration(time.Second*30),
		"Interval between keep-alive pings sent for websocket streams")
	rootCmd.PersistentFlags().BoolVar(&insecureSkipVerify, "insecure", true,
		"Skip verification of UniFi TLS certificate.")

	initConfig()
}

func initConfig() {
	// TODO: Implement logic to draw client configuration from:
	// config, then env, then CLI args, in that order.
	viper.MustBindEnv("UNIFI_API_KEY")
}

func MarshalAndPrintJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func configureLog() {
	// TODO: Allow configuration or flag to set logging to different level.
	log.Level = logrus.InfoLevel
	log.Out = os.Stderr
}
