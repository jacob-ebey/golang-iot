package core

var (
	MessageTypeError           = "ERROR"
	MessageTypeWarning         = "WARNING"
	MessageTypeData            = "DATA"
	MessageTypeAnalytics       = "ANALYTICS"
	MessageTypeSetupPeripheral = "SETUP_PERIPHERAL"
)

type MessageSerializer interface {
	SerialiseMessage(message Message) ([]byte, error)
}

type MessageDeserializer interface {
	DeserialiseMessage(data []byte) (Message, error)
}

type Message struct {
	Type    string
	Payload interface{}
}

type ErrorPayload struct {
	Code    string
	Message string
}

type SetupPeripheralPayload struct {
	NewDeviceID string
}
