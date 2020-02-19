#include "array.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define EXPAND_RATE 100

// creates a new empty array
array *array_create(unsigned int capacity, size_t item_size)
{
    array *_array = malloc(sizeof(array));

    if (!_array) {
        fputs("[array_create] Not enough memory.", stderr);
        return NULL;
    }

    _array->length = 0;
    _array->capacity = capacity;
    _array->expand_rate = EXPAND_RATE;
    _array->item_size = item_size;
    _array->contents = calloc(_array->capacity, _array->item_size);

    if (!_array->contents) {
        fputs("[array_create] Not enough memory.", stderr);
        return NULL;
    }

    return _array;
}

// empties and destroys the array completely
void array_destroy(array *_array)
{
    if (!_array) {
        fputs("[array_destroy] Must provide an array.", stderr);
        return;
    }

    array_clear(_array);
    free(_array->contents);
    free(_array);
}

// removes all the elements on the array, leaving it empty
void array_clear(array *_array)
{
    if (!_array) {
        fputs("[array_clear] Must provide an array.", stderr);
        return;
    }

    unsigned int array_length = _array->length;

    for (unsigned int i = 0; i < array_length; i++) {
        if (_array->contents[i]) {
            _array->contents[i] = NULL;
        }
        _array->length--;
    }
}

void array_expand(array *_array)
{
    int new_capacity = _array->capacity + EXPAND_RATE;
    void *contents = realloc(_array->contents, new_capacity * _array->item_size);

    if (!contents) {
        fputs("[array_expand] Not enough memory.", stderr);
        return;
    }

    _array->contents = contents;
    _array->capacity = new_capacity;
}

// add element to the end
void array_push(array *_array, void *value)
{
    if (!_array) {
        fputs("[array_push] Must provide an array.", stderr);
        return;
    }

    _array->contents[_array->length] = value;
    _array->length++;

    // expand if necessary
    if (_array->length >= _array->capacity) {
        array_expand(_array);
    }
}

void *remove_element_at(array *_array, unsigned int index)
{
    if (_array->contents[index] != NULL) {
        void *element = _array->contents[index];
        _array->contents[index] = NULL;
        _array->length--;

        return element;
    }

    return NULL;
}

// remove last element and return it
void *array_pop(array *_array)
{
    if (!_array) {
        fputs("[array_pop] Must provide an array.", stderr);
        return NULL;
    }

    void *element = NULL;
    if (_array->length > 0) {
        element = remove_element_at(_array, _array->length - 1);
    }

    return element;
}

// add/set element at index
void array_set(array *_array, void *elem, unsigned int index)
{
    if (!_array) {
        fputs("[array_set] Must provide an array.", stderr);
        return;
    }

    // index is too large
    if (index >= _array->capacity) {
        return;
    }

    if (index >= _array->length) {
        _array->length = index + 1;
    }

    _array->contents[index] = elem;
}

// get element at index
void *array_get(array *_array, unsigned int index)
{
    if (!_array) {
        fputs("[array_get] Must provide an array.", stderr);
        return NULL;
    }

    // index is too large
    if (index >= _array->length) {
        return NULL;
    }

    return _array->contents[index];
}

// remove element at index and return it
void *array_remove(array *_array, unsigned int index)
{
    if (!_array) {
        fputs("[array_remove] Must provide an array.", stderr);
        return NULL;
    }

    // index is too large
    if (index >= _array->length) {
        return NULL;
    }

    void *element = remove_element_at(_array, index);

    if (element != NULL && _array->contents[index + 1] != NULL) {
        memmove(
            &_array->contents[index],
            &_array->contents[index + 1],
            sizeof(_array->item_size) * (_array->length - index)
        );
    }

    return element;
}

// add element to the beginning
void array_shift(array *_array, void *value)
{
    if (!_array) {
        fputs("[array_shift] Must provide an array.", stderr);
        return;
    }

    if (_array->length > 0) {
        memmove(
            &_array->contents[1],
            _array->contents,
            sizeof(_array->item_size) * _array->length
        );
    }

    _array->contents[0] = value;
    _array->length++;

    // expand if necessary
    if (_array->length >= _array->capacity) {
        array_expand(_array);
    }
}

// remove first element and return it
void *array_unshift(array *_array)
{
    if (!_array) {
        fputs("[array_unshift] Must provide an array.", stderr);
        return NULL;
    }

    void *element = NULL;

    if (_array->length > 0) {
        element = remove_element_at(_array, 0);

        memmove(
            _array->contents,
            &_array->contents[1],
            sizeof(_array->item_size) * _array->length
        );
    }

    return element;
}
