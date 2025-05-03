package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ClifHouck/unified/types"
)

func init() {
	protectCmd.AddCommand(protectInfoCmd)
	protectCmd.AddCommand(camerasCmd)
	protectCmd.AddCommand(subscribeCmd)

	// Subscriptions
	subscribeCmd.AddCommand(deviceEventsCmd)
	subscribeCmd.AddCommand(protectEventsCmd)

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
	Run: func(_ *cobra.Command, _ []string) {
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

				err = MarshalAndPrintJSON(item)
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

				err = MarshalAndPrintJSON(item)
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
			err = MarshalAndPrintJSON(cameras)
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
		err = MarshalAndPrintJSON(camera)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}
