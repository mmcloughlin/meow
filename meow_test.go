package meow

import "testing"

func TestVectors(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		got := Checksum(v.Seed, v.Input)
		AssertBytesEqual(t, v.Hash, got[:])
	}
}

func TestVersions(t *testing.T) {
	testdata := LoadTestData(t)
	if Version != testdata.Version {
		t.Errorf("version mismatch got=%d reference=%d", Version, testdata.Version)
	}
	if VersionName != testdata.VersionName {
		t.Errorf("version name mismatch got=%s reference=%s", VersionName, testdata.VersionName)
	}
}

func TestDisplayImplementation(t *testing.T) {
	t.Logf("implementation=%s", implementation)
}
