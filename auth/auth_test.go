package auth

import "testing"

func TestHash(t *testing.T) {
	result := Hash("bees", "honey")

	expected := "03b5c64d26016b15f14dac64a3deafccaa9ef72655ae09753c554e6b940ac5a2"

	if result != expected {
		t.Log("Expected %s", expected)
		t.Log("Received %s", result)
		t.FailNow()
	}
	if len(result) != 64 {
		t.Errorf("Incorrect Length!")
		t.FailNow()
	}
	t.Log("Passed!")
}
