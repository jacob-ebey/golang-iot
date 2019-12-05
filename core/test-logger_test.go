package core

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
)

func TestTestLoggerDoesNotBlock(t *testing.T) {
	logger := TestLogger{}
	err := fmt.Errorf("ROFL")
	logger.LogError(err)

	if diff := deep.Equal(err, logger.LoggedErrors[0]); diff != nil {
		t.Fatal(diff)
	}
}

func TestTestLoggerBlocksUntilError(t *testing.T) {
	errors := make(chan error)
	logger := TestLogger{
		LoggedErrorsChannel: errors,
	}
	err := fmt.Errorf("ROFL")
	go func() {
		logger.LogError(err)
	}()

	if diff := deep.Equal(err, <-errors); diff != nil {
		t.Fatal(diff)
	}
}
