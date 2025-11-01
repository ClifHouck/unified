package types_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ClifHouck/unified/types"
)

func rawTypeNameToJSONEventName(rawTypeName string) string {
		jsonTypeID := strings.ToLower(rawTypeName[0:1]) + rawTypeName[1:len(rawTypeName)-5]
		jsonTypeID = strings.ReplaceAll(jsonTypeID, "camera", "")
		jsonTypeID = strings.ToLower(jsonTypeID[0:1]) + jsonTypeID[1:]
		return jsonTypeID
}

func TestAllProtectEventTypesUnmarshalJSON(t *testing.T) {

	for _, eventObj := range types.AllProtectEvents {
		jsonEventName := rawTypeNameToJSONEventName(reflect.TypeOf(eventObj).Name())
		testCase := struct {
			json        	string
			jsonEventName 	string
			eventObj    	interface{}
		}{
			`{
			  "type": "add",
			  "item": {
				"id": "66d025b301ebc903e80003ea",
				"modelKey": "event",
				"start": 1445408038748,
				"end": 1445408048748,
				"device": "66d025b301ebc903e80003ea",
				"type": "%s"
			  }
			}`,
			jsonEventName,
			eventObj,
		}
		testCase.json = fmt.Sprintf(testCase.json, testCase.jsonEventName)
		t.Run(jsonEventName, func(t *testing.T) {
			var event types.ProtectEvent
			err := event.UnmarshalJSON([]byte(testCase.json))
			require.NoError(t, err)

			assert.Equal(t, "add", event.Type)
			assert.Equal(t, testCase.jsonEventName, event.ItemType)
			assert.Equal(t, reflect.TypeOf(eventObj).String(),
							reflect.TypeOf(event.Item).String()[1:])
		})
	}
}
