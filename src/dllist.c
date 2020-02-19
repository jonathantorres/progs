#include "dllist.h"
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

dllist_node *create_node(void *value)
{
    dllist_node *node = malloc(sizeof(dllist_node));

    if (!node) {
        fputs("Not enough memory.", stderr);
        return NULL;
    }

    node->prev = NULL;
    node->next = NULL;
    node->value = value;

    return node;
}

void destroy_node(dllist_node *node)
{
    if (!node) {
        fputs("A valid node must be provided.", stderr);
        return;
    }

    node->prev = NULL;
    node->next = NULL;
    node->value = NULL;

    free(node);
}

// create a new list
dllist *dllist_new()
{
    dllist *new_list = malloc(sizeof(dllist));

    if (!new_list) {
        fputs("Not enough memory.", stderr);
        return NULL;
    }

    new_list->first = NULL;

    return new_list;
}

// remove all the values in the list
void dllist_clear(dllist *list)
{
    if (list->first != NULL) {
        dllist_node *current_node = list->first;

        while (current_node->next != NULL) {
            current_node = current_node->next;

            if (current_node->prev) {
                destroy_node(current_node->prev);
            }
        }

        if (current_node) {
            destroy_node(current_node);
        }

        list->first = NULL;
    }
}

// destroy the list
void dllist_destroy(dllist *list)
{
    if (!list) {
        fputs("Must provide a dllist.", stderr);
        return;
    }

    dllist_clear(list);
    free(list);
}

// get the length of the list
int dllist_length(dllist *list)
{
    if (!list) {
        fputs("Must provide a dllist.", stderr);
        return -1;
    }

    int length = 0;

    if (list->first != NULL) {
        dllist_node *current_node = list->first;
        length++;

        while (current_node->next != NULL) {
            current_node = current_node->next;
            length++;
        }
    }

    return length;
}

// insert at the end
void dllist_push(dllist *list, void *value)
{
    if (!list) {
        fputs("Must provide a dllist.", stderr);
        return;
    }

    dllist_node *new_node = create_node(value);

    // list is empty, this is the first element
    if (list->first == NULL) {
        list->first = new_node;
        return;
    }

    dllist_node *current_node = list->first;

    while (current_node->next != NULL) {
        current_node = current_node->next;
    }

    new_node->prev = current_node;
    current_node->next = new_node;
}

// insert at the beginning
void dllist_shift(dllist *list, void *value)
{
    if (!list) {
        fputs("Must provide a dllist.", stderr);
        return;
    }

    dllist_node *new_node = create_node(value);

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
void *dllist_unshift(dllist *list)
{
    if (!list) {
        fputs("Must provide a dllist.", stderr);
        return NULL;
    }

    // list is empty, return nothing
    if (list->first == NULL) {
        return NULL;
    }

    // list has just 1 node
    if (list->first->next == NULL) {
        void *value = list->first->value;
        destroy_node(list->first);

        return value;
    }

    dllist_node *new_first = list->first->next;
    void *value = list->first->value;

    destroy_node(list->first);
    new_first->prev = NULL;
    list->first = new_first;

    return value;
}

// remove the last node and return it
void *dllist_pop(dllist *list)
{
    if (!list) {
        fputs("Must provide a dllist.", stderr);
        return NULL;
    }

    // list is empty, return nothing
    if (list->first == NULL) {
        return NULL;
    }

    // list has just 1 node
    if (list->first->next == NULL) {
        void *value = list->first->value;
        destroy_node(list->first);
        list->first = NULL;

        return value;
    }

    dllist_node *current_node = list->first;

    while (current_node->next != NULL) {
        current_node = current_node->next;
    }

    void *value = current_node->value;
    current_node->prev->next = NULL;
    destroy_node(current_node);

    return value;
}

// remove node whose value is {value}
void dllist_remove(dllist *list, void *value, dllist_cmp cmp)
{
    if (!list) {
        fputs("Must provide a valid dllist.", stderr);
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
            dllist_pop(list);
        }

        return;
    }

    dllist_node *current_node = list->first;

    // check the first one
    if (cmp(current_node->value, value) == 0) {
        dllist_node *next_node = current_node->next;
        next_node->prev = NULL;
        list->first = next_node;
        destroy_node(current_node);
        return;
    }

    while (current_node->next != NULL) {
        current_node = current_node->next;

        if (cmp(current_node->value, value) == 0) {
            // remove the node
            current_node->prev->next = current_node->next;
            current_node->next->prev = current_node->prev;
            destroy_node(current_node);
            break;
        }
    }
}

// check to see if value {value} exists in the list
bool dllist_exists(dllist *list, void *value, dllist_cmp cmp)
{
    if (!list) {
        fputs("Must provide a valid dllist.", stderr);
        return -1;
    }

    // list is empty, not found
    if (list->first == NULL) {
        return false;
    }

    dllist_node *current_node = list->first;

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
