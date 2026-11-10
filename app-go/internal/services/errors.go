package services

import "errors"

var (
	ErrUserNotFound = errors.New("el usuario no existe")
	ErrInvalidPage = errors.New("número de página inválido")
)
