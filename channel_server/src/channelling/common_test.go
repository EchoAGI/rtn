package channelling

import (
	"testing"
)

func assertDataError(t *testing.T, err error, code string) {
	if err == nil {
		t.Error("Expected an error, but none was returned")
		return
	}

	dataError, ok := err.(*DataError)
	if !ok {
		t.Errorf("Expected error %#v to be a *DataError", err)
		return
	}

	if code != dataError.Code {
		t.Errorf("Expected error code to be %v, but was %v", code, dataError.Code)
	}
}
