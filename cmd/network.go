package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ClifHouck/unified/types"
)

var filter = ""
var filterFlagSet = pflag.NewFlagSet("filter", pflag.ExitOnError)

var hidePage = false
var pageArgs = &types.PageArguments{}
var pageFlagSet = pflag.NewFlagSet("page", pflag.ExitOnError)

var voucherGenerateReq = &types.VoucherGenerateRequest{}

func init() {
	filterFlagSet.StringVar(&filter, "filter", "", "Filter results based on expression")

	listingFlagSet.BoolVar(&idOnly, "id-only", false, "List only the ID of listed entities, one per line.")

	pageFlagSet.BoolVar(&hidePage, "hide-page", false, "Hides the returned current page information")
	pageFlagSet.Uint32Var(&pageArgs.Offset, "page-offset", 0, "Offset of page to request")
	pageFlagSet.Uint32Var(&pageArgs.Limit, "page-limit", 0, "Limit of items per page")

	networkCmd.AddCommand(networkInfoCmd)
	networkCmd.AddCommand(devicesCmd)
	networkCmd.AddCommand(sitesCmd)
	networkCmd.AddCommand(clientsCmd)
	networkCmd.AddCommand(vouchersCmd)

	// Sites
	listSitesCmd.Flags().AddFlagSet(listingFlagSet)
	listSitesCmd.Flags().AddFlagSet(filterFlagSet)
	listSitesCmd.Flags().AddFlagSet(pageFlagSet)
	sitesCmd.AddCommand(listSitesCmd)

	// Clients
	listClientsCmd.Flags().AddFlagSet(listingFlagSet)
	listClientsCmd.Flags().AddFlagSet(filterFlagSet)
	listClientsCmd.Flags().AddFlagSet(pageFlagSet)
	clientsCmd.AddCommand(listClientsCmd)
	clientsCmd.AddCommand(clientDetailsCmd)
	clientsCmd.AddCommand(actionClientCmd)

	// Devices
	listDevicesCmd.Flags().AddFlagSet(listingFlagSet)
	listDevicesCmd.Flags().AddFlagSet(pageFlagSet)
	devicesCmd.AddCommand(listDevicesCmd)
	devicesCmd.AddCommand(deviceDetailsCmd)
	devicesCmd.AddCommand(statsDevicesCmd)
	devicesCmd.AddCommand(actionDevicesCmd)
	devicesCmd.AddCommand(actionDevicePortCmd)

	// Vouchers
	listVouchersCmd.Flags().AddFlagSet(listingFlagSet)
	listVouchersCmd.Flags().AddFlagSet(filterFlagSet)
	listVouchersCmd.Flags().AddFlagSet(pageFlagSet)
	vouchersCmd.AddCommand(listVouchersCmd)
	vouchersCmd.AddCommand(voucherDetailsCmd)

	voucherGenerateCmd.Flags().IntVar(&voucherGenerateReq.Count, "count", 0, "Number of vouchers")
	voucherGenerateCmd.Flags().StringVar(&voucherGenerateReq.Name, "name", "", "Name of vouchers")
	voucherGenerateCmd.Flags().IntVar(&voucherGenerateReq.AuthorizedGuestLimit, "guest-limit", 0, "Authorized guest limit")
	voucherGenerateCmd.Flags().IntVar(&voucherGenerateReq.TimeLimitMinutes, "time-limit", 0, "Time limit in minutes")
	voucherGenerateCmd.Flags().IntVar(&voucherGenerateReq.DataUsageLimitMBytes, "data-limit", 0, "Data limit in megabytes")
	voucherGenerateCmd.Flags().IntVar(&voucherGenerateReq.RxRateLimitKbps, "rx-limit", 0, "Recieve rate limit in kilobytes")
	voucherGenerateCmd.Flags().IntVar(&voucherGenerateReq.TxRateLimitKbps, "tx-limit", 0, "Transmit rate limit in kilobytes")
	vouchersCmd.AddCommand(voucherGenerateCmd)

	vouchersCmd.AddCommand(voucherDeleteCmd)

	voucherDeleteByFilterCmd.Flags().AddFlagSet(filterFlagSet)
	vouchersCmd.AddCommand(voucherDeleteByFilterCmd)
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
		devices, page, err := c.Network.Devices(types.SiteID(args[0]), pageArgs)
		if err != nil {
			log.Error(err.Error())
			return
		}
		if idOnly {
			for _, device := range devices {
				if idOnly {
					fmt.Println(device.ID)
				}
			}
		} else {
			if !hidePage {
				err := MarshalAndPrintJSON(page)
				if err != nil {
					log.Error(err.Error())
					return
				}
			}

			err := MarshalAndPrintJSON(devices)
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

var actionDevicePortCmd = &cobra.Command{
	Use:   "action [site ID] [device ID] [portIdx] [action]",
	Short: "Execute an action on a specific adopted device",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()

		action := &types.DevicePortActionRequest{
			Action: args[3],
		}

		port, err := strconv.Atoi(args[2])
		if err != nil {
			log.Error(err.Error())
			return
		}

		err = c.Network.DevicePortExecuteAction(types.SiteID(args[0]),
			types.DeviceID(args[1]),
			types.PortIdx(port),
			action)
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
		sites, page, err := c.Network.Sites(types.Filter(filter), pageArgs)
		if err != nil {
			log.Error(err.Error())
			return
		}

		if idOnly {
			for _, site := range sites {
				fmt.Println(site.ID)
			}
		} else {
			if !hidePage {
				err := MarshalAndPrintJSON(page)
				if err != nil {
					log.Error(err.Error())
					return
				}
			}

			err := MarshalAndPrintJSON(sites)
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
		clients, page, err := c.Network.Clients(
			types.SiteID(args[0]),
			types.Filter(filter),
			pageArgs)
		if err != nil {
			log.Error(err.Error())
			return
		}
		if idOnly {
			for _, client := range clients {
				fmt.Println(client.ID)
			}
		} else {
			if !hidePage {
				err := MarshalAndPrintJSON(page)
				if err != nil {
					log.Error(err.Error())
					return
				}
			}

			err := MarshalAndPrintJSON(clients)
			if err != nil {
				log.Error(err.Error())
				return
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

var actionClientCmd = &cobra.Command{
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
		vouchers, page, err := c.Network.Vouchers(types.SiteID(args[0]),
			types.Filter(filter), pageArgs)
		if err != nil {
			log.Error(err.Error())
			return
		}
		if idOnly {
			for _, voucher := range vouchers {
				fmt.Println(voucher.ID)
			}
		} else {
			if !hidePage {
				err := MarshalAndPrintJSON(page)
				if err != nil {
					log.Error(err.Error())
					return
				}
			}

			err := MarshalAndPrintJSON(vouchers)
			if err != nil {
				log.Error(err.Error())
				return
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

var voucherGenerateCmd = &cobra.Command{
	Use:   "generate [site ID]",
	Short: "Generate one or more hotspot vouchers for a site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		vouchers, err := c.Network.VoucherGenerate(types.SiteID(args[0]), voucherGenerateReq)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = MarshalAndPrintJSON(vouchers)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var voucherDeleteCmd = &cobra.Command{
	Use:   "delete [site ID] [voucher ID]",
	Short: "Delete a specific voucher",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		voucherDeleteResp, err := c.Network.VoucherDelete(types.SiteID(args[0]), types.VoucherID(args[1]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = MarshalAndPrintJSON(voucherDeleteResp)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var voucherDeleteByFilterCmd = &cobra.Command{
	Use:   "delete-filter [site ID]",
	Short: "Delete many vouchers by way of filter - BE CAREFUL!",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()

		// TODO: Should there be a "--allow-empty-filter" flag or similar?
		if len(filter) == 0 {
			log.Error("Filter may not be empty for delete request")
			return
		}

		voucherDeleteResp, err := c.Network.VoucherDeleteByFilter(types.SiteID(args[0]), types.Filter(filter))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = MarshalAndPrintJSON(voucherDeleteResp)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}
