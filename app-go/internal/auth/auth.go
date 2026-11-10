package auth

import "errors"

var ErrInvalidToken = errors.New("token inválido o vencido")

type Authenticator interface {
	Validate(token string) (string, error)
	UserExists(username string) bool
}
