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

$(progs): meow_hash.h meow_intrinsics.h

meow_%.h:
	wget https://raw.githubusercontent.com/cmuratori/meow_hash/master/$@

testvectors.json: testvectors
	./$< | tee $@ | python -m json.tool > /dev/null # pipe through python to catch a json mistake

clean:
	$(RM) $(progs) testvectors.json

distclean: clean
	$(RM) meow_hash.h