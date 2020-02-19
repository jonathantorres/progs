#include <stdlib.h>
#include <stdio.h>
#include "htable.h"

size_t fnv1a_hash(const char* cp)
{
    size_t hash = 0x811c9dc5;
    while (*cp) {
        hash ^= (unsigned char) *cp++;
        hash *= 0x01000193;
    }

    return hash;
}

htable *htable_create(htable_compare cmp)
{
    htable *_htable = malloc(sizeof(htable));

    if (!_htable) {
        fputs("[htable_create] Not enough memory.", stderr);
        return NULL;
    }

    _htable->cmp = cmp;
    _htable->buckets = array_create(NUM_OF_BUCKETS, sizeof(array*));
    _htable->buckets->length = _htable->buckets->capacity;

    return _htable;
}

void htable_destroy(htable *_htable)
{
    if (!_htable) {
        fputs("[htable_destroy] Must provide a hashtable.", stderr);
        return;
    }

    if (_htable->buckets) {
        for (unsigned int i = 0; i < _htable->buckets->length; i++) {
            array *bucket = array_get(_htable->buckets, i);
            if (bucket) {
                for (unsigned int j = 0; j < bucket->length; j++) {
                    htable_node *elem = array_get(bucket, j);
                    if (elem) {
                        if (elem->key) {
                            free(elem->key);
                        }
                        free(elem);
                    }
                }
                if (bucket->contents) {
                    free(bucket->contents);
                }
                free(bucket);
            }
        }
        if (_htable->buckets->contents) {
            free(_htable->buckets->contents);
        }
        free(_htable->buckets);
    }
    free(_htable);
}

bool bucket_index_is_valid(unsigned int bucket_index)
{
    if (bucket_index >= NUM_OF_BUCKETS) {
        return false;
    }

    return true;
}

array *find_bucket(htable *_htable, void *key, size_t *bucket_hash, bool create_new)
{
    *bucket_hash = fnv1a_hash(key);
    unsigned int bucket_index = *bucket_hash % NUM_OF_BUCKETS;

    if (!bucket_index_is_valid(bucket_index)) {
        return NULL;
    }

    array *bucket = array_get(_htable->buckets, bucket_index);
    if (!bucket && create_new) {
        bucket = array_create(NUM_OF_BUCKETS, sizeof(htable_node*));
        array_set(_htable->buckets, bucket, bucket_index);
    }

    return bucket;
}

bool htable_set(htable *_htable, void *key, void *value)
{
    if (!_htable) {
        fputs("[htable_set] Must provide a hashtable.", stderr);
        return false;
    }

    size_t bucket_hash = 0;
    array *bucket = find_bucket(_htable, key, &bucket_hash, true);

    if (!bucket) {
        return false;
    }

    // create the node
    htable_node *node = malloc(sizeof(htable_node));
    if (!node) {
        return false;
    }

    node->key = strdup(key);
    node->value = value;
    node->hash = bucket_hash;
    array_push(bucket, node);

    return true;
}

void *htable_get(htable *_htable, void *key)
{
    if (!_htable) {
        fputs("[htable_get] Must provide a hashtable.", stderr);
        return NULL;
    }

    size_t bucket_hash = 0;
    array *bucket = find_bucket(_htable, key, &bucket_hash, false);
    if (!bucket) {
        return NULL;
    }

    for (unsigned int i = 0; i < bucket->length; i++) {
        htable_node *node = array_get(bucket, i);
        if (node && node->hash == bucket_hash && _htable->cmp(node->key, key) == 0) {
            return node->value;
        }
    }

    return NULL;
}

void *htable_remove(htable *_htable, void *key)
{
    if (!_htable) {
        fputs("[htable_remove] Must provide a hashtable.", stderr);
        return NULL;
    }

    void *value = NULL;
    size_t bucket_hash = 0;
    array *bucket = find_bucket(_htable, key, &bucket_hash, false);
    if (!bucket) {
        return value;
    }

    for (unsigned int i = 0; i < bucket->length; i++) {
        htable_node *node = array_get(bucket, i);
        if (node && node->hash == bucket_hash && _htable->cmp(node->key, key) == 0) {
            value = node->value;
            void *last = array_pop(bucket);
            if (last != node) {
                array_set(bucket, last, i);
            }
            free(node->key);
            free(node);
            break;
        }
    }

    return value;
}

bool htable_traverse(htable *_htable, htable_node_cb cb)
{
    if (!_htable) {
        fputs("[htable_traverse] Must provide a hashtable.", stderr);
        return false;
    }
    unsigned int traverse_ok = false;

    for (unsigned int i = 0; i < _htable->buckets->length; i++) {
        array *bucket = array_get(_htable->buckets, i);
        if (bucket) {
            for (unsigned int j = 0; j < bucket->length; j++) {
                htable_node *node = array_get(bucket, j);
                if (node) {
                    if (cb(node)) {
                        traverse_ok = true;
                        break;
                    }
                }
            }
        }
    }

    return traverse_ok;
}
