package reconnect

import (
	"time"
)

// Try one Reconnect on r and return whether we should give up with err or retry.
func reconnectOnce(r Reconnectable) (err error, retry bool) {
	err = r.Reconnect()
	switch err {
	case nil:
		return nil, false
	case ErrAlreadyConnected:
		return nil, false
	default:
		if final, ok := err.(FinalError); ok {
			return final.RealError, false
		}
		return err, true
	}
	panic("Not reached")
}

// Retry a Reconnect on r after period of time
func After(r Reconnectable, after time.Duration) error {
	time.Sleep(after)
	err, _ := reconnectOnce(r)
	return err
}

// Retry a Reconnect on r as long as the interactive function again returns true or we receive a final state.
// Useful, if you like to give the user a choice after he changed the connection parameters
func Interactive(r Reconnectable, again func(r Reconnectable) bool) (err error) {
	for {
		var retry bool
		err, retry = reconnectOnce(r)
		if retry && again(r) {
			continue
		} else {
			break
		}
	}
	return err
}

// Retry a Reconnect on r for c times, sleeping interval between each try
func NTimes(r Reconnectable, interval time.Duration, max uint) (err error) {
	for retries := uint(0); retries < max; retries++ {
		var retry bool
		if err, retry = reconnectOnce(r); !retry {
			return err
		}
		time.Sleep(interval)
	}
	return err
}


