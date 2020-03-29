#include "dl_list.h"

dl_list_node *_dl_list_create_node(void *value)
{
    dl_list_node *node = malloc(sizeof(dl_list_node));

    if (!node) {
        fputs("Not enough memory.", stderr);
        return NULL;
    }

    node->prev = NULL;
    node->next = NULL;
    node->value = value;

    return node;
}

void _dl_list_free_node(dl_list_node *node, dl_list_free_cb cb)
{
    if (!node) {
        fputs("A valid node must be provided.", stderr);
        return;
    }
    if (cb) {
        cb(node->value);
    }

    node->prev = NULL;
    node->next = NULL;
    node->value = NULL;

    free(node);
}

// create a new list
dl_list *dl_list_new()
{
    dl_list *new_list = malloc(sizeof(dl_list));

    if (!new_list) {
        fputs("Not enough memory.", stderr);
        return NULL;
    }

    new_list->first = NULL;

    return new_list;
}

// remove all the values in the list
void dl_list_clear(dl_list *list, dl_list_free_cb cb)
{
    if (list->first != NULL) {
        dl_list_node *current_node = list->first;

        while (current_node->next != NULL) {
            current_node = current_node->next;

            if (current_node->prev) {
                _dl_list_free_node(current_node->prev, cb);
            }
        }

        if (current_node) {
            _dl_list_free_node(current_node, cb);
        }

        list->first = NULL;
    }
}

// destroy the list
void dl_list_free(dl_list *list, dl_list_free_cb cb)
{
    if (!list) {
        fputs("Must provide a dl_list.", stderr);
        return;
    }

    dl_list_clear(list, cb);
    free(list);
}

// get the len of the list
int dl_list_len(dl_list *list)
{
    if (!list) {
        fputs("Must provide a dl_list.", stderr);
        return -1;
    }

    int len = 0;

    if (list->first != NULL) {
        dl_list_node *current_node = list->first;
        len++;

        while (current_node->next != NULL) {
            current_node = current_node->next;
            len++;
        }
    }

    return len;
}

// insert at the end
void dl_list_push(dl_list *list, void *value)
{
    if (!list) {
        fputs("Must provide a dl_list.", stderr);
        return;
    }

    dl_list_node *new_node = _dl_list_create_node(value);

    // list is empty, this is the first element
    if (list->first == NULL) {
        list->first = new_node;
        return;
    }

    dl_list_node *current_node = list->first;

    while (current_node->next != NULL) {
        current_node = current_node->next;
    }

    new_node->prev = current_node;
    current_node->next = new_node;
}

// insert at the beginning
void dl_list_shift(dl_list *list, void *value)
{
    if (!list) {
        fputs("Must provide a dl_list.", stderr);
        return;
    }

    dl_list_node *new_node = _dl_list_create_node(value);

    // list is empty, this is the first element
    if (list->first == NULL) {
        list->first = new_node;
        return;
    }

    list->first->prev = new_node;
    new_node->next = list->first;
    list->first = new_node;
}

// remove the first node and return it
void *dl_list_unshift(dl_list *list)
{
    if (!list) {
        fputs("Must provide a dl_list.", stderr);
        return NULL;
    }

    // list is empty, return nothing
    if (list->first == NULL) {
        return NULL;
    }

    // list has just 1 node
    if (list->first->next == NULL) {
        void *value = list->first->value;
        _dl_list_free_node(list->first, NULL);
        list->first = NULL;

        return value;
    }

    dl_list_node *new_first = list->first->next;
    void *value = list->first->value;

    _dl_list_free_node(list->first, NULL);
    new_first->prev = NULL;
    list->first = new_first;

    return value;
}

// remove the last node and return it
void *dl_list_pop(dl_list *list)
{
    if (!list) {
        fputs("Must provide a dl_list.", stderr);
        return NULL;
    }

    // list is empty, return nothing
    if (list->first == NULL) {
        return NULL;
    }

    // list has just 1 node
    if (list->first->next == NULL) {
        void *value = list->first->value;
        _dl_list_free_node(list->first, NULL);
        list->first = NULL;

        return value;
    }

    dl_list_node *current_node = list->first;

    while (current_node->next != NULL) {
        current_node = current_node->next;
    }

    void *value = current_node->value;
    current_node->prev->next = NULL;
    _dl_list_free_node(current_node, NULL);

    return value;
}

// remove node whose value is {value}
void dl_list_remove(dl_list *list, void *value, dl_list_cmp cmp, dl_list_free_cb cb)
{
    if (!list) {
        fputs("Must provide a valid dl_list.", stderr);
        return;
    }

    // list is empty, return nothing
    if (list->first == NULL) {
        return;
    }

    // list has just 1 node
    if (list->first->next == NULL) {
        void *node_value = list->first->value;

        if (cmp(node_value, value) == 0) {
            _dl_list_free_node(list->first, cb);
        }

        return;
    }

    dl_list_node *current_node = list->first;

    // check the first one
    if (cmp(current_node->value, value) == 0) {
        dl_list_node *next_node = current_node->next;
        next_node->prev = NULL;
        list->first = next_node;
        _dl_list_free_node(current_node, cb);
        return;
    }

    while (current_node->next != NULL) {
        current_node = current_node->next;

        if (cmp(current_node->value, value) == 0) {
            // remove the node
            current_node->prev->next = current_node->next;
            current_node->next->prev = current_node->prev;
            _dl_list_free_node(current_node, cb);
            break;
        }
    }
}

// check to see if value {value} exists in the list
bool dl_list_exists(dl_list *list, void *value, dl_list_cmp cmp)
{
    if (!list) {
        fputs("Must provide a valid dl_list.", stderr);
        return -1;
    }

    // list is empty, not found
    if (list->first == NULL) {
        return false;
    }

    dl_list_node *current_node = list->first;

    // check the first one
    if (cmp(current_node->value, value) == 0) {
        return true;
    }

    bool found = false;

    // check the rest
    while (current_node->next != NULL) {
        current_node = current_node->next;

        if (cmp(current_node->value, value) == 0) {
            found = true;
            break;
        }
    }

    return found;
}
