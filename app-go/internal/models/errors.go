package models

import "errors"

var (
	ErrInvalidUsername = errors.New("nombre de usuario inválido (3-20 caracteres: letras minúsculas, números o _)")
	ErrInvalidUser     = errors.New("remitente o destinatario inválido")
	ErrSameUser        = errors.New("no puedes enviarte un mensaje a ti mismo")
	ErrEmptyMessage    = errors.New("el mensaje no puede estar vacío")
	ErrMessageTooLong  = errors.New("el mensaje supera la longitud máxima permitida")
)
