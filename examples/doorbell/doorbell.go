package main

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ClifHouck/unified/client"
	"github.com/ClifHouck/unified/types"
)

var log *logrus.Logger

var mp3Filename string
var apiKeyFilename string

var rootCmd = &cobra.Command{
	Use:   "doorbell",
	Short: "Watch for ring events and play an MP3 when they occur",
	Run: func(cmd *cobra.Command, args []string) {
		Doorbell()
	},
}

func init() {
	rootCmd.Flags().StringVar(&mp3Filename, "mp3", "", "Filename of MP3 to load for doorbell sound")
	rootCmd.Flags().StringVar(&apiKeyFilename, "api-key", "", "File containing UniFi API key")
}

func main() {
	rootCmd.Execute()
}

func Doorbell() {
	log = logrus.New()

	data, err := os.ReadFile(apiKeyFilename)
	if err != nil {
		log.Error(err.Error())
		return
	}

	apiKey := strings.TrimSpace(string(data))

	config := client.NewDefaultConfig(apiKey)

	valid, reasons := config.IsValid()
	if !valid {
		log.Error("Please fix the following unified client configuration problems:")
		for _, reason := range reasons {
			log.Error(reason)
		}
		log.Fatal("Unifi client configuration is invalid, aborting.")
	}

	ctx, cancel := context.WithCancel(context.Background())
	unifiClient := client.NewClient(ctx, config, log)

	info, err := unifiClient.Protect.Info()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("ProtectInfo encountered error")
		return
	}
	log.WithFields(logrus.Fields{
		"version": info.ApplicationVersion,
	}).Info("Unifi Protect Info")

	eventChan, err := unifiClient.Protect.SubscribeProtectEvents()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("SubscribeProtectEvent encountered error")
		cancel()
		return
	}
	defer cancel()

	stream, format, err := LoadMP3(mp3Filename)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":    err.Error(),
			"filename": mp3Filename,
		}).Fatal("LoadMP3 encountered error")
	}
	log.WithFields(logrus.Fields{
		"filename": mp3Filename,
	}).Info("Load MP3 success")

	streamHandler := client.NewProtectEventStreamHandler(ctx, eventChan)

	// Sync this because the event handler will be called asynchronously.
	var handlerMutex sync.Mutex
	streamHandler.SetRingEventHandler(func(eventType string, _ *types.RingEvent) {
		handlerMutex.Lock()
		defer handlerMutex.Unlock()
		if eventType == "add" {
			PlayMP3(stream, format) // This is where the magic happens!
		}
	})

	go streamHandler.Process()

	<-ctx.Done()
	log.Warn("Got context.Done!")
}

func LoadMP3(filename string) (beep.StreamSeekCloser, *beep.Format, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return nil, nil, err
	}

	return streamer, &format, nil
}

func PlayMP3(streamer beep.StreamSeekCloser, format *beep.Format) {
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Errorf("PlayMP3: speaker.Init failed with error: %s", err.Error())
		return
	}

	err = streamer.Seek(0)
	if err != nil {
		log.Errorf("PlayMP3: streamer.Seek failed with error: %s", err.Error())
		return
	}

	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   -0.5,
		Silent:   false,
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		done <- true
	})))
	<-done
	log.Info("MP3 done playing.")
}
