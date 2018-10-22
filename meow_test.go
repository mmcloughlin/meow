package meow_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/mmcloughlin/meow"
)

type TestVector struct {
	Seed  uint64
	Input []byte
	Hash  []byte
}

func LoadTestVectors(t *testing.T) []TestVector {
	t.Helper()
	b, err := ioutil.ReadFile("testdata/testvectors.json")
	if err != nil {
		t.Fatal(err)
	}

	var data []struct {
		SeedLo   uint64 `json:"seed_lo"`
		SeedHi   uint64 `json:"seed_hi"`
		InputHex string `json:"input_hex"`
		HashHex  string `json:"hash_hex"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatal(err)
	}

	var vs []TestVector
	for _, d := range data {
		in, err := hex.DecodeString(d.InputHex)
		if err != nil {
			t.Fatal(err)
		}

		h, err := hex.DecodeString(d.HashHex)
		if err != nil {
			t.Fatal(err)
		}

		vs = append(vs, TestVector{
			Seed:  (d.SeedHi << 32) | d.SeedLo,
			Input: in,
			Hash:  h,
		})
	}

	return vs
}

func TestVectors(t *testing.T) {
	vs := LoadTestVectors(t)
	for _, v := range vs {
		got := meow.Checksum(v.Seed, v.Input)
		AssertBytesEqual(t, v.Hash, got[:])
	}
}

func AssertBytesEqual(t *testing.T, expect, got []byte) {
	if len(expect) != len(got) {
		t.Fatalf("length mismatch got=%d expect=%d", len(expect), len(got))
	}

	if bytes.Equal(expect, got) {
		return
	}

	n := len(expect)
	delta := make([]byte, n)
	for i := 0; i < len(expect); i++ {
		delta[i] = expect[i] ^ got[i]
	}

	t.Fatalf("expected equal\n   got=%x\nexpect=%x\n delta=%x", got, expect, delta)
}
