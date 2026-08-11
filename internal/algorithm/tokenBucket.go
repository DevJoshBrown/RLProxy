package algorithm

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	mu                sync.Mutex
	Capacity          float64
	RefillRate        float64
	CurrentTokens     float64
	TimeOfLastRequest time.Time
}

func NewTokenBucket(Capacity float64, RefillRate float64, CurrentTokens float64) *TokenBucket {
	return &TokenBucket{
		Capacity:          Capacity,
		RefillRate:        RefillRate,
		CurrentTokens:     CurrentTokens,
		TimeOfLastRequest: time.Now(),
	}
}

func (t *TokenBucket) Allow() bool {

	t.mu.Lock()
	defer t.mu.Unlock()

	// calculate time the request is made.
	requestTime := time.Now()

	// calculate time since last request
	timeSinceLast := requestTime.Sub(t.TimeOfLastRequest)

	// calculate the seconds passed, and convert to 'x' tokens per second
	newTokens := t.RefillRate * timeSinceLast.Seconds()

	// check bucket not already at capacity
	if t.CurrentTokens < t.Capacity {

		// add the new tokens to the pool
		t.CurrentTokens += newTokens

		// limit new token overflow to capacity
		if t.CurrentTokens > t.Capacity {
			t.CurrentTokens = t.Capacity
		}
	}

	//update the time of last request the time the request was made.
	t.TimeOfLastRequest = requestTime

	// attempt a request
	if t.CurrentTokens >= 1 {
		t.CurrentTokens -= 1
		fmt.Print("request allowed")

		return true
	} else {
		fmt.Print("request denied")
		return false
	}

}
