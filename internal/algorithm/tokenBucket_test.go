package algorithm

import (
	"testing"
	"time"
)

func TestTokenBucket_BlockOnce(t *testing.T) {

	tb := NewTokenBucket(1.0, 0.0, 0.0)
	result, retry := tb.Allow()
	if result != false {
		t.Errorf("expected: false, got: %t", result)
	}
	if retry != NoRetry {
		t.Errorf("expected: NoRetry, got: %v", retry)
	}
}

func TestTokenBucket_AllowOnce(t *testing.T) {

	tb := NewTokenBucket(1.0, 0.0, 1.0)
	result, retry := tb.Allow()

	if result != true {
		t.Errorf("expected: true, got: %t", result)
	}
	if retry != 0.0 {
		t.Errorf("expected: 0.0, got %v", retry)
	}
}

func TestTokenBucket_AllowOnceThenBlock(t *testing.T) {
	tb := NewTokenBucket(1.0, 0.0, 1.0)
	var result bool
	var retry time.Duration
	for i := 0; i < 2; i++ {
		result, retry = tb.Allow()

		if i == 0 {
			if result != true {
				t.Errorf("expected true, got: %t", result)
			}
			if retry != 0.0 {
				t.Errorf("expected: 0.0, got %v", retry)
			}
		}
		if i == 1 {
			if result != false {
				t.Errorf("expected false, got: %t", result)
			}
			if retry != NoRetry {
				t.Errorf("expected: NoRetry, got: %v", retry)
			}
		}
	}
}

func TestTokenBucket_EqualRefill(t *testing.T) {

	tb := NewTokenBucket(10.0, 1.0, 1.0)

	var result bool
	var retry time.Duration
	allowedReqs := 0

	// initial request to set a time of last request
	result, retry = tb.Allow()
	allowedReqs += 1

	for i := 0; i < 4; i++ {
		tb.TimeOfLastRequest = tb.TimeOfLastRequest.Add(time.Second * -1)
		result, retry = tb.Allow()
		allowedReqs += 1
		if result != true {
			t.Errorf("expected true, got: %t", result)
		}
		if retry != 0.0 {
			t.Errorf("expected 0.0, got %v", retry)
		}
	}

	if allowedReqs != 5 {
		t.Errorf("expected 5 requests, got: %d", allowedReqs)
	}
}

func TestTokenBucket_fiveThenFail(t *testing.T) {

	tb := NewTokenBucket(5.0, 0.5, 3.0)

	var result bool
	var retry time.Duration
	allowedReqs := 0

	result, retry = tb.Allow()
	allowedReqs += 1

	for i := 0; i < 5; i++ {
		tb.TimeOfLastRequest = tb.TimeOfLastRequest.Add(time.Second * -1)
		result, retry = tb.Allow()
		if i < 4 {
			if result != true {
				t.Errorf("expected true, got: %t", result)
			}
			if retry != 0.0 {
				t.Errorf("expected 0.0, got %v", retry)
			}
			allowedReqs += 1
		} else {
			if result != false {
				t.Errorf("expected false, got: %t", result)
			}
			if retry <= 0.0 {
				t.Errorf("expected a positive retry time, got: %v", retry)
			}
		}
	}
	if allowedReqs != 5 {
		t.Errorf("expected 5 requests, got: %d", allowedReqs)
	}
}

func TestTokenBucket_ConcurrentTenCapFive(t *testing.T) {

	tb := NewTokenBucket(5.0, 0.0, 5.0)
	ch := make(chan bool)
	iterations := 10
	trueCount := 0
	falseCount := 0

	for i := 0; i < iterations; i++ {
		go func() {
			result, _ := tb.Allow()
			ch <- result
		}()
	}

	for i := 0; i < iterations; i++ {
		result := <-ch
		if result == true {
			trueCount++
		} else {
			falseCount++
		}
	}

	if (trueCount != 5) || (falseCount != 5) {
		t.Errorf("expected 5 true and 5 false, got: %v true and %v false", trueCount, falseCount)
	}
}
