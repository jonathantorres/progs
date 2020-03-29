#ifndef _array_h
#define _array_h

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct array {
    unsigned int len;
    unsigned int capacity;
    unsigned int expand_rate;
    size_t item_size;
    void **contents;
} array;

typedef void(*array_free_cb)(void *value);

array *array_new(unsigned int capacity, size_t item_size);
void array_free(array *_array, array_free_cb cb);
void array_clear(array *_array, array_free_cb cb);
void array_push(array *_array, void *value);
void *array_pop(array *_array);
void array_set(array *_array, void *elem, unsigned int index);
void *array_get(array *_array, unsigned int index);
void *array_remove(array *_array, unsigned int index);
void array_shift(array *_array, void *value);
void *array_unshift(array *_array);
void array_swap(array *_array, unsigned int a, unsigned int b);

// Macro usage:
// ARRAY_FOREACH(array) {
    // your code here
    // you can use the variable "i" which contains the current index
    // you'll have to assign the current element to a variable to use it
// }
#define ARRAY_FOREACH(array) for (unsigned int i = 0; i < (array)->len; i++)

#endif
