package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtectEventUnmarshalJSON(t *testing.T) {
	// TODO: More test cases, test all types.
	rawJSON := []byte(`{
  "type": "add",
  "item": {
    "id": "66d025b301ebc903e80003ea",
    "modelKey": "event",
    "start": 1445408038748,
    "end": 1445408048748,
    "device": "66d025b301ebc903e80003ea",
	"type": "ring"
  }
}`)
	var event ProtectEvent
	err := event.UnmarshalJSON(rawJSON)
	assert.NoError(t, err)

	assert.Equal(t, "add", event.Type)
	assert.Equal(t, "ring", event.ItemType)
	assert.IsType(t, &RingEvent{}, event.Item)
}
