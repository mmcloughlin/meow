package meow_test

import (
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
		got := meow.Sum(v.Seed, v.Input)
		t.Logf("got=%x expect=%x", got, v.Hash)
	}
}
