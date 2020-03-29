#include "h_table.h"

#define NUM_OF_BUCKETS 100

size_t _h_table_fnv1a_hash(const char* cp)
{
    size_t hash = 0x811c9dc5;
    while (*cp) {
        hash ^= (unsigned char) *cp++;
        hash *= 0x01000193;
    }

    return hash;
}

bool _h_table_bucket_index_is_valid(unsigned int bucket_index)
{
    if (bucket_index >= NUM_OF_BUCKETS) {
        return false;
    }

    return true;
}

array *_h_table_find_bucket(h_table *_h_table, char *key, size_t *bucket_hash, bool create_new)
{
    *bucket_hash = _h_table_fnv1a_hash(key);
    unsigned int bucket_index = *bucket_hash % NUM_OF_BUCKETS;

    if (!_h_table_bucket_index_is_valid(bucket_index)) {
        return NULL;
    }

    array *bucket = array_get(_h_table->buckets, bucket_index);
    if (!bucket && create_new) {
        bucket = array_new(NUM_OF_BUCKETS, sizeof(h_table_node*));
        array_set(_h_table->buckets, bucket, bucket_index);
    }

    return bucket;
}

h_table *h_table_new(h_table_compare cmp)
{
    h_table *_h_table = malloc(sizeof(h_table));

    if (!_h_table) {
        fputs("[h_table_new] Not enough memory.", stderr);
        return NULL;
    }

    _h_table->cmp = cmp;
    _h_table->buckets = array_new(NUM_OF_BUCKETS, sizeof(array*));
    _h_table->buckets->len = _h_table->buckets->capacity;

    return _h_table;
}

void h_table_free(h_table *_h_table, h_table_node_cb cb)
{
    if (!_h_table) {
        fputs("[h_table_free] Must provide a hash_table.", stderr);
        return;
    }

    if (_h_table->buckets) {
        for (unsigned int i = 0; i < _h_table->buckets->len; i++) {
            array *bucket = array_get(_h_table->buckets, i);
            if (bucket) {
                for (unsigned int j = 0; j < bucket->len; j++) {
                    h_table_node *elem = array_get(bucket, j);
                    if (elem) {
                        if (cb) {
                            cb(elem);
                        }
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
        if (_h_table->buckets->contents) {
            free(_h_table->buckets->contents);
        }
        free(_h_table->buckets);
    }
    free(_h_table);
}

void h_table_set(h_table *_h_table, char *key, void *value)
{
    if (!_h_table) {
        fputs("[h_table_set] Must provide a hash_table.", stderr);
        return;
    }

    size_t bucket_hash = 0;
    array *bucket = _h_table_find_bucket(_h_table, key, &bucket_hash, true);

    if (!bucket) {
        return;
    }

    // create the node
    h_table_node *node = malloc(sizeof(h_table_node));
    if (!node) {
        return;
    }

    node->key = strdup(key);
    node->value = value;
    node->hash = bucket_hash;
    array_push(bucket, node);
}

void *h_table_get(h_table *_h_table, char *key)
{
    if (!_h_table) {
        fputs("[h_table_get] Must provide a hash_table.", stderr);
        return NULL;
    }

    size_t bucket_hash = 0;
    array *bucket = _h_table_find_bucket(_h_table, key, &bucket_hash, false);
    if (!bucket) {
        return NULL;
    }

    for (unsigned int i = 0; i < bucket->len; i++) {
        h_table_node *node = array_get(bucket, i);
        if (node && node->hash == bucket_hash && _h_table->cmp(node->key, key) == 0) {
            return node->value;
        }
    }

    return NULL;
}

void *h_table_remove(h_table *_h_table, char *key, h_table_node_cb cb)
{
    if (!_h_table) {
        fputs("[h_table_remove] Must provide a hash_table.", stderr);
        return NULL;
    }

    void *value = NULL;
    size_t bucket_hash = 0;
    array *bucket = _h_table_find_bucket(_h_table, key, &bucket_hash, false);
    if (!bucket) {
        return value;
    }

    for (unsigned int i = 0; i < bucket->len; i++) {
        h_table_node *node = array_get(bucket, i);
        if (node && node->hash == bucket_hash && _h_table->cmp(node->key, key) == 0) {
            value = node->value;
            void *last = array_pop(bucket);
            if (last != node) {
                array_set(bucket, last, i);
            }
            if (cb) {
                cb(node);
            }
            free(node->key);
            free(node);
            break;
        }
    }

    return value;
}

void h_table_traverse(h_table *_h_table, h_table_node_cb cb)
{
    if (!_h_table) {
        fputs("[h_table_traverse] Must provide a hash_table.", stderr);
        return;
    }

    for (unsigned int i = 0; i < _h_table->buckets->len; i++) {
        array *bucket = array_get(_h_table->buckets, i);
        if (bucket) {
            for (unsigned int j = 0; j < bucket->len; j++) {
                h_table_node *node = array_get(bucket, j);
                if (node) {
                    if (cb) {
                        cb(node);
                    }
                }
            }
        }
    }
}
