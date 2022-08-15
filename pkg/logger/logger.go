package logger

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

var once sync.Once

type logger struct {
	correlationId string
	transactionId string
	message       string
	metadata      interface{}
	timestamp     string
}

var loggerInstance *logger

func getInstance() *logger {
	if loggerInstance == nil {
		once.Do(
			func() {
				loggerInstance = &logger{}
				loggerInstance.transactionId = uuid.New().String()
			})
	} else {
		fmt.Println("Single instance already created.")
	}

	return loggerInstance
}

func SetCorrelationID(correlationId string) {
	loggerInstance.correlationId = correlationId
}

func Info(message string, metadata ...interface{}) {
	logMessage := generateLogMessage(message, metadata)
	log.Println(logMessage)
}

func Error(message string, metadata ...interface{}) {
	logMessage := generateLogMessage(message, metadata)
	log.Fatalf(logMessage)
}

func generateLogMessage(message string, metadata ...interface{}) string {
	log := getInstance()
	log.message = message
	log.metadata = metadata
	log.timestamp = time.Now().Format("2022-01-01 13:01:01")
	logMessage, _ := json.Marshal(loggerInstance)
	log.message = ""
	log.timestamp = ""
	log.metadata = nil
	return string(logMessage)
}
