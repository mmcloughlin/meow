CXX=clang++
CXXFLAGS=-maes

all: testvectors.json

testvectors: testvectors.cc

testvectors.json: testvectors
	./$< > $@

meow_hash.h:
	wget https://raw.githubusercontent.com/cmuratori/meow_hash/master/meow_hash.h
