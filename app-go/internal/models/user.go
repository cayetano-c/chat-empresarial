package models

import (
	"regexp"
	"strings"
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9_]{3,20}$`)

// User representa a un empleado que participa en el chat.
// Sus campos no son exportados: solo se modifican mediante métodos.
type User struct {
	username    string
	displayName string
}

// NewUser es el constructor: garantiza que un User siempre tenga un estado válido.
func NewUser(username, displayName string) (*User, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	displayName = strings.TrimSpace(displayName)

	if !usernamePattern.MatchString(username) {
		return nil, ErrInvalidUsername
	}
	if displayName == "" {
		displayName = username
	}
	return &User{username: username, displayName: displayName}, nil
}

// Username devuelve el identificador único del usuario.
func (u *User) Username() string { return u.username }

// DisplayName devuelve el nombre para mostrar.
func (u *User) DisplayName() string { return u.displayName }

// SetDisplayName usa receptor por puntero porque modifica el estado.
func (u *User) SetDisplayName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrInvalidUsername
	}
	u.displayName = name
	return nil
}
