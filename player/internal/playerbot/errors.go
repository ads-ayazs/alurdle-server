package playerbot

import "errors"

var (
	ErrInvalidBotName = errors.New("invalid bot name")
	ErrFailedAttempt  = errors.New("failed attempt")
)
