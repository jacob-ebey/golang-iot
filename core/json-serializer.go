package core

import (
	"encoding/json"
)

type JsonSerializer struct{}

func (serializer *JsonSerializer) SerialiseMessage(message Message) ([]byte, error) {
	return json.Marshal(&message)
}

func (serializer *JsonSerializer) DeserialiseMessage(data []byte) (Message, error) {
	result := Message{}
	err := json.Unmarshal(data, &result)
	return result, err
}
