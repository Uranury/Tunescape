package apperrors

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrBadRequest          = errors.New("bad request")
	ErrEmailTaken          = errors.New("email already in use")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrSpotifyIDTaken      = errors.New("spotify account already linked to another user")
	ErrSpotifyNotConnected = errors.New("spotify account not connected")
	ErrNoSnapshot          = errors.New("no snapshot found for user")
)
