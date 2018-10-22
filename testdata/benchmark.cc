#include <vector>

#include <benchmark/benchmark.h>

#include <immintrin.h>
#include "meow_hash.h"

static void meow_hash_benchmark(benchmark::State& state) {
    size_t            buffer_size = state.range(0);
    std::vector<char> buffer(buffer_size, 0);

    for (auto _ : state)
    {
      benchmark::DoNotOptimize(MeowHash1(0, buffer_size, buffer.data()));
    }

    state.SetBytesProcessed(state.iterations() * buffer_size);
}

BENCHMARK(meow_hash_benchmark)->RangeMultiplier(2)->Range(1, 1<<24);

BENCHMARK_MAIN();
