package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {

	networkCmd.AddCommand(devicesCmd)
	networkCmd.AddCommand(sitesCmd)

	devicesCmd.AddCommand(listDevicesCmd)
	devicesCmd.AddCommand(deviceDetailsCmd)
	// devicesCmd.AddCommand(latestStatisticsDevicesCmd)
	// devicesCmd.AddCommand(actionDevicesCmd)

	sitesCmd.AddCommand(listSitesCmd)
}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Make UniFi Network API calls",
	Long:  `Complete access to UniFi's Network API from the command line`,
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Make UniFi Network `devices` calls",
	Long:  `Call device endpoints under UniFi Network's API.`,
}

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Make UniFi Network `sites` calls",
	Long:  `Call sites endpoints under UniFi Network's API.`,
}

var listDevicesCmd = &cobra.Command{
	Use:   "list [site ID]",
	Short: "List all adopted UniFi Network devices by a specific site",
	Long: `Calls the devices UniFi Network API endpoint for a specific site ID
and prints the results to stdout.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		devices, err := c.ListAllDevices(args[0])
		if err != nil {
			log.Error(err.Error())
			return
		}
		for _, device := range devices {
			log.Infof("%+v", device)
		}
	},
}

var deviceDetailsCmd = &cobra.Command{
	Use:   "details [site ID] [device ID]",
	Short: "Get detailed information about a specific adopted device",
	Long: `Get detailed information about a specific adopted device. 
Response includes more information about a single device, as well 
as more detailed information about device features, such as switch 
ports and/or access point radios`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		device, err := c.GetDeviceDetails(args[0], args[1])
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Infof("%+v", device)
	},
}

var listSitesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sites managed by the Network application",
	Long: `List local sites managed by this Network application (paginated). 
Setups using Multi-Site option enabled will return all created sites, 
while if option is disabled it will return just the default site.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		sites, err := c.ListAllSites()
		if err != nil {
			log.Error(err.Error())
			return
		}
		for _, site := range sites {
			log.Infof("%+v", site)
		}
	},
}
