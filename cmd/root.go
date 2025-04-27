/*
  Copyright © 2025 Clif Houck <me@clifhouck.com>
*/
package cmd

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ClifHouck/unified/client"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "unified",
	Short: "Make Unifi Network or Protect API calls",
	Long:  `Allows a user to talk to Unifi application APIs.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var (
	cfgFile            string
	hostname           string
	apiKey             string
	keepAliveInterval  time.Duration
	insecureSkipVerify bool
)

func getClientConfig() *client.Config {
	return &client.Config{
		Hostname:                   hostname,
		ApiKey:                     apiKey,
		WebSocketKeepAliveInterval: keepAliveInterval,
		InsecureSkipVerify:         insecureSkipVerify,
	}
}

func getClient() *client.Client {
	return client.NewClient(getClientConfig())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(networkCmd)

	if apiKey == "" {
		apiKey = viper.GetString("unifi_api_key")
		log.Info("Unifi API key set from viper.")
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Config file.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.unified.yaml)")

	// Universal command flags
	rootCmd.PersistentFlags().StringVar(&hostname, "host", "unifi",
		"Hostname of unifi API")
	rootCmd.PersistentFlags().DurationVar(&keepAliveInterval, "keep-alive-interval", time.Duration(time.Second*30),
		"Interval between keep-alive pings sent for websocket streams")
	rootCmd.PersistentFlags().BoolVar(&insecureSkipVerify, "insecure", true,
		"Skip verification of unifi TLS certificate.")

	initConfig()
}

func initConfig() {
	viper.MustBindEnv("UNIFI_API_KEY")
}
