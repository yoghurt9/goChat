package main

import (
	"time"
)

type Message struct {
	// Sender's ID
	FromId string
	// Recipient's ID
	ToId string
	// Content
	Content string

	CreateTime time.Time
}

