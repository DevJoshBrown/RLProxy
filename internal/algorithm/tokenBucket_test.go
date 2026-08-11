package algorithm

import (
	"testing"
	"time"
)

func TestTokenBucket_BlockOnce(t *testing.T) {

	tb := NewTokenBucket(1.0, 0.0, 0.0)
	result := tb.Allow()

	if result != false {
		t.Errorf("expected: false, got: %t", result)
	}
}

func TestTokenBucket_AllowOnce(t *testing.T) {

	tb := NewTokenBucket(1.0, 0.0, 1.0)
	result := tb.Allow()

	if result != true {
		t.Errorf("expected: true, got: %t", result)
	}
}

func TestTokenBucket_AllowOnceThenBlock(t *testing.T) {
	tb := NewTokenBucket(1.0, 0.0, 1.0)
	var result bool
	for i := 0; i < 2; i++ {
		result = tb.Allow()

		if i == 0 {
			if result != true {
				t.Errorf("expected true, got: %t", result)
			}
		}
		if i == 1 {
			if result != false {
				t.Errorf("expected false, got: %t", result)
			}
		}
	}
}

func TestTokenBucket_EqualRefill(t *testing.T) {

	tb := NewTokenBucket(10.0, 1.0, 1.0)

	var result bool
	allowedReqs := 0

	result = tb.Allow()
	allowedReqs += 1

	for i := 0; i < 4; i++ {
		tb.TimeOfLastRequest = tb.TimeOfLastRequest.Add(time.Second * -1)
		result = tb.Allow()
		allowedReqs += 1
		if result != true {
			t.Errorf("expected true, got: %t", result)
		}
	}

	if allowedReqs != 5 {
		t.Errorf("expected 5 requests, got: %d", allowedReqs)
	}
}

func TestTokenBucket_fiveThenFail(t *testing.T) {

	tb := NewTokenBucket(5.0, 0.5, 3.0)

	var result bool
	allowedReqs := 0

	result = tb.Allow()
	allowedReqs += 1

	for i := 0; i < 5; i++ {
		tb.TimeOfLastRequest = tb.TimeOfLastRequest.Add(time.Second * -1)
		result = tb.Allow()
		if i < 4 {
			if result != true {
				t.Errorf("expected true, got: %t", result)
			}
			allowedReqs += 1
		} else {
			if result != false {
				t.Errorf("expected false, got: %t", result)
			}
		}
	}
	if allowedReqs != 5 {
		t.Errorf("expected 5 requests, got: %d", allowedReqs)
	}
}
