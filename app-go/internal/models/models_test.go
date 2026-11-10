package models

import (
	"errors"
	"strings"
	"testing"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		name, from, to, text string
		wantErr              error
	}{
		{"válido", "ana", "luis", "hola", nil},
		{"vacío", "ana", "luis", "   ", ErrEmptyMessage},
		{"sin remitente", "", "luis", "hola", ErrInvalidUser},
		{"mismo usuario", "ana", "ana", "hola", ErrSameUser},
		{"muy largo", "ana", "luis", strings.Repeat("a", MaxMessageLength+1), ErrMessageTooLong},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewMessage(tc.from, tc.to, tc.text)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, esperado %v", err, tc.wantErr)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name, username string
		wantErr        error
	}{
		{"válido", "ana_01", nil},
		{"con mayúsculas se normaliza", "ANA", nil},
		{"muy corto", "ab", ErrInvalidUsername},
		{"con espacios", "ana perez", ErrInvalidUsername},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewUser(tc.username, "")
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, esperado %v", err, tc.wantErr)
			}
		})
	}
}
