package reconnect

import (
	"math/rand"
	"time"
)

// Retry a Reconnect on r with exponential backoff algorithm, at most c times and slot time
// set to interval. (see http://en.wikipedia.org/wiki/Exponential_backoff for details)
func ExponentialBackoff(r Reconnectable, interval time.Duration, max uint) (err error) {
	for retries := uint(0); retries < max; retries++ {
		var retry bool
		if err, retry = reconnectOnce(r); !retry {
			return err
		}
		switch retries {
		case 0:
			time.Sleep(time.Duration(rand.Int63n(int64(interval))))
		case 1:
			time.Sleep(time.Duration(rand.Int63n(7)) * interval)
		default:
			time.Sleep(time.Duration(rand.Int63n((1<<retries)-1)) * interval)
		}
	}
	return err
}
