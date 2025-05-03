// Package cmd exposes unifi API command through an CLI interface
//
// Copyright © 2025 Clif Houck <me@clifhouck.com>
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ClifHouck/unified/client"
)

var log *logrus.Logger

var rootCmd = &cobra.Command{
	Use:   "unified",
	Short: "Make UniFi Network or Protect API calls",
	Long:  `Allows a user to talk to UniFi application APIs.`,
}

var (
	ctx                context.Context
	cfgFile            string
	hostname           string
	apiKey             string
	keepAliveInterval  time.Duration
	insecureSkipVerify bool
	debugLogging       bool
)

func getClientConfig() *client.Config {
	config := &client.Config{
		Hostname:                   hostname,
		APIKey:                     apiKey,
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
	return client.NewClient(ctx, getClientConfig(), log)
}

func Execute() {
	rootCmd.AddCommand(networkCmd)
	rootCmd.AddCommand(protectCmd)

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
	ctx = context.Background()
	log = logrus.New()

	// Config file.
	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is $HOME/.unified.yaml)")

	// Universal command flags
	rootCmd.PersistentFlags().StringVar(&hostname, "host", "unifi",
		"Hostname of UniFi API")
	// TODO: Maybe only expose this for websocket calls
	rootCmd.PersistentFlags().DurationVar(&keepAliveInterval, "keep-alive-interval", time.Second*30,
		"Interval between keep-alive pings sent for websocket streams")
	rootCmd.PersistentFlags().BoolVar(&insecureSkipVerify, "insecure", true,
		"Skip verification of UniFi TLS certificate.")

	rootCmd.PersistentFlags().BoolVar(&debugLogging, "debug", false, "Enable debug logging")

	initConfig()

	cobra.OnInitialize(configureLog)
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
	log.Out = os.Stderr

	log.SetLevel(logrus.InfoLevel)
	if debugLogging {
		log.SetLevel(logrus.DebugLevel)
	}
}
