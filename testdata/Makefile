CXX=clang++
CXXFLAGS=-O3 -maes -std=c++11

ifdef GOOGLE_BENCHMARK_DIR
	CXXFLAGS+=-O3 -maes -std=c++11 -I$(GOOGLE_BENCHMARK_DIR)/include
	LDFLAGS=-L$(GOOGLE_BENCHMARK_DIR)/lib
	LDLIBS=-lbenchmark
endif

progs=benchmark testvectors

all: $(progs) testvectors.json

%: %.cc
	$(CXX) -o $@ $(CXXFLAGS) $< $(LDFLAGS) $(LDLIBS)

$(progs): meow_hash.h

meow_hash.h:
	wget https://raw.githubusercontent.com/cmuratori/meow_hash/master/meow_hash.h

testvectors.json: testvectors
	./$< > $@

clean:
	$(RM) $(progs) testvectors.json