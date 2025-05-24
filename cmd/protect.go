package cmd

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ClifHouck/unified/types"
)

var viewerSettingsReq = &types.ViewerSettingsRequest{
	Name: "string",
}

var (
	highQuality    bool
	mediumQuality  bool
	lowQuality     bool
	packageQuality bool

	snapshotLowQuality  bool
	snapshotJPEGQuality int
)

var qualitiesFlagSet = pflag.NewFlagSet("qualities", pflag.ExitOnError)

func init() { //nolint:funlen
	// Top level commands
	protectCmd.AddCommand(protectInfoCmd)
	protectCmd.AddCommand(camerasCmd)
	protectCmd.AddCommand(subscribeCmd)
	protectCmd.AddCommand(viewersCmd)
	protectCmd.AddCommand(liveViewsCmd)
	protectCmd.AddCommand(lightsCmd)
	protectCmd.AddCommand(nvrCmd)
	protectCmd.AddCommand(chimesCmd)
	protectCmd.AddCommand(sensorsCmd)
	protectCmd.AddCommand(filesCmd)

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
	camerasCmd.AddCommand(cameraPatchCmd)

	cameraGetSnapshotCmd.Flags().BoolVar(&snapshotLowQuality, "low-quality", false, "snapshot low quality")
	cameraGetSnapshotCmd.Flags().IntVar(&snapshotJPEGQuality, "jpeg-quality", 100, "JPEG Quality from 1 to 100")
	camerasCmd.AddCommand(cameraGetSnapshotCmd)

	// TODO: Should RTSPS Streams related commands be a sub-command?
	qualitiesFlagSet.BoolVar(&highQuality, "high", false, "high stream quality")
	qualitiesFlagSet.BoolVar(&mediumQuality, "medium", false, "medium stream quality")
	qualitiesFlagSet.BoolVar(&lowQuality, "low", false, "low stream quality")
	qualitiesFlagSet.BoolVar(&packageQuality, "package", false, "package stream quality")

	cameraRTSPSStreamCreateCmd.Flags().AddFlagSet(qualitiesFlagSet)
	camerasCmd.AddCommand(cameraRTSPSStreamCreateCmd)
	cameraRTSPSStreamDeleteCmd.Flags().AddFlagSet(qualitiesFlagSet)
	camerasCmd.AddCommand(cameraRTSPSStreamDeleteCmd)
	camerasCmd.AddCommand(cameraRTSPSStreamGetCmd)
	camerasCmd.AddCommand(cameraDisableMicPermanentlyCmd)
	camerasCmd.AddCommand(cameraTalkbackSessionCmd)

	// Lights
	lightListCmd.Flags().AddFlagSet(listingFlagSet)
	lightsCmd.AddCommand(lightListCmd)
	lightsCmd.AddCommand(lightDetailsCmd)
	lightsCmd.AddCommand(lightPatchCmd)

	// Chimes
	chimeListCmd.Flags().AddFlagSet(listingFlagSet)
	chimesCmd.AddCommand(chimeListCmd)
	chimesCmd.AddCommand(chimeDetailsCmd)
	chimesCmd.AddCommand(chimePatchCmd)

	// Sensors
	sensorListCmd.Flags().AddFlagSet(listingFlagSet)
	sensorsCmd.AddCommand(sensorListCmd)
	sensorsCmd.AddCommand(sensorDetailsCmd)
	sensorsCmd.AddCommand(sensorPatchCmd)

	// Device Asset File Management
	filesListCmd.Flags().AddFlagSet(listingFlagSet)
	filesCmd.AddCommand(filesListCmd)
	filesCmd.AddCommand(fileUploadCmd)
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

var lightsCmd = &cobra.Command{
	Use:   "lights",
	Short: "Make UniFi Protect `lights` calls",
	Long:  `Call camera endpoints under UniFi Protect's API.`,
}

var chimesCmd = &cobra.Command{
	Use:   "chimes",
	Short: "Make UniFi Protect `chimes` calls",
	Long:  `Call chimes endpoints under UniFi Protect's API.`,
}

var sensorsCmd = &cobra.Command{
	Use:   "sensors",
	Short: "Make UniFi Protect `sensors` calls",
	Long:  `Call sensors endpoints under UniFi Protect's API.`,
}

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Make UniFi Protect device asset `files` calls",
	Long:  `Call device asset files endpoints under UniFi Protect's API.`,
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

var cameraPatchCmd = &cobra.Command{
	Use:   "patch [camera ID] [camera JSON filename]",
	Short: "Patch the configuration of an existing camera",
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		data, err := os.ReadFile(args[1])
		if err != nil {
			log.Error(err.Error())
			return
		}

		var cameraReq types.CameraPatchRequest
		err = json.Unmarshal(data, &cameraReq)
		if err != nil {
			log.Error(err.Error())
			return
		}

		c := getClient()
		modifiedCamera, err := c.Protect.CameraPatch(types.CameraID(args[0]), &cameraReq)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(modifiedCamera)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var cameraGetSnapshotCmd = &cobra.Command{
	Use:   "snapshot [camera ID] [filename]",
	Short: "Get a live snapshot image from a specified camera and save it to a file",
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		if snapshotJPEGQuality < 1 || snapshotJPEGQuality > 100 {
			log.Errorf("--jpeg-quality must be between 1 and 100, got '%d'", snapshotJPEGQuality)
			return
		}

		c := getClient()
		image, err := c.Protect.CameraGetSnapshot(types.CameraID(args[0]), !snapshotLowQuality)
		if err != nil {
			log.Error(err.Error())
			return
		}

		outfile, err := os.Create(args[1])
		if err != nil {
			log.Error(err.Error())
			return
		}

		err = jpeg.Encode(outfile, image, &jpeg.Options{Quality: snapshotJPEGQuality})
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.WithFields(logrus.Fields{
			"filename": args[1],
		}).Infof("Saved snapshot to file")
	},
}

func qualities() []string {
	quals := make([]string, 0, 4)
	if lowQuality {
		quals = append(quals, "low")
	}
	if mediumQuality {
		quals = append(quals, "medium")
	}
	if highQuality {
		quals = append(quals, "high")
	}
	if packageQuality {
		quals = append(quals, "package")
	}
	return quals
}

var cameraRTSPSStreamCreateCmd = &cobra.Command{
	Use:   "stream-create [camera ID]",
	Short: "Create RTSPS stream(s), based on qualities specified, for a camera",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		resp, err := c.Protect.CameraCreateRTSPSStream(
			types.CameraID(args[0]),
			&types.CameraCreateRTSPSStreamRequest{Qualities: qualities()},
		)
		if err != nil {
			log.Error(err.Error())
			return
		}

		err = marshalAndPrintJSON(resp)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var cameraRTSPSStreamDeleteCmd = &cobra.Command{
	Use:   "stream-delete [camera ID]",
	Short: "Delete RTSPS stream(s), based on qualities specified, for a camera",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		err := c.Protect.CameraDeleteRTSPSStream(types.CameraID(args[0]), &types.CameraDeleteRTSPSStreamRequest{
			Qualities: qualities(),
		})
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.Info("204 - No Content - OK")
	},
}

var cameraRTSPSStreamGetCmd = &cobra.Command{
	Use:   "stream-get [camera ID]",
	Short: "Get RTSPS streams that exist for a camera",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		resp, err := c.Protect.CameraGetRTSPSStream(types.CameraID(args[0]))
		if err != nil {
			log.Error(err.Error())
			return
		}

		err = marshalAndPrintJSON(resp)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var cameraDisableMicPermanentlyCmd = &cobra.Command{
	Use:   "disable-mic-permanently [camera ID]",
	Short: "Permanently disable the microphone for a specific camera",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		camera, err := c.Protect.CameraDisableMicPermanently(types.CameraID(args[0]))
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

var cameraTalkbackSessionCmd = &cobra.Command{
	Use:   "talkback-session [camera ID]",
	Short: "Get the talkback stream URL and audio config for a camera",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		cameraTalkbackResp, err := c.Protect.CameraTalkbackSession(types.CameraID(args[0]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(cameraTalkbackResp)
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

var lightListCmd = &cobra.Command{
	Use:   "list",
	Short: "List adopted Protect lights",
	Run: func(_ *cobra.Command, _ []string) {
		c := getClient()
		lights, err := c.Protect.Lights()
		if err != nil {
			log.Error(err.Error())
			return
		}
		if idOnly {
			for _, light := range lights {
				if idOnly {
					fmt.Println(light.ID)
				}
			}
		} else {
			err = marshalAndPrintJSON(lights)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	},
}

var lightDetailsCmd = &cobra.Command{
	Use:   "details [light ID]",
	Short: "Get detailed information about a specific adopted device",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		light, err := c.Protect.LightDetails(types.LightID(args[0]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(light)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var lightPatchCmd = &cobra.Command{
	Use:   "patch [light ID] [light JSON filename]",
	Short: "Patch the configuration of an existing light",
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		data, err := os.ReadFile(args[1])
		if err != nil {
			log.Error(err.Error())
			return
		}

		var lightReq types.LightPatchRequest
		err = json.Unmarshal(data, &lightReq)
		if err != nil {
			log.Error(err.Error())
			return
		}

		c := getClient()
		modifiedlight, err := c.Protect.LightPatch(types.LightID(args[0]), &lightReq)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(modifiedlight)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var nvrCmd = &cobra.Command{
	Use:   "nvrs",
	Short: "Get information about the NVR",
	Run: func(_ *cobra.Command, _ []string) {
		c := getClient()
		nvr, err := c.Protect.NVRs()
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(nvr)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var chimeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List adopted Protect chimes",
	Run: func(_ *cobra.Command, _ []string) {
		c := getClient()
		chimes, err := c.Protect.Chimes()
		if err != nil {
			log.Error(err.Error())
			return
		}
		if idOnly {
			for _, chime := range chimes {
				if idOnly {
					fmt.Println(chime.ID)
				}
			}
		} else {
			err = marshalAndPrintJSON(chimes)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	},
}

var chimeDetailsCmd = &cobra.Command{
	Use:   "details [chime ID]",
	Short: "Get detailed information about a specific adopted device",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		chime, err := c.Protect.ChimeDetails(types.ChimeID(args[0]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(chime)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var chimePatchCmd = &cobra.Command{
	Use:   "patch [chime ID] [chime JSON filename]",
	Short: "Patch the configuration of an existing chime",
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		data, err := os.ReadFile(args[1])
		if err != nil {
			log.Error(err.Error())
			return
		}

		var chimeReq types.ChimePatchRequest
		err = json.Unmarshal(data, &chimeReq)
		if err != nil {
			log.Error(err.Error())
			return
		}

		c := getClient()
		modifiedchime, err := c.Protect.ChimePatch(types.ChimeID(args[0]), &chimeReq)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(modifiedchime)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var sensorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List adopted Protect sensors",
	Run: func(_ *cobra.Command, _ []string) {
		c := getClient()
		sensors, err := c.Protect.Sensors()
		if err != nil {
			log.Error(err.Error())
			return
		}
		if idOnly {
			for _, sensor := range sensors {
				if idOnly {
					fmt.Println(sensor.ID)
				}
			}
		} else {
			err = marshalAndPrintJSON(sensors)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	},
}

var sensorDetailsCmd = &cobra.Command{
	Use:   "details [sensor ID]",
	Short: "Get detailed information about a specific adopted device",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		c := getClient()
		sensor, err := c.Protect.SensorDetails(types.SensorID(args[0]))
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(sensor)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var sensorPatchCmd = &cobra.Command{
	Use:   "patch [sensor ID] [sensor JSON filename]",
	Short: "Patch the configuration of an existing sensor",
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		data, err := os.ReadFile(args[1])
		if err != nil {
			log.Error(err.Error())
			return
		}

		var sensorReq types.SensorPatchRequest
		err = json.Unmarshal(data, &sensorReq)
		if err != nil {
			log.Error(err.Error())
			return
		}

		c := getClient()
		modifiedsensor, err := c.Protect.SensorPatch(types.SensorID(args[0]), &sensorReq)
		if err != nil {
			log.Error(err.Error())
			return
		}
		err = marshalAndPrintJSON(modifiedsensor)
		if err != nil {
			log.Error(err.Error())
			return
		}
	},
}

var filesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Protect device asset files",
	Run: func(_ *cobra.Command, _ []string) {
		c := getClient()
		files, err := c.Protect.Files(types.FileTypeAnimations)
		if err != nil {
			log.Error(err.Error())
			return
		}
		if idOnly {
			for _, file := range files {
				if idOnly {
					fmt.Println(file.Name)
				}
			}
		} else {
			err = marshalAndPrintJSON(files)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	},
}

var fileUploadCmd = &cobra.Command{
	Use:   "upload [filename]",
	Short: "Upload a device asset file",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		data, err := os.ReadFile(args[0])
		if err != nil {
			log.Error(err.Error())
			return
		}

		_, filename := filepath.Split(args[0])

		c := getClient()
		err = c.Protect.FileUpload(types.FileTypeAnimations, filename, data)
		if err != nil {
			log.Error(err.Error())
		}
	},
}
