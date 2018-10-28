// +build !noasm

package meow

import (
	"encoding/json"
	"testing"
)

// TestDisplayCPUFeatures is purely for debuging purposes.
func TestDisplayCPUFeatures(t *testing.T) {
	b, err := json.MarshalIndent(cpu, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
