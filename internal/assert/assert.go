package assert

import (
	"strings"
	"testing"
)

// test helper to check equality between the value we expected and the value we got
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got %v; want %v", actual, expected)
	}
}

// test helper to check if a string contains a certain substring
func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("expected string %s to contain %s", actual, expectedSubstring)
	}
}

func NilError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("got %v; expected: nil", err)
	}
}
