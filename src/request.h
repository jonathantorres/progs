#ifndef _REQUEST_H
#define _REQUEST_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <ctype.h>
#include <sys/select.h>
#include <sys/socket.h>
#include "array.h"
#include "htable.h"

#define REQUEST_METHOD_GET "GET"
#define REQUEST_METHOD_HEAD "HEAD"
#define REQUEST_METHOD_POST "POST"
#define REQUEST_METHOD_PUT "PUT"
#define REQUEST_METHOD_DELETE "DELETE"
#define REQUEST_METHOD_OPTIONS "OPTIONS"
#define REQUEST_METHOD_TRACE "TRACE"
#define REQUEST_METHOD_PATCH "PATCH"

typedef struct {
    char *method;
    char *uri;
    htable *headers;
    htable *uri_params;
    char *body;
} request;

request *request_create(char *buffer);
void request_destroy(request *req);
void request_get_method(); // GET or POST
void request_get_uri(); // / or /some or /some/another
void request_get_query(); // get's the foo=bar&bar=baz in a data structure
void request_get_headers();
void request_get_header(); // provide the header name?
void request_set_header();
void request_get_body(); // a string with the body

#endif
