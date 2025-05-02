package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ClifHouck/unified/types"
)

func init() {
	protectCmd.AddCommand(protectInfoCmd)
	protectCmd.AddCommand(camerasCmd)
	protectCmd.AddCommand(subscribeCmd)

	// Subscriptions
	// subscribeCmd.AddCommand(deviceEventsCmd)
	// subscribeCmd.AddComannd(protectEventsCmd)

	// Cameras
	cameraListCmd.Flags().AddFlagSet(listingFlagSet)
	camerasCmd.AddCommand(cameraListCmd)
	camerasCmd.AddCommand(cameraDetailsCmd)
}

var protectCmd = &cobra.Command{
	Use:   "protect",
	Short: "Make UniFi Protect API calls",
	Long:  `Complete access to UniFi's Protect API from the command line`,
}

var camerasCmd = &cobra.Command{
	Use:   "cameras",
	Short: "Make UniFi Protect `cameras` calls",
	Long:  `Call camera endpoints under UniFi Protect's API.`,
}

var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Make UniFi Protect `subscribe` calls",
	Long:  `Call subscribe endpoints under UniFi Protect's API.`,
}

var protectInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get protect application info",
	Long:  `Get generic information about the Protect application`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		info, err := c.Protect.Info()
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

var cameraListCmd = &cobra.Command{
	Use:   "list",
	Short: "List adopted Protect cameras",
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		cameras, err := c.Protect.Cameras()
		if err != nil {
			log.Error(err.Error())
			return
		}
		if idOnly {
			for _, camera := range cameras {
				if idOnly {
					fmt.Println(camera.ID)
				}
			}
		} else {
			err := MarshalAndPrintJSON(cameras)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	},
}

var cameraDetailsCmd = &cobra.Command{
	Use:   "details [camera ID]",
	Short: "Get detailed information about a specific adopted device",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		camera, err := c.Protect.CameraDetails(types.CameraID(args[0]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = MarshalAndPrintJSON(camera)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}
