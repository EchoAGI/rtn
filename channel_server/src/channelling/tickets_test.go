package channelling

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func getRandom(n int) ([]byte, error) {
	result := make([]byte, n)
	if _, err := rand.Read(result); err != nil {
		return nil, err
	}
	return result, nil
}

func Test_ReverseBase64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		data, err := getRandom(64)
		if err != nil {
			t.Errorf("Could not get random data: %v", err)
			continue
		}

		s := base64.URLEncoding.EncodeToString(data)
		reversed, err := reverseBase64String(s)
		if err != nil {
			t.Errorf("Could not reverse %s: %v", s, err)
			continue
		}

		if s == reversed {
			t.Errorf("Reversing should be different for %s", s)
		}

		original, err := reverseBase64String(reversed)
		if err != nil {
			t.Errorf("Could not reverse back %s: %v", reversed, err)
			continue
		}

		if s != original {
			t.Errorf("Reversing back should have restored %s from %s but got %s", s, reversed, original)
		}
	}
}

func Test_Sessions(t *testing.T) {
	sessionSecret, err := getRandom(64)
	if err != nil {
		t.Fatalf("Could not create session secret: %v", err)
		return
	}

	encryptionSecret, err := getRandom(32)
	if err != nil {
		t.Fatalf("Could not create encryption secret: %v", err)
		return
	}

	tickets := NewTickets(sessionSecret, encryptionSecret, "test")
	silentOutput = true
	for i := 0; i < 1000; i++ {
		st := tickets.DecodeSessionToken("")
		if st == nil {
			t.Error("Could not create session")
			continue
		}

		if !tickets.ValidateSession(st.Id, st.Sid) {
			t.Errorf("Session is invalid: %v", st)
			continue
		}
	}
	silentOutput = false
}
