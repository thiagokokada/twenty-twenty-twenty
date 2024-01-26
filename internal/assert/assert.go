package assert

import (
	"cmp"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func GreaterOrEqual[T cmp.Ordered](t *testing.T, actual, expected T) {
	t.Helper()
	if actual < expected {
		t.Errorf("got: %v; want: >=%v", actual, expected)
	}
}
