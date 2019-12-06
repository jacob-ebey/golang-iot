package azuregateway

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/amenzhinsky/iothub/iotdevice"
	iotmqtt "github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
	"github.com/amenzhinsky/iothub/iotservice"
	"github.com/amenzhinsky/iothub/logger"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/jacob-ebey/golang-iot/core"
)

// TODO: Test fail cases for serialize, send, subscribe, deserialize

func TestAzureReaderReceivesMessage(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Log(err)
	}

	client, err := iotdevice.NewFromConnectionString(
		iotmqtt.New(), os.Getenv("IOTHUB_DEVICE_CONNECTION_STRING"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err = client.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}

	logger := core.TestLogger{}
	reader, err := NewAzureReader(context.Background(), &logger, client)
	if err != nil {
		t.Fatal(err)
	}

	service, err := iotservice.NewFromConnectionString(
		os.Getenv("IOTHUB_SERVICE_CONNECTION_STRING"),
	)
	if err != nil {
		t.Fatal(err)
	}

	serializer := &core.JsonSerializer{}
	toSend := core.Message{
		Type: "LOL",
		Payload: map[string]interface{}{
			"LMFAO": "ROFLAO",
		},
	}
	payload, err := serializer.SerialiseMessage(toSend)
	if err != nil {
		t.Fatal(err)
	}

	if err := service.SendEvent(context.Background(), client.DeviceID(), payload); err != nil {
		t.Fatal(err)
	}

	received := <-reader
	if diff := deep.Equal(toSend, received); diff != nil {
		t.Fatal(diff)
	}
}

func TestAzureWriterSendsMessage(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Log(err)
	}

	client, err := iotdevice.NewFromConnectionString(
		iotmqtt.New(), os.Getenv("IOTHUB_DEVICE_CONNECTION_STRING"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err = client.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}

	log := core.TestLogger{}

	correlationID := uuid.New().String()
	writer := NewAzureWriter(context.Background(), &log, client, iotdevice.WithSendCorrelationID(correlationID))

	service, err := iotservice.NewFromConnectionString(
		os.Getenv("IOTHUB_SERVICE_CONNECTION_STRING"),
		iotservice.WithLogger(logger.New(1, func(lvl logger.Level, s string) {
			fmt.Println("IOT_HUB: " + s)
		})),
	)
	if err != nil {
		t.Fatal(err)
	}

	toSend := core.Message{
		Type: "LOL",
		Payload: map[string]interface{}{
			"LMFAO": "ROFLAO",
		},
	}

	receive := make(chan core.Message)
	serializer := &core.JsonSerializer{}
	go func() {
		service.SubscribeEvents(context.Background(), func(msg *iotservice.Event) error {
			if msg.CorrelationID != correlationID {
				return nil
			}
			if msg == nil {
				return nil
			}

			message, err := serializer.DeserialiseMessage(msg.Payload)
			if err != nil {
				t.Fatal(err)
			}

			receive <- message

			return nil
		})
	}()

	<-time.After(10 * time.Second)
	writer <- toSend

	message := <-receive
	if diff := deep.Equal(toSend, message); diff != nil {
		t.Fatal(diff)
	}
}
