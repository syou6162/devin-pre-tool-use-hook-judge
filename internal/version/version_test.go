package version

import "testing"

func TestString(t *testing.T) {
	if got := String(); got != Version {
		t.Errorf("String() = %q, want %q", got, Version)
	}
}
