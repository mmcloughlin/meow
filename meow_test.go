package meow

import (
	"math/rand"
	"testing"
	"testing/quick"
)

// ChecksumSlice adapts Checksum to return a slice instead of an array.
func ChecksumSlice(seed uint64, data []byte) []byte {
	cksum := Checksum(seed, data)
	return cksum[:]
}

// ChecksumHash implements Checksum with the hash.Hash interface. Intended to facilitate comparison between the two.
func ChecksumHash(seed uint64, data []byte) []byte {
	h := New(seed)
	h.Write(data)
	return h.Sum(nil)
}

// ChecksumHashWithReset is intended to confirm hash.Hash Reset() behavior.
// Hashes some random data, resets and then computes the desired hash.
func ChecksumHashWithReset(seed uint64, data []byte) []byte {
	n := rand.Intn(9 << 10)
	r := make([]byte, n)
	rand.Read(r)

	h := New(seed)
	h.Write(r)
	h.Reset()
	h.Write(data)
	return h.Sum(nil)
}

// ChecksumRandomBatchedHash implements Checksum by writing random amounts to a hash.Hash.
func ChecksumRandomBatchedHash(seed uint64, data []byte) []byte {
	h := New(seed)
	for len(data) > 0 {
		n := rand.Intn(len(data) + 1)
		h.Write(data[:n])
		data = data[n:]
	}
	return h.Sum(nil)
}

func TestVectorsChecksum(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		AssertBytesEqual(t, v.Hash, ChecksumSlice(v.Seed, v.Input))
	}
}

func TestVectorsHash(t *testing.T) {
	testdata := LoadTestData(t)
	for _, v := range testdata.TestVectors {
		AssertBytesEqual(t, v.Hash, ChecksumHash(v.Seed, v.Input))
	}
}

func TestQuickChecksumMatchesHash(t *testing.T) {
	if err := quick.CheckEqual(ChecksumSlice, ChecksumHash, nil); err != nil {
		t.Fatal(err)
	}
}

func TestQuickHashBatching(t *testing.T) {
	if err := quick.CheckEqual(ChecksumHash, ChecksumRandomBatchedHash, nil); err != nil {
		t.Fatal(err)
	}
}

func TestQuickHashReset(t *testing.T) {
	if err := quick.CheckEqual(ChecksumHash, ChecksumHashWithReset, nil); err != nil {
		t.Fatal(err)
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
