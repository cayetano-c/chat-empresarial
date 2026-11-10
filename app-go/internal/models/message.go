package models

import (
	"strings"
	"time"
	"unicode/utf8"
)

const MaxMessageLength = 500

type Message struct {
	ID        int       `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMessage(from, to, text string) (*Message, error) {
	from = strings.ToLower(strings.TrimSpace(from))
	to = strings.ToLower(strings.TrimSpace(to))
	text = strings.TrimSpace(text)

	if from == "" || to == "" {
		return nil, ErrInvalidUser
	}
	if from == to {
		return nil, ErrSameUser
	}
	if text == "" {
		return nil, ErrEmptyMessage
	}
	if utf8.RuneCountInString(text) > MaxMessageLength {
		return nil, ErrMessageTooLong
	}
	return &Message{From: from, To: to, Text: text, CreatedAt: time.Now()}, nil
}

func (m *Message) IsBetween(a, b string) bool {
	return (m.From == a && m.To == b) || (m.From == b && m.To == a)
}
