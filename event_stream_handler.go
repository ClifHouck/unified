package events

import "context"
import "encoding/json"
import "sync"

import log "github.com/sirupsen/logrus"

import "github.com/ClifHouck/unified/types"
import "github.com/ClifHouck/unified/client"

type ProtectEventStreamHandler struct {
	ctx    context.Context
	stream <-chan *client.ProtectEventMessage

	ringHandler func(string, *types.RingEvent)
	ringMutex   sync.Mutex

	sensorExtremeValuesHandler func(string, *types.SensorExtremeValuesEvent)
	sensorExtremeValuesMutex   sync.Mutex

	SensorWaterLeakHandler func(string, *types.SensorWaterLeakEvent)
	SensorWaterLeakMutex   sync.Mutex

	SensorTamperHandler func(string, *types.SensorTamperEvent)
	SensorTamperMutex   sync.Mutex

	SensorBatteryLowHandler func(string, *types.SensorBatteryLowEvent)
	SensorBatteryLowMutex   sync.Mutex
}

func NewProtectEventStreamHandler(ctx context.Context,
	stream <-chan *client.ProtectEventMessage) *ProtectEventStreamHandler {
	handler := &ProtectEventStreamHandler{
		ctx:    ctx,
		stream: stream,
	}

	// Should this be here, or should clients of this class call it
	// explicitly?
	go handler.processStream()

	return handler
}

func (pesh *ProtectEventStreamHandler) processStream() {
	log.Info("Waiting for events...")
	for {
		select {
		case message := <-pesh.stream:
			if message == nil {
				log.Warn("Got nil message. Bailing out!")
				return
			}

			var item types.ProtectEventItem
			err := json.Unmarshal(message.Event.RawItem, &item)
			if err != nil {
				log.Error("Couldn't parse RawItem!")
				log.Error(err.Error())
			}

			log.WithFields(log.Fields{
				"ID":           item.ID,
				"device":       item.Device,
				"event.type":   item.Type,
				"message.type": message.Event.Type,
			}).Info("Received ProtectEvent")

			switch event := message.Event.Item.(type) {
			case *types.RingEvent:
				go pesh.invokeRingEventHandler(message.Event.Type, event)
			case *types.SensorExtremeValuesEvent:
			case *types.SensorWaterLeakEvent:
			case *types.SensorTamperEvent:
			case *types.SensorBatteryLowEvent:
			case *types.SensorAlarmEvent:
			case *types.SensorOpenedEvent:
			case *types.SensorClosedEvent:
			case *types.LightMotionEvent:
			case *types.CameraMotionEvent:
			case *types.CameraSmartDetectAudioEvent:
			case *types.CameraSmartDetectZoneEvent:
			case *types.CameraSmartDetectLineEvent:
			case *types.CameraSmartDetectLoiterEvent:
			default:
				log.Error("Unknown type encountered: '%s'", message.Event.ItemType)
			}

			if message.Error != nil {
				log.Error(message.Error.Error())
				return
			}
		case <-pesh.ctx.Done():
			log.Warn("Got context.Done!")
			return
		}
	}
}

func (pesh *ProtectEventStreamHandler) SetRingEventHandler(handler func(string, *types.RingEvent)) {
	pesh.ringMutex.Lock()
	defer pesh.ringMutex.Unlock()

	pesh.ringHandler = handler
}

func (pesh *ProtectEventStreamHandler) invokeRingEventHandler(eventType string, event *types.RingEvent) {
	pesh.ringMutex.Lock()
	defer pesh.ringMutex.Unlock()

	if pesh.ringHandler != nil {
		go pesh.ringHandler(eventType, event)
	}
}
