package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ClifHouck/unified/types"
)

var idOnly = false
var listingFlagSet = pflag.NewFlagSet("listing", pflag.ExitOnError)

var filter = ""
var filterFlagSet = pflag.NewFlagSet("filter", pflag.ExitOnError)

// FIXME: Pagination flag set/args
// FIXME: Filter flag set/args

func init() {
	filterFlagSet.StringVar(&filter, "filter", "", "Filter results based on expression")
	listingFlagSet.BoolVar(&idOnly, "id-only", false, "List only the ID of listed entities, one per line.")

	networkCmd.AddCommand(networkInfoCmd)
	networkCmd.AddCommand(devicesCmd)
	networkCmd.AddCommand(sitesCmd)
	networkCmd.AddCommand(clientsCmd)
	networkCmd.AddCommand(vouchersCmd)

	listDevicesCmd.Flags().AddFlagSet(listingFlagSet)
	devicesCmd.AddCommand(listDevicesCmd)
	devicesCmd.AddCommand(deviceDetailsCmd)
	devicesCmd.AddCommand(statsDevicesCmd)
	devicesCmd.AddCommand(actionDevicesCmd)
	//devicesCmd.AddCommand(actionDevicePortCmd)

	listSitesCmd.Flags().AddFlagSet(listingFlagSet)
	sitesCmd.AddCommand(listSitesCmd)

	listClientsCmd.Flags().AddFlagSet(listingFlagSet)
	listClientsCmd.Flags().AddFlagSet(filterFlagSet)
	clientsCmd.AddCommand(listClientsCmd)
	clientsCmd.AddCommand(clientDetailsCmd)

	listVouchersCmd.Flags().AddFlagSet(listingFlagSet)
	listVouchersCmd.Flags().AddFlagSet(filterFlagSet)
	vouchersCmd.AddCommand(listVouchersCmd)
	vouchersCmd.AddCommand(voucherDetailsCmd)
	//vouchersCmd.AddCommand(voucherGenerateCmd)
	//vouchersCmd.AddCommand(voucherDeleteCmd)
	// voucherDeleteByFilterCmd.Flags().AddFlagSet(filterFlagSet)
	//vouchersCmd.AddCommand(voucherDeleteByFilterCmd)
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
		info, err := c.Network.Info()
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

var vouchersCmd = &cobra.Command{
	Use:   "vouchers",
	Short: "Make UniFi Network `vouchers` calls",
	Long:  `Call vouchers endpoints under UniFi Network's API.`,
}

var listDevicesCmd = &cobra.Command{
	Use:   "list [site ID]",
	Short: "List all adopted UniFi Network devices by a specific site",
	Long: `Calls the devices UniFi Network API endpoint for a specific site ID
and prints the results to stdout.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		devices, err := c.Network.Devices(types.SiteID(args[0]), &types.PageArguments{})
		if err != nil {
			log.Error(err.Error())
			return
		}
		for _, device := range devices {
			if idOnly {
				fmt.Println(device.ID)
			} else {
				err := MarshalAndPrintJSON(device)
				if err != nil {
					log.Error(err.Error())
					return
				}
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
		device, err := c.Network.DeviceDetails(types.SiteID(args[0]), types.DeviceID(args[1]))
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
		stats, err := c.Network.DeviceStatistics(types.SiteID(args[0]), types.DeviceID(args[1]))
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

		err := c.Network.DeviceExecuteAction(types.SiteID(args[0]), types.DeviceID(args[1]), action)
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
		// FIXME: Deal with filter and page options!!
		sites, err := c.Network.Sites(types.Filter(filter), &types.PageArguments{})
		if err != nil {
			log.Error(err.Error())
			return
		}
		for _, site := range sites {
			if idOnly {
				fmt.Println(site.ID)
			} else {
				err := MarshalAndPrintJSON(site)
				if err != nil {
					log.Error(err.Error())
					return
				}
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
		// FIXME: Deal with page options!!
		clients, err := c.Network.Clients(
			types.SiteID(args[0]),
			types.Filter(filter),
			&types.PageArguments{})
		if err != nil {
			log.Error(err.Error())
			return
		}
		for _, client := range clients {
			if idOnly {
				fmt.Println(client.ID)
			} else {
				err := MarshalAndPrintJSON(client)
				if err != nil {
					log.Error(err.Error())
					return
				}
			}
		}
	},
}

var clientDetailsCmd = &cobra.Command{
	Use:   "details [site ID] [client ID]",
	Short: "Get detailed information about a specific connected client",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		client, err := c.Network.ClientDetails(types.SiteID(args[0]), types.ClientID(args[1]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = MarshalAndPrintJSON(client)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var clientActionCommand = &cobra.Command{
	Use:   "action [site ID] [client ID] [action]",
	Short: "Execute an action on a specific client",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()

		action := &types.ClientActionRequest{
			Action: args[2],
		}

		err := c.Network.ClientExecuteAction(
			types.SiteID(args[0]), types.ClientID(args[1]), action)
		if err != nil {
			log.Error(err.Error())
			return
		}
		log.Info("Request success: 200 OK")
	},
}

var listVouchersCmd = &cobra.Command{
	Use:   "list [site ID]",
	Short: "List hotspot vouchers of a site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		// FIXME: Deal with page options!!
		vouchers, err := c.Network.Vouchers(types.SiteID(args[0]),
			types.Filter(filter), &types.PageArguments{})
		if err != nil {
			log.Error(err.Error())
			return
		}
		for _, voucher := range vouchers {
			if idOnly {
				fmt.Println(voucher.ID)
			} else {
				err := MarshalAndPrintJSON(voucher)
				if err != nil {
					log.Error(err.Error())
					return
				}
			}
		}
	},
}

var voucherDetailsCmd = &cobra.Command{
	Use:   "details [site ID] [voucher ID]",
	Short: "Get detailed information about a specific voucher",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		voucher, err := c.Network.VoucherDetails(types.SiteID(args[0]), types.VoucherID(args[1]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = MarshalAndPrintJSON(voucher)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}
