package helper

import "testing"

// TestWantError tests wantError.
func TestWantError(t *testing.T, err error, want bool) {
	// mark as helper function
	t.Helper()

	if err != nil && !want {
		t.Errorf("got an error %v, want nothing happened", err)
	} else if err == nil && want {
		t.Error("got nothing happened, want an error")
	}
}
