package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ClifHouck/unified/types"
)

func init() {
	networkCmd.AddCommand(devicesCmd)
	networkCmd.AddCommand(sitesCmd)
	networkCmd.AddCommand(clientsCmd)
	networkCmd.AddCommand(networkInfoCmd)

	devicesCmd.AddCommand(listDevicesCmd)
	devicesCmd.AddCommand(deviceDetailsCmd)
	devicesCmd.AddCommand(statsDevicesCmd)
	devicesCmd.AddCommand(actionDevicesCmd)

	sitesCmd.AddCommand(listSitesCmd)

	clientsCmd.AddCommand(listClientsCmd)
}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Make UniFi Network API calls",
	Long:  `Complete access to UniFi's Network API from the command line`,
}

var networkInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get network application info",
	Long:  `Get generic information about the Network application`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		info, err := c.NetworkInfo()
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = MarshalAndPrintJSON(info)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Make UniFi Network `devices` calls",
	Long:  `Call devices endpoints under UniFi Network's API.`,
}

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Make UniFi Network `sites` calls",
	Long:  `Call sites endpoints under UniFi Network's API.`,
}

var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "Make UniFi Network `clients` calls",
	Long:  `Call clients endpoints under UniFi Network's API.`,
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
			err := MarshalAndPrintJSON(device)
			if err != nil {
				log.Error(err.Error())
				return
			}
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
		err = MarshalAndPrintJSON(device)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var statsDevicesCmd = &cobra.Command{
	Use:   "stats [site ID] [device ID]",
	Short: "Get latest (live) statistics of a specific adopted device.",
	Long: `Get latest (live) statistics of a specific adopted device.
Response contains latest readings from a single device, such as CPU and
memory utilization, uptime, uplink tx/rx rates etc`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		stats, err := c.GetDeviceStats(args[0], args[1])
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = MarshalAndPrintJSON(stats)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var actionDevicesCmd = &cobra.Command{
	Use:   "action [site ID] [device ID] [action]",
	Short: "Execute an action on a specific adopted device",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()

		action := &types.DeviceActionRequest{
			Action: args[2],
		}

		err := c.ExecuteDeviceAction(args[0], args[1], action)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("Request success: 200 OK")
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
			err := MarshalAndPrintJSON(site)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	},
}

var listClientsCmd = &cobra.Command{
	Use:   "list [site ID]",
	Short: "List connected clients of a site",
	Long: `List connected clients of a site (paginated). Clients are either
physical devices (computers, smartphones, connected by wire or wirelessly),
or active VPN connections.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		clients, err := c.ListAllClients(args[0])
		if err != nil {
			log.Error(err.Error())
			return
		}
		for _, client := range clients {
			err := MarshalAndPrintJSON(client)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	},
}
