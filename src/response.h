#ifndef _RESPONSE_H
#define _RESPONSE_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "h_table.h"

#define HTTP_VERSION 1.1

typedef struct http_status_code {
    int code;
    char *description;
} http_status_code_t;

typedef struct {
    int status_code;
    char *status_message;
    h_table *headers;
    char *raw_headers;
    char *body;
    int body_len;
} response;

response *response_create(void);
void response_set_status_code(response *res, int code);
char *response_get_start_line(response *res);
void response_set_header(response *res, char *key, char *value);
char *response_get_headers(response *res);
void response_set_body(response *res, char *body);
void response_set_body_len(response *res, size_t body_len);
char *response_get_body(response *res);
void response_destroy(response *res);

#endif
