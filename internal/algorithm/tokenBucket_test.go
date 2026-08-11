package algorithm

import (
	"testing"
)

var requests int

func TestTokenBucket_BlockWhenEmpty(t *testing.T) {

	tb := NewTokenBucket(1.0, 0.0, 0.0)
	result := tb.Allow()

	if result != false {
		t.Errorf("expected: false, got: %t", result)
	}
}

func TestTokenBucket_AllowWhenFilled(t *testing.T) {

	tb := NewTokenBucket(1.0, 0.0, 1.0)
	result := tb.Allow()

	if result != true {
		t.Errorf("expected: true, got: %t", result)
	}
}
