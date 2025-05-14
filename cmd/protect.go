package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ClifHouck/unified/types"
)

var viewerSettingsReq = &types.ViewerSettingsRequest{
	Name: "string",
}

func init() {
	protectCmd.AddCommand(protectInfoCmd)
	protectCmd.AddCommand(camerasCmd)
	protectCmd.AddCommand(subscribeCmd)
	protectCmd.AddCommand(viewersCmd)
	protectCmd.AddCommand(liveViewsCmd)

	// Subscriptions
	subscribeCmd.AddCommand(deviceEventsCmd)
	subscribeCmd.AddCommand(protectEventsCmd)

	// Viewers
	viewerListCmd.Flags().AddFlagSet(listingFlagSet)
	viewersCmd.AddCommand(viewerListCmd)
	viewersCmd.AddCommand(viewerDetailsCmd)
	viewerSettingsCmd.Flags().StringVar(&viewerSettingsReq.Liveview, "liveview", "", "A live view ID to set")
	viewersCmd.AddCommand(viewerSettingsCmd)

	// Live Views
	liveViewListCmd.Flags().AddFlagSet(listingFlagSet)
	liveViewsCmd.AddCommand(liveViewListCmd)
	liveViewsCmd.AddCommand(liveViewDetailsCmd)
	liveViewsCmd.AddCommand(liveViewCreateCmd)
	liveViewsCmd.AddCommand(liveViewPatchCmd)

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

var viewersCmd = &cobra.Command{
	Use:   "viewers",
	Short: "Make UniFi Protect `viewers` calls",
	Long:  `Call viewer endpoints under UniFi Protect's API.`,
}

var liveViewsCmd = &cobra.Command{
	Use:   "liveviews",
	Short: "Make UniFi Protect `liveviews` calls",
	Long:  `Call liveview endpoints under UniFi Protect's API.`,
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
	Run: func(_ *cobra.Command, _ []string) {
		c := getClient()
		info, err := c.Protect.Info()
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(info)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var deviceEventsCmd = &cobra.Command{
	Use:   "device-events",
	Short: "Stream device events from Protect API",
	Run: func(_ *cobra.Command, _ []string) {
		c := getClient()
		events, err := c.Protect.SubscribeDeviceEvents()
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.Info("Streaming device events...")
		for {
			select {
			case streamEvent := <-events:
				if streamEvent == nil {
					log.Warn("Got nil event. Bailing out!")
					return
				}

				var item types.ProtectDeviceEventItem
				err = json.Unmarshal(streamEvent.RawItem, &item)
				if err != nil {
					log.Error("Couldn't parse RawItem!")
					log.Error(err.Error())
				}

				log.WithFields(logrus.Fields{
					"ID":           item.ID,
					"event.type":   streamEvent.ItemType,
					"message.type": streamEvent.Type,
				}).Info("Received ProtectDeviceEvent")

				err = marshalAndPrintJSON(item)
				if err != nil {
					log.Error(err.Error())
					return
				}

			case <-ctx.Done():
				log.Warn("Got context.Done!")
				return
			}
		}
	},
}

var protectEventsCmd = &cobra.Command{
	Use:   "protect-events",
	Short: "Stream protect events from Protect API",
	Run: func(_ *cobra.Command, _ []string) {
		c := getClient()
		events, err := c.Protect.SubscribeProtectEvents()
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.Info("Streaming protect events...")
		for {
			select {
			case streamEvent := <-events:
				if streamEvent == nil {
					log.Warn("Got nil event. Bailing out!")
					return
				}

				var item types.ProtectEventItem
				err = json.Unmarshal(streamEvent.RawItem, &item)
				if err != nil {
					log.Error("Couldn't parse RawItem!")
					log.Error(err.Error())
				}

				log.WithFields(logrus.Fields{
					"ID":           item.ID,
					"event.type":   streamEvent.ItemType,
					"message.type": streamEvent.Type,
				}).Info("Received ProtectEvent")

				err = marshalAndPrintJSON(item)
				if err != nil {
					log.Error(err.Error())
					return
				}

			case <-ctx.Done():
				log.Warn("Got context.Done!")
				return
			}
		}
	},
}

var cameraListCmd = &cobra.Command{
	Use:   "list",
	Short: "List adopted Protect cameras",
	Run: func(_ *cobra.Command, _ []string) {
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
			err = marshalAndPrintJSON(cameras)
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
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		camera, err := c.Protect.CameraDetails(types.CameraID(args[0]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(camera)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var viewerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all viewers",
	Run: func(_ *cobra.Command, _ []string) {
		c := getClient()
		viewers, err := c.Protect.Viewers()
		if err != nil {
			log.Error(err.Error())
			return
		}
		if idOnly {
			for _, viewer := range viewers {
				if idOnly {
					fmt.Println(viewer.ID)
				}
			}
		} else {
			err = marshalAndPrintJSON(viewers)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	},
}

var viewerDetailsCmd = &cobra.Command{
	Use:   "details [viewer ID]",
	Short: "Get detailed information about a viewer",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		viewer, err := c.Protect.ViewerDetails(types.ViewerID(args[0]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(viewer)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var viewerSettingsCmd = &cobra.Command{
	Use:   "settings [viewer ID]",
	Short: "Patch the settings for a specific viewer",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		viewer, err := c.Protect.ViewerSettings(types.ViewerID(args[0]), viewerSettingsReq)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(viewer)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var liveViewListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all liveviews",
	Run: func(_ *cobra.Command, _ []string) {
		c := getClient()
		liveViews, err := c.Protect.LiveViews()
		if err != nil {
			log.Error(err.Error())
			return
		}
		if idOnly {
			for _, liveView := range liveViews {
				if idOnly {
					fmt.Println(liveView.ID)
				}
			}
		} else {
			err = marshalAndPrintJSON(liveViews)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	},
}

var liveViewDetailsCmd = &cobra.Command{
	Use:   "details [liveview ID]",
	Short: "Get detailed information about a liveview",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		liveView, err := c.Protect.LiveViewDetails(types.LiveViewID(args[0]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(liveView)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var liveViewCreateCmd = &cobra.Command{
	Use:   "create [liveview JSON filename]",
	Short: "Create a new liveview",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		data, err := os.ReadFile(args[0])
		if err != nil {
			log.Error(err.Error())
			return
		}

		var liveView types.LiveView
		err = json.Unmarshal(data, &liveView)
		if err != nil {
			log.Error(err.Error())
			return
		}

		c := getClient()
		newLiveView, err := c.Protect.LiveViewCreate(&liveView)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(newLiveView)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var liveViewPatchCmd = &cobra.Command{
	Use:   "patch [liveview ID] [liveview JSON filename]",
	Short: "Patch the configuration of an existing liveview",
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		data, err := os.ReadFile(args[1])
		if err != nil {
			log.Error(err.Error())
			return
		}

		var liveView types.LiveView
		err = json.Unmarshal(data, &liveView)
		if err != nil {
			log.Error(err.Error())
			return
		}

		c := getClient()
		modifiedLiveView, err := c.Protect.LiveViewPatch(types.LiveViewID(args[0]), &liveView)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(modifiedLiveView)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}
