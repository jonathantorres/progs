#ifndef _dl_list_h
#define _dl_list_h

#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

typedef struct dl_list_node {
    void *value;
    struct dl_list_node *next;
    struct dl_list_node *prev;
} dl_list_node;

typedef struct dl_list {
    dl_list_node *first;
} dl_list;

typedef int(*dl_list_cmp)(void *a, void *b);
typedef void(*dl_list_free_cb)(void *value);

dl_list *dl_list_new();
void dl_list_clear(dl_list *list, dl_list_free_cb cb);
void dl_list_free(dl_list *list, dl_list_free_cb cb);
int dl_list_len(dl_list *list);
void dl_list_push(dl_list *list, void *value);
void dl_list_shift(dl_list *list, void *value);
void *dl_list_unshift(dl_list *list);
void *dl_list_pop(dl_list *list);
void dl_list_remove(dl_list *list, void *value, dl_list_cmp cmp, dl_list_free_cb cb);
bool dl_list_exists(dl_list *list, void *value, dl_list_cmp cmp);

// Macro usage:
// DL_LIST_FOREACH(list) {
    // your code here
    // you can use the variable "cur"
    // inside of it that references the current item of the list
// }
#define DL_LIST_FOREACH(list) \
    for (dl_list_node *cur = (list)->first; cur != NULL; cur = cur->next)

#endif
