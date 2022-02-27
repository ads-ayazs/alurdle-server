package game

import "errors"

// Error codes returned by functions implemented in this package
var (
	ErrSerialization = errors.New("game serialization error")
	ErrGameOver      = errors.New("game is finished")
	ErrOutOfTurns    = errors.New("out of turns")
	ErrNilResult     = errors.New("nil result provided")
	ErrWordLength    = errors.New("invalid word length")
	ErrInvalidWord   = errors.New("word is not in dictionary")
	// ErrInvalidId     = errors.New("invalid id")
)
