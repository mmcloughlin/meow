package meow_test

import (
	"testing"

	"github.com/mmcloughlin/meow"
)

func TestVectors(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		got := meow.Checksum(v.Seed, v.Input)
		AssertBytesEqual(t, v.Hash, got[:])
	}
}

func TestVersions(t *testing.T) {
	testdata := LoadTestData(t)
	if meow.Version != testdata.Version {
		t.Errorf("version mismatch got=%d reference=%d", meow.Version, testdata.Version)
	}
	if meow.VersionName != testdata.VersionName {
		t.Errorf("version name mismatch got=%s reference=%s", meow.VersionName, testdata.VersionName)
	}
}
