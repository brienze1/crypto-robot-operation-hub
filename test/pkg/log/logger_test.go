package log

import (
	"bytes"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/adapters"
	logg "github.com/brienze1/crypto-robot-operation-hub/pkg/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var (
	logger        adapters.LoggerAdapter
	testMessage   string
	testMetadata1 interface{}
	testMetadata2 interface{}
	buf           bytes.Buffer
)

func setup() {
	logger = logg.Logger()

	log.SetOutput(&buf)

	testMessage = uuid.NewString()
	testMetadata1 = map[string]interface{}{"Key1": 123, "Key2": "value2", "Test": uuid.NewString()}
	testMetadata2 = map[string]interface{}{"Key3": 1234, "Key4": "value3", "Test2": uuid.NewString()}
}

func teardown() {
	log.SetOutput(os.Stderr)
}

func TestLoggerInfoMessageOnlySuccess(t *testing.T) {
	setup()
	defer func() {
		teardown()
	}()

	logger.Info(testMessage)

	assert.Contains(t, buf.String(), testMessage)
}

func TestLoggerInfoOneMetadataSuccess(t *testing.T) {
	setup()
	defer func() {
		teardown()
	}()

	logger.Info(testMessage, testMetadata1)

	assert.Contains(t, buf.String(), testMessage)
}
