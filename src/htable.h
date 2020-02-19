#ifndef _htable_h
#define _htable_h

#include "array.h"
#include <stdbool.h>
#include <string.h>

#define NUM_OF_BUCKETS 100

typedef int (*htable_compare)(void *a, void *b);

typedef struct htable {
    array *buckets;
    htable_compare cmp;
} htable;

typedef struct htable_node {
    void *key;
    void *value;
    size_t hash;
} htable_node;

typedef bool (*htable_node_cb)(htable_node *node);

htable *htable_create(htable_compare cmp);
void htable_destroy(htable *_htable);
void *htable_get(htable *_htable, void *key);
bool htable_set(htable *_htable, void *key, void *value);
void *htable_remove(htable *_htable, void *key);
bool htable_traverse(htable *_htable, htable_node_cb cb);

// Macro Usage:
// HTABLE_FOREACH(htable) {
    // use elem
    // elem is the current htable_node
// } HTABLE_FOREACH_END
#define HTABLE_FOREACH(htable) \
    if ((htable)->buckets) { \
        for (unsigned int i = 0; i < (htable)->buckets->length; i++) { \
            array *bucket = array_get((htable)->buckets, i); \
            if (bucket) { \
                for (unsigned int j = 0; j < bucket->length; j++) { \
                    htable_node *elem = array_get(bucket, j); \
                    if (elem)

#define HTABLE_FOREACH_END } } } }

#endif
