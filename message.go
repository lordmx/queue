package main

import (
	"time"
)

type Message struct {
	ID        string
	Body      []byte
	CreatedAt time.Time
	Queue     *Queue
	Client    *Client
}

func NewMessage(queue *Queue, client *Client, body []byte) *Message {
	return &Message{
		ID:        RandString(32),
		Body:      body,
		CreatedAt: time.Now(),
		Queue:     queue,
		Client:    client,
	}
}
