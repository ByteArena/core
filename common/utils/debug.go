package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Context map[string]interface{}

type Message struct {
	Time    string  `json:"time"`
	Level   string  `json:"lvl"`
	Service string  `json:"service"`
	Message string  `json:"message"`
	Context Context `json:"context"`
}

var (
	LogFn = debug
)

func Debug(service string, message string) {
	LogFn(service, message)
}

func debug(service string, message string) {
	context := make(Context, 0)

	if hostname, err := os.Hostname(); err == nil {
		context["hostname"] = hostname
	}

	messageStruct := Message{
		Time:    time.Now().Format(time.RFC3339),
		Service: service,
		Message: message,
		Context: context,
		Level:   "debug",
	}

	data, _ := json.Marshal(messageStruct)

	fmt.Println(string(data))
}

func RecoverableError(service, message string) {
	context := make(Context, 0)

	if hostname, err := os.Hostname(); err == nil {
		context["hostname"] = hostname
	}

	messageStruct := Message{
		Time:    time.Now().Format(time.RFC3339),
		Service: service,
		Message: message,
		Context: context,
		Level:   "error",
	}

	data, _ := json.Marshal(messageStruct)

	fmt.Println(string(data))
}
