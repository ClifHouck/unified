// Package cmd exposes unifi API command through an CLI interface
//
// Copyright © 2025 Clif Houck <me@clifhouck.com>
package cmd

import (
	"context"
	"encoding/json"
	"errors"
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
	traceLogging       bool
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
	rootCmd.PersistentFlags().BoolVar(&traceLogging, "trace", false, "Enable trace logging")

	cobra.OnInitialize(configureLog)
	cobra.OnInitialize(initConfig)
}

func getAPIKey() string {
	if viper.IsSet("UNIFI_API_KEY") {
		log.Debug("UniFi API key set from environment.")
		return viper.GetString("UNIFI_API_KEY")
	} else if viper.IsSet("apiKey") {
		log.Debug("UniFi API key set from configuration file.")
		return viper.GetString("apikey")
	}
	log.Fatal("Couldn't retrieve API key from configuration.")
	return ""
}

func tryReadConfig(filename string) bool {
	inFile, err := os.Open(filename)
	if err != nil {
		log.WithFields(logrus.Fields{
			"filename": filename,
		}).Errorf("Couldn't open specified config file: %s", err.Error())
		return false
	}

	err = viper.ReadConfig(inFile)
	if err != nil {
		log.WithFields(logrus.Fields{
			"filename": filename,
		}).Errorf("Couldn't read specified config file: %s", err.Error())
		return false
	}

	return true
}

func initConfig() {
	err := viper.BindEnv("UNIFI_API_KEY")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = viper.BindPFlag("keepAliveInterval", rootCmd.Flags().Lookup("keep-alive-interval"))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = viper.BindPFlag("insecure", rootCmd.Flags().Lookup("insecure"))
	if err != nil {
		log.Fatal(err.Error())
	}

	viper.SetConfigName(".unified.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath("$HOME/.unified/")
	viper.AddConfigPath(".")

	// If a config file is specified as a flag, try to load it first.
	if cfgFile != "" {
		if tryReadConfig(cfgFile) {
			log.WithFields(logrus.Fields{
				"file": cfgFile,
			}).Debug("Unified config file loaded")

			apiKey = getAPIKey()
			return
		}
	}

	// Fallback to default config locations.
	err = viper.ReadInConfig()
	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Debug("Unified config file not found")
		}
		log.Error(err.Error())
	}

	log.WithFields(logrus.Fields{
		"file": viper.ConfigFileUsed(),
	}).Debug("Unified config file loaded")

	apiKey = getAPIKey()

	logConfig()
}

func logConfig() {
	log.WithFields(logrus.Fields{
		"host":               hostname,
		"isAPIKeySet":        len(apiKey) > 0,
		"insecureSkipVerify": insecureSkipVerify,
		"keepAliveInterval":  keepAliveInterval.String(),
	}).Debug("Config values")
}

func marshalAndPrintJSON(v any) error {
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
	if traceLogging {
		log.SetLevel(logrus.TraceLevel)
	}
}
