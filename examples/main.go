package main

import "strings"
import "time"
import "os"
import "context"

import "github.com/ClifHouck/unified/client"
import "github.com/ClifHouck/unified/types"

import "github.com/gopxl/beep"
import "github.com/gopxl/beep/effects"
import "github.com/gopxl/beep/mp3"
import "github.com/gopxl/beep/speaker"

import log "github.com/sirupsen/logrus"

func main() {
	data, err := os.ReadFile("unifi_api.key")
	if err != nil {
		log.Error(err.Error())
		return
	}

	apiKey := strings.TrimSpace(string(data))

	config := client.NewDefaultConfig(apiKey)

	valid, reasons := config.IsValid()
	if !valid {
		log.Error("Unifi client configuration is invalid, aborting:")
		for _, reason := range reasons {
			log.Error(reason)
		}
		return
	}
	unifiClient := client.NewClient(config)

	info, err := unifiClient.ProtectInfo()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("ProtectInfo encountered error")
		return
	}
	log.WithFields(log.Fields{
		"version": info.ApplicationVersion,
	}).Info("Unifi Protect Info")

	ctx, cancel := context.WithCancel(context.Background())

	eventChan, err := unifiClient.SubscribeProtectEvents(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("SubscribeProtectEvent encountered error")
		cancel()
		return
	}
	defer cancel()

	stream, format, err := LoadMP3("./Ding-dong-sound-effect.mp3")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("LoadMP3 encountered error")
		return
	}

	streamHandler := client.NewProtectEventStreamHandler(ctx, eventChan)

	streamHandler.SetRingEventHandler(func(eventType string, event *types.RingEvent) {
		if eventType == "add" {
			PlayDingDong(stream, format)
		}
	})

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
		log.Fatal(err)
	}

	return streamer, &format, nil
}

func PlayDingDong(streamer beep.StreamSeekCloser, format *beep.Format) {
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Errorf("PlayDingDong: speaker.Init failed with error: %s", err.Error())
		return
	}

	err = streamer.Seek(0)
	if err != nil {
		log.Errorf("PlayDingDong: streamer.Seek failed with error: %s", err.Error())
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
	log.Info("Ding dong done playing.")
}
