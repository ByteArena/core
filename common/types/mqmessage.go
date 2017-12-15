package types

import (
	"os"
	"runtime"
	"time"
)

type MQPayload map[string]interface{}

type MQMessage struct {
	Date    time.Time  `json:"date"`
	Host    string     `json:"host"`
	App     string     `json:"app"`
	File    string     `json:"file"`
	Line    int        `json:"line"`
	Message string     `json:"message"`
	Payload *MQPayload `json:"payload"`
	IsError bool       `json:"error"`
}

func NewMQError(app string, msg string) *MQMessage {
	mqerr := NewMQMessage(app, msg)
	mqerr.IsError = true

	return mqerr
}

func NewMQMessage(app string, msg string) *MQMessage {

	mqerr := &MQMessage{
		Date:    time.Now(),
		App:     app,
		Message: msg,
	}

	if _, file, line, ok := runtime.Caller(1); ok {
		mqerr.File = file
		mqerr.Line = line
	}

	if hostname, err := os.Hostname(); err == nil {
		mqerr.Host = hostname
	}

	return mqerr
}

func (mqerr *MQMessage) SetPayload(payload MQPayload) *MQMessage {
	mqerr.Payload = &payload
	return mqerr
}
