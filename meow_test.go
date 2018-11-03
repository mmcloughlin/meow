package meow

import (
	"testing"
)

func TestVectorsChecksum(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		AssertBytesEqual(t, v.Hash, checksumSlice(v.Seed, v.Input))
	}
}

func TestVectorsHash(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		AssertBytesEqual(t, v.Hash, checksumHash(v.Seed, v.Input))
	}
}

func TestChecksumMatchesHash(t *testing.T) {
	CheckEqual(t, checksumSlice, checksumHash)
}

func TestHashBatching(t *testing.T) {
	CheckEqual(t, checksumHash, checksumRandomBatchedHash)
}

func TestHashReset(t *testing.T) {
	CheckEqual(t, checksumHash, checksumHashWithReset)
}

func TestHashSumPreservesState(t *testing.T) {
	CheckEqual(t, checksumHash, checksumHashWithIntermediateSum)
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
