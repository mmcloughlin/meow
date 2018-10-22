Test data generation using the [official implementation](https://github.com/cmuratori/meow_hash).

* `testvectors.cc` generates a set of test vectors in JSON format
* `benchmark.cc` benchmarks `MeowHash1`

Note benchmarks depend on [Google benchmark](https://github.com/google/benchmark), and the `Makefile` expects to find the installation at `$GOOGLE_BENCHMARK_DIR`. On Mac with homebrew:

```sh
brew install google-benchmark
export GOOGLE_BENCHMARK_DIR=$(brew --prefix google-benchmark)
```