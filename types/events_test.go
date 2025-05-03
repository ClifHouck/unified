package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ClifHouck/unified/types"
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
	var event types.ProtectEvent
	err := event.UnmarshalJSON(rawJSON)
	require.NoError(t, err)

	assert.Equal(t, "add", event.Type)
	assert.Equal(t, "ring", event.ItemType)
	assert.IsType(t, &types.RingEvent{}, event.Item)
}
