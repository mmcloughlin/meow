CXX=clang++
CXXFLAGS=-O3 -march=native -std=c++11 -I$(GOOGLE_BENCHMARK_DIR)/include
LDFLAGS=-L$(GOOGLE_BENCHMARK_DIR)/lib
LDLIBS=-lbenchmark

all: testvectors.json benchmark

testvectors: testvectors.cc

testvectors.json: testvectors
	./$< > $@

meow_hash.h:
	wget https://raw.githubusercontent.com/cmuratori/meow_hash/master/meow_hash.h

clean:
	$(RM) testvectors testvectors.json