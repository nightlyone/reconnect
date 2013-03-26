package reconnect

import (
	"math/rand"
	"time"
)

// Given a retry count and a base retry interval, return the next retry interval
// according to the exponential backoff algorithm.
// Maximum returned retry interval is interval * 2^15.
func nextExponentialBackoff(retries uint, interval time.Duration) time.Duration {
	switch {
	case retries == 0:
		return time.Duration(rand.Int63n(int64(interval)))
	case retries == 1:
		return time.Duration(rand.Int63n(7)) * interval
	case retries >= 2 && retries < 15:
		return time.Duration(rand.Int63n((1<<retries)-1)) * interval
	default:
	}
	return time.Duration(rand.Int63n(1<<15)) * interval
}

// Retry a Reconnect on r with exponential backoff algorithm, at most c times and slot time
// set to interval. (see http://en.wikipedia.org/wiki/Exponential_backoff for details)
func ExponentialBackoff(r Reconnectable, interval time.Duration, max uint) (err error) {
	for retries := uint(0); retries < max; retries++ {
		var retry bool
		if err, retry = reconnectOnce(r); !retry {
			return err
		}
		time.Sleep(nextExponentialBackoff(retries, interval))
	}
	return err
}
