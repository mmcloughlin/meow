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

func TestVectorsChecksum64(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		sum := Checksum64(v.Seed, v.Input)
		if sum != v.Hash64 {
			t.Fatalf("got=%016x expect%016x", sum, v.Hash32)
		}
	}
}

func TestVectorsChecksum32(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		sum := Checksum32(v.Seed, v.Input)
		if sum != v.Hash32 {
			t.Fatalf("got=%08x expect%08x", sum, v.Hash32)
		}
	}
}

func TestVectorsHash(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		AssertBytesEqual(t, v.Hash, checksumHash(v.Seed, v.Input))
	}
}

func TestVectorsHash64(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		h := New64(v.Seed)
		h.Write(v.Input)
		sum := h.Sum64()
		if sum != v.Hash64 {
			t.Fatalf("got=%016x expect%016x", sum, v.Hash32)
		}
	}
}

func TestVectorsHash32(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		h := New32(v.Seed)
		h.Write(v.Input)
		sum := h.Sum32()
		if sum != v.Hash32 {
			t.Fatalf("got=%08x expect%08x", sum, v.Hash32)
		}
	}
}

func TestHashSizes(t *testing.T) {
	AssertHashSize(t, "New", New(0), Size)
	AssertHashSize(t, "New64", New64(0), 8)
	AssertHashSize(t, "New32", New32(0), 4)
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
