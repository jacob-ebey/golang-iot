package core

import (
	"testing"

	"github.com/go-test/deep"
)

func TestJsonSerializerSerializesAndDeserializes(t *testing.T) {
	message := Message{
		Type: "LOL",
		Payload: map[string]interface{}{
			"LMFAO": "ROFLAO",
		},
	}
	serializer := &JsonSerializer{}

	serialized, err := serializer.SerialiseMessage(message)
	if err != nil {
		t.Fatal(err)
	}

	deserialized, err := serializer.DeserialiseMessage(serialized)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(message, deserialized); diff != nil {
		t.Fatal(diff)
	}
}
