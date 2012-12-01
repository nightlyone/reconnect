package reconnect

import (
	"errors"
)

// Just a marker for final connect errors
type FinalError struct {
	RealError error
}

func (f FinalError) Error() string {
	return f.RealError.Error()
}

var ErrAlreadyConnected = errors.New("connection already established")

// A Reconnectable will re-establish a connection using internal state
//
// Reconnect tries to re-establishing a broken connection exactly once.
// If called, while the connection is already established, it should return ErrAlreadyConnected.
// If calling Reconnect once more makes no sense, if should return a FinalError type of error.
type Reconnectable interface {
	Reconnect() error
}
