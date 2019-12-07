package azuregateway

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/amenzhinsky/iothub/iotdevice"
	iotmqtt "github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
	"github.com/amenzhinsky/iothub/iotservice"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

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

	ctx, done := context.WithCancel(context.TODO())
	reader, errors := NewAzureReader(ctx, client)

	service, err := iotservice.NewFromConnectionString(os.Getenv("IOTHUB_SERVICE_CONNECTION_STRING"))
	if err != nil {
		t.Fatal(err)
	}

	toSend := []byte{1, 2, 3, 4}

	if err := service.SendEvent(context.Background(), client.DeviceID(), toSend); err != nil {
		t.Fatal(err)
	}

	select {
	case payload := <-reader:
		if len(payload) != len(toSend) {
			t.Fatal("payload length is not correct")
		}

		for index, b := range payload {
			if b != toSend[index] {
				t.Fatal("payload value is not correct")
			}
		}
	case err := <-errors:
		t.Fatal(err)
	}

	done()
}

func TestAzureWriterSendsMessage(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Log(err)
	}

	service, err := iotservice.NewFromConnectionString(os.Getenv("IOTHUB_SERVICE_CONNECTION_STRING"))
	if err != nil {
		t.Fatal(err)
	}

	ctx, done := context.WithCancel(context.TODO())
	receive := make(chan []byte)
	correlationID := uuid.New().String()
	go func() {
		if err := service.SubscribeEvents(context.Background(), func(msg *iotservice.Event) error {
			if msg.CorrelationID != correlationID {
				return nil
			}
			if msg == nil {
				return nil
			}

			receive <- msg.Payload
			done()

			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}()
	<-time.After(5 * time.Second)

	client, err := iotdevice.NewFromConnectionString(iotmqtt.New(), os.Getenv("IOTHUB_DEVICE_CONNECTION_STRING"))
	if err != nil {
		t.Fatal(err)
	}

	if err = client.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	toSend := []byte{1, 2, 3, 4}

	writer, errors := NewAzureWriter(ctx, client, iotdevice.WithSendCorrelationID(correlationID))
	writer <- toSend

	select {
	case payload := <-receive:
		if len(payload) != len(toSend) {
			t.Fatal("payload length is not correct")
		}

		for index, b := range payload {
			if b != toSend[index] {
				t.Fatal("payload value is not correct")
			}
		}

	case err := <-errors:
		t.Fatal(err)

	case <-ctx.Done():
		break
	}
}

func TestAzureWriterForwardsErrors(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Log(err)
	}

	ctx := context.TODO()

	client, err := iotdevice.NewFromConnectionString(iotmqtt.New(), os.Getenv("IOTHUB_DEVICE_CONNECTION_STRING"))
	if err != nil {
		t.Fatal(err)
	}

	if err = client.Connect(ctx); err != nil {
		t.Fatal(err)
	}
	client.Close()

	toSend := []byte{1, 2, 3, 4}

	writer, errors := NewAzureWriter(ctx, client)
	writer <- toSend

	select {
	case err := <-errors:
		if err == nil {
			t.Fatal("error was nil")
		}
	}
}
