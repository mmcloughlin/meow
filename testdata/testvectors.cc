#include <stdio.h>
#include <inttypes.h>
#include <stdlib.h>
#include <assert.h>
#include <immintrin.h>
#include <string.h>
#include <getopt.h>

#include "meow_intrinsics.h"
#include "meow_hash.h"

// Generate random uint64 value.
uint64_t rand_uint64()
{
    uint64_t x = 0;
    x = (x << 16) | (rand() & 0xffff);
    x = (x << 16) | (rand() & 0xffff);
    x = (x << 16) | (rand() & 0xffff);
    x = (x << 16) | (rand() & 0xffff);
    return x;
}

// Generate a random byte array of the given length.
uint8_t *rand_uint8_array(size_t len)
{
    uint8_t *x = (uint8_t *)malloc(len);
    assert(x);
    for (size_t i = 0; i < len; i++)
    {
        x[i] = rand() & 0xff;
    }
    return x;
}

// Print array as hex.
void printx(uint8_t *x, size_t len)
{
    for (size_t i = 0; i < len; i++)
    {
        printf("%02x", x[i]);
    }
}

// Meow hash test vector.
typedef struct
{
    uint64_t seed;
    size_t len;
    uint8_t *input;
#define HASH_LEN (128 / 8)
    uint8_t hash[HASH_LEN];
    uint64_t hash64;
    uint32_t hash32;
} test_vector_t;

// Generate a random test vector. Allocates memory which must be free'd with test_vector_free.
test_vector_t *test_vector_rand(size_t len)
{
    test_vector_t *tv = (test_vector_t *)malloc(sizeof(test_vector_t));
    assert(tv);
    tv->seed = rand_uint64();
    tv->len = len;
    tv->input = rand_uint8_array(tv->len);
    meow_u128 hash = MeowHash1(tv->seed, tv->len, tv->input);
    memcpy(tv->hash, &hash, HASH_LEN);
    tv->hash64 = MeowU64From(hash);
    tv->hash32 = MeowU32From(hash);
    return tv;
}

// Write test vector as JSON to stdout.
void test_vector_json(test_vector_t *tv, const char *prefix, const char *indent)
{
    printf("%s{\n", prefix);
    printf("%s%s\"seed\": \"%016" PRIx64 "\",\n", prefix, indent, tv->seed);

    printf("%s%s\"input\": \"", prefix, indent);
    printx(tv->input, tv->len);
    printf("\",\n");

    printf("%s%s\"hash\": \"", prefix, indent);
    printx(tv->hash, HASH_LEN);
    printf("\",\n");

    printf("%s%s\"hash64\": \"%016" PRIx64 "\",\n", prefix, indent, tv->hash64);
    printf("%s%s\"hash32\": \"%08" PRIx32 "\"\n", prefix, indent, tv->hash32);

    printf("%s}", prefix);
}

// Free memory allocated to test_vector.
void test_vector_free(test_vector_t *tv)
{
    free(tv->input);
    free(tv);
}

// Output test vectors for a given range of lengths.
void output_test_vectors(size_t *lengths, size_t n)
{
    printf("{\n");
    printf("\t\"version_number\": %d,\n", MEOW_HASH_VERSION);
    printf("\t\"version_name\": \"%s\",\n", MEOW_HASH_VERSION_NAME);
    printf("\t\"test_vectors\": [\n");
    for (size_t i = 0; i < n; i++)
    {
        test_vector_t *tv = test_vector_rand(lengths[i]);
        test_vector_json(tv, "\t\t", "\t");
        test_vector_free(tv);
        printf(i != n - 1 ? ",\n" : "\n");
    }
    printf("\t]\n");
    printf("}\n");
}

// Populate array with lengths of the form (a * i)%m for i < n
void modulo_lengths(size_t *lengths, size_t a, size_t m, size_t n)
{
    for (size_t i = 0; i < n; i++)
    {
        lengths[i] = (i * a) % m;
    }
}

// Generate collection of test vector lengths.
size_t *generate_lengths(size_t *n)
{
    size_t num_small = 32;
    size_t num_modulo = 256;

    *n = num_small + num_modulo;
    size_t *lengths = (size_t *)malloc(*n * sizeof(size_t));
    assert(lengths);

    // Ensure we have examples of small hash inputs.
    for (size_t i = 0; i < num_small; i++)
    {
        lengths[i] = i;
    }

    // Modulo lengths provides larger hash inputs of all possible values modulo 256.
    modulo_lengths(lengths + num_small, 251, 8 << 10, num_modulo);

    return lengths;
}

int main(int argc, char **argv)
{
    int opt;
    while ((opt = getopt(argc, argv, "s:")) != -1)
    {
        switch (opt)
        {
        case 's':
            srand(atoi(optarg));
            break;
        }
    }

    size_t n;
    size_t *lengths = generate_lengths(&n);
    output_test_vectors(lengths, n);
    free(lengths);
}
