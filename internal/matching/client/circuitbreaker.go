package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/s3f4/locationmatcher/pkg/log"
)

type Circuit func(ctx context.Context, url string, reader io.Reader) (*http.Response, error)

// The Breaker function accepts any function that conforms to the Circuit type
// definition, and an unsigned integer representing the number of consecutive failures
// allowed before the circuit automatically opens.
// In return it provides another function,
// which also conforms to the Circuit type definition:
func Breaker(circuit Circuit, failureTreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context, url string, reader io.Reader) (*http.Response, error) {
		m.RLock() // Establish a "read lock"
		d := consecutiveFailures - int(failureTreshold)

		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			log.Info(time.Now())
			log.Info(shouldRetryAt)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				log.Info("service unreachable")
				return nil, errors.New("service unreachable")
			}
		}

		m.RUnlock() // Release read lock

		response, err := circuit(ctx, url, reader) // Issue request proper

		m.Lock() // Lock around shared resources
		defer m.Unlock()

		lastAttempt = time.Now() // Record time of attempt

		if err != nil { // Circuit returned an error,
			consecutiveFailures++ // so we count the failure
			return response, err  // and return
		}

		consecutiveFailures = 0 // Reset failures counter
		return response, nil
	}
}
