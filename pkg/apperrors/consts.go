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
	ErrNoSnapshot            = errors.New("no snapshot found for user")
	ErrUpstreamUnavailable   = errors.New("upstream service unavailable")
	ErrAlreadyFriends        = errors.New("users are already friends")
	ErrRequestAlreadySent    = errors.New("friend request already sent")
	ErrRequestNotFound       = errors.New("friend request not found")
	ErrCannotAddSelf         = errors.New("cannot add yourself as a friend")
	ErrNotFriends            = errors.New("users are not friends")
)
