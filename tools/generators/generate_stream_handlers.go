package main

import (
	"os"
	"reflect"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/ClifHouck/unified/types"
)

type StreamHandlerArguments struct {
	PackageName         string
	StreamType          string
	EventType           string
	EventTypeFirstLower string
	AllEventTypes       []interface{}
	Filename            string
}

const topOfFileComment = `
// (!) DO NOT EDIT (!) Generated by generate_stream_handlers
`

const streamHandlerPackage = "package {{.PackageName}}\n\n"

const imports = `import (
	"context"
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/ClifHouck/unified/types"
)
`

const streamHandlerStructBegin = `
type {{.StreamType}}StreamHandler struct {
	ctx    context.Context
	stream <-chan *types.{{.StreamType}}
`

const streamHandlerStructEventTypeMembers = `
	{{.EventTypeFirstLower}}Handler func(string, *types.{{.EventType}})
	{{.EventTypeFirstLower}}Mutex   sync.Mutex
`

const streamHandlerStructEnd = `} // {{.StreamType}}StreamHandler
`

const newStreamHandlerObjectFunction = `
func New{{.StreamType}}StreamHandler(ctx context.Context,
	stream <-chan *types.{{.StreamType}}) *{{.StreamType}}StreamHandler {
	handler := &{{.StreamType}}StreamHandler{
		ctx:    ctx,
		stream: stream,
	}

	return handler
}
`

const processStreamMethodBegin = `
func (esh *{{.StreamType}}StreamHandler) Process() {
	log.Info("Waiting for events...")
	for {
		select {
		case streamEvent := <-esh.stream:
			if streamEvent == nil {
				log.Warn("Got nil event. Bailing out!")
				return
			}

			var item types.{{.StreamType}}Item
			err := json.Unmarshal(streamEvent.RawItem, &item)
			if err != nil {
				log.Error("Couldn't parse RawItem!")
				log.Error(err.Error())
			}

			log.WithFields(log.Fields{
				"ID":           item.ID,
				"event.type":   streamEvent.ItemType,
				"message.type": streamEvent.Type,
			}).Info("Received {{.StreamType}}")

			switch event := streamEvent.Item.(type) {
`

const processStreamCase = `			case *types.{{.EventType}}:
				go esh.invoke{{.EventType}}Handler(streamEvent.Type, event)
`

const processStreamMethodEnd = `
			default:
				log.Errorf("Unknown type encountered: '%s'", streamEvent.ItemType)
			}

		case <-esh.ctx.Done():
			log.Warn("Got context.Done!")
			return
		}
	}
}
`

const setEventHandlerMethod = `
func (esh *{{.StreamType}}StreamHandler) Set{{.EventType}}Handler(handler func(string, *types.{{.EventType}})) {
	esh.{{.EventTypeFirstLower}}Mutex.Lock()
	defer esh.{{.EventTypeFirstLower}}Mutex.Unlock()

	esh.{{.EventTypeFirstLower}}Handler = handler
}
`

const invokeEventHandlerMethod = `
func (esh *{{.StreamType}}StreamHandler) invoke{{.EventType}}Handler(eventType string, event *types.{{.EventType}}) {
	esh.{{.EventTypeFirstLower}}Mutex.Lock()
	defer esh.{{.EventTypeFirstLower}}Mutex.Unlock()

	if esh.{{.EventTypeFirstLower}}Handler != nil {
		go esh.{{.EventTypeFirstLower}}Handler(eventType, event)
	}
}
`

var allTemplateDefinitions = map[string]string{
	"topOfFileComment":                    topOfFileComment,
	"streamHandlerPackage":                streamHandlerPackage,
	"imports":                             imports,
	"streamHandlerStructBegin":            streamHandlerStructBegin,
	"streamHandlerStructEventTypeMembers": streamHandlerStructEventTypeMembers,
	"streamHandlerStructEnd":              streamHandlerStructEnd,
	"newStreamHandlerObjectFunction":      newStreamHandlerObjectFunction,
	"processStreamMethodBegin":            processStreamMethodBegin,
	"processStreamCase":                   processStreamCase,
	"processStreamMethodEnd":              processStreamMethodEnd,
	"setEventHandlerMethod":               setEventHandlerMethod,
	"invokeEventHandlerMethod":            invokeEventHandlerMethod,
}

func renderStreamHandlerToFile(args *StreamHandlerArguments) error {
	outFile, err := os.OpenFile(args.Filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	templates := map[string]*template.Template{}
	for k, v := range allTemplateDefinitions {
		templates[k] = template.Must(template.New(k).Parse(v))
	}

	eventTypeNames := []string{}
	for _, eventObj := range args.AllEventTypes {
		eventTypeNames = append(eventTypeNames, reflect.TypeOf(eventObj).Name())
	}

	for _, templateName := range []string{
		"streamHandlerPackage",
		"imports",
		"topOfFileComment",
		"streamHandlerStructBegin",
	} {
		err = templates[templateName].Execute(outFile, args)
		if err != nil {
			return err
		}
	}

	for _, eventTypeName := range eventTypeNames {
		args.EventType = eventTypeName
		args.EventTypeFirstLower = strings.ToLower(eventTypeName[:1]) + eventTypeName[1:]
		err = templates["streamHandlerStructEventTypeMembers"].Execute(outFile, args)
		if err != nil {
			return err
		}
	}

	for _, templateName := range []string{
		"streamHandlerStructEnd",
		"newStreamHandlerObjectFunction",
		"processStreamMethodBegin",
	} {
		err = templates[templateName].Execute(outFile, args)
		if err != nil {
			return err
		}
	}

	for _, eventTypeName := range eventTypeNames {
		args.EventType = eventTypeName
		args.EventTypeFirstLower = strings.ToLower(eventTypeName[:1]) + eventTypeName[1:]
		err = templates["processStreamCase"].Execute(outFile, args)
		if err != nil {
			return err
		}
	}

	err = templates["processStreamMethodEnd"].Execute(outFile, args)
	if err != nil {
		return err
	}

	for _, eventTypeName := range eventTypeNames {
		args.EventType = eventTypeName
		args.EventTypeFirstLower = strings.ToLower(eventTypeName[:1]) + eventTypeName[1:]
		err = templates["setEventHandlerMethod"].Execute(outFile, args)
		if err != nil {
			return err
		}

		err = templates["invokeEventHandlerMethod"].Execute(outFile, args)
		if err != nil {
			return err
		}
	}

	log.WithFields(log.Fields{
		"filename":   args.Filename,
		"streamType": args.StreamType,
	}).Info("Stream Handler successfully rendered to file.")
	return nil
}

func main() {
	streamArgs := []*StreamHandlerArguments{
		{
			Filename:      "client/protect_device_update_stream_handler.go",
			PackageName:   "client",
			StreamType:    "ProtectDeviceEvent",
			AllEventTypes: types.AllProtectDeviceEvents,
		},
		{
			Filename:      "client/protect_event_stream_handler.go",
			PackageName:   "client",
			StreamType:    "ProtectEvent",
			AllEventTypes: types.AllProtectEvents,
		},
	}

	for _, streamArg := range streamArgs {
		err := renderStreamHandlerToFile(streamArg)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}
