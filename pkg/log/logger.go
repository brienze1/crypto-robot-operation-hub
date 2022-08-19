package log

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

var once sync.Once

type (
	logger struct {
		Timestamp     string      `json:"timestamp"`
		Level         string      `json:"level"`
		ErrorMsg      string      `json:"error,omitempty"`
		Message       string      `json:"message"`
		Metadata      interface{} `json:"metadata,omitempty"`
		CorrelationId string      `json:"correlationId"`
		TransactionId string      `json:"transactionId"`
	}
)

var loggerInstance *logger

func Logger() *logger {
	if loggerInstance == nil {
		once.Do(
			func() {
				loggerInstance = &logger{}
				loggerInstance.TransactionId = uuid.New().String()
			})
	}

	return loggerInstance
}

func (l *logger) SetCorrelationID(correlationId string) {
	l.CorrelationId = correlationId
}

func (l *logger) Info(message string, metadata ...interface{}) {
	logMessage := l.generateLogMessage("INFO ", message, nil, metadata)
	fmt.Println(logMessage)
}

func (l *logger) Error(err error, message string, metadata ...interface{}) {
	logMessage := l.generateLogMessage("ERROR", message, err, metadata)
	fmt.Println(logMessage)
}

func (l *logger) generateLogMessage(level string, message string, err error, metadata ...[]interface{}) string {
	log := l.clone()
	log.Level = level
	log.Message = message
	log.Timestamp = time.Now().Format("2022-01-01 13:01:01")

	if len(metadata[0]) > 0 {
		log.Metadata = metadata
	}

	if err != nil {
		log.ErrorMsg = err.Error()
	}

	logMessage, _ := json.Marshal(log)
	return string(logMessage)
}

func (l *logger) clone() logger {
	return *l
}
