#include <stdio.h>
#include <inttypes.h>
#include <stdlib.h>
#include <assert.h>
#include <immintrin.h>
#include <string.h>
#include <getopt.h>

#include "meow_hash.h"

// Generate random uint64 value.
uint64_t rand_uint64()
{
    uint64_t x = 0;
    x = (x << 16) | (rand() & 0xff);
    x = (x << 16) | (rand() & 0xff);
    x = (x << 16) | (rand() & 0xff);
    x = (x << 16) | (rand() & 0xff);
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
#define HASH_LEN (512 / 8)
    uint8_t hash[HASH_LEN];
} test_vector_t;

// Generate a random test vector. Allocates memory which must be free'd with test_vector_free.
test_vector_t *test_vector_rand(size_t len)
{
    test_vector_t *tv = (test_vector_t *)malloc(sizeof(test_vector_t));
    assert(tv);
    tv->seed = rand_uint64();
    tv->len = len;
    tv->input = rand_uint8_array(tv->len);
    meow_lane lane = MeowHash1(tv->seed, tv->len, tv->input);
    memcpy(tv->hash, &lane.Sub[0], HASH_LEN);
    return tv;
}

// Write test vector as JSON to stdout.
void test_vector_json(test_vector_t *tv, const char *prefix, const char *indent)
{
    printf("%s{\n", prefix);
    printf("%s%s\"seed_lo\": %" PRId64 ",\n", prefix, indent, tv->seed & 0xffffffff);
    printf("%s%s\"seed_hi\": %" PRId64 ",\n", prefix, indent, (tv->seed >> 32) & 0xffffffff);

    printf("%s%s\"input_hex\": \"", prefix, indent);
    printx(tv->input, tv->len);
    printf("\",\n");

    printf("%s%s\"hash_hex\": \"", prefix, indent);
    printx(tv->hash, HASH_LEN);
    printf("\"\n");

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
    test_vector_t *tv;
    printf("[\n");
    for (size_t i = 0; i < n; i++)
    {
        tv = test_vector_rand(lengths[i]);
        test_vector_json(tv, "    ", "    ");
        test_vector_free(tv);
        printf(i != n - 1 ? ",\n" : "\n");
    }
    printf("]\n");
}

// Generate lengths populates an array with lengths of the form (a*i)%m for i <= n.
size_t *generate_lengths(size_t a, size_t m, size_t n)
{
    size_t *lengths = (size_t *)malloc(n * sizeof(size_t));
    assert(lengths);
    for (size_t i = 0; i < n; i++)
    {
        lengths[i] = (i * a) % m;
    }
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

    size_t n = 256;
    size_t *lengths = generate_lengths(251, 8 << 10, n);
    output_test_vectors(lengths, n);
    free(lengths);
}