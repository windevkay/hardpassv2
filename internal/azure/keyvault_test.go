package azure

import (
	"testing"
)

func TestGenSecureString(t *testing.T) {
	res, err := genSecureString()

	if err != nil {
		t.Error("genSecureString() failed")
	}

	if len(res) != 64 {
		t.Error("genSecureString() failed to generate 64 character string")
	}
}
