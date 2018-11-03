package meow

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"hash"
	"io"
	"io/ioutil"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"testing/quick"
)

var testVectorsFilename = flag.String("testvectors", "testdata/testvectors.json", "test vectors filename")

type TestVector struct {
	Seed   uint64
	Input  []byte
	Hash   []byte
	Hash64 uint64
	Hash32 uint32
}

func (v *TestVector) UnmarshalJSON(b []byte) error {
	var data struct {
		Seed   string `json:"seed"`
		Input  string `json:"input"`
		Hash   string `json:"hash"`
		Hash64 string `json:"hash64"`
		Hash32 string `json:"hash32"`
	}

	var err error

	if err = json.Unmarshal(b, &data); err != nil {
		return err
	}

	if v.Seed, err = strconv.ParseUint(data.Seed, 16, 64); err != nil {
		return err
	}

	if v.Input, err = hex.DecodeString(data.Input); err != nil {
		return err
	}

	if v.Hash, err = hex.DecodeString(data.Hash); err != nil {
		return err
	}

	if v.Hash64, err = strconv.ParseUint(data.Hash64, 16, 64); err != nil {
		return err
	}

	h32, err := strconv.ParseUint(data.Hash32, 16, 32)
	if err != nil {
		return err
	}
	v.Hash32 = uint32(h32)

	return nil
}

type TestData struct {
	Version     int          `json:"version_number"`
	VersionName string       `json:"version_name"`
	TestVectors []TestVector `json:"test_vectors"`
}

var testdata *TestData

func LoadTestData(t *testing.T) *TestData {
	t.Helper()

	if testdata != nil {
		return testdata
	}

	b, err := ioutil.ReadFile(*testVectorsFilename)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(b, &testdata); err != nil {
		t.Fatal(err)
	}

	return testdata
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

	t.Fatalf("expected equal\n   got=%x\nexpect=%x\n delta=%x\n", got, expect, delta)
}

func Trials() int {
	if testing.Short() {
		return 1 << 7
	}
	return 1 << 14
}

func CheckEqual(t *testing.T, f, g checksumFunc) {
	cfg := &quick.Config{}
	cfg.MaxCount = Trials()
	cfg.Values = func(args []reflect.Value, r *rand.Rand) {
		n := r.Intn(8 << 10)
		b := make([]byte, n)
		r.Read(b)

		args[0] = reflect.ValueOf(r.Uint64())
		args[1] = reflect.ValueOf(b)
	}

	t.Logf("trials=%d", cfg.MaxCount)

	if err := quick.CheckEqual(f, g, cfg); err != nil {
		t.Fatal(err)
	}
}

func AssertHashSize(t *testing.T, name string, h hash.Hash, size int) {
	if h.Size() != size {
		t.Errorf("%s Size() got=%d expect=%d", name, h.Size(), size)
	}

	io.WriteString(h, "Confirm the hash Sum() is what you expect")
	n := len(h.Sum(nil))
	if n != size {
		t.Errorf("%s len(Sum()) got=%d expect=%d", name, n, size)
	}
}
