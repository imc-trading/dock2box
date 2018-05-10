package model

import (
	"time"

	"github.com/pborman/uuid"
)

type Message struct {
	UUID     string    `json:"uuid"`
	Created  time.Time `json:"created"`
	HostUUID string    `json:"hostUUID"`
	TaskUUID string    `json:"taskUUID"`
	Line     int       `json:"line"`
	Message  string    `json:"message"`
}

type Messages []*Message

func NewMessage(hostUUID string, taskUUID string, line int, message string) *Message {
	return &Message{
		UUID:     uuid.New(),
		Created:  time.Now(),
		HostUUID: hostUUID,
		TaskUUID: taskUUID,
		Line:     line,
		Message:  message,
	}
}
