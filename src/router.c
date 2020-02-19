#include "router.h"

typedef void (*route_action)(request *req, response *res);

typedef struct {
    char *path;
    route_action action;
} route;

static dllist *routes = NULL;

void default_action(request *req, response *res)
{
    char *response = "<h1>Hello World!</h1>";
    response_set_status_code(res, 200);
    response_set_header(res, "action", strdup("default"));
    response_set_body(res, strdup(response));
    response_set_body_len(res, strlen(response));
}

void test_action(request *req, response *res)
{
    char *response = "<h1>This is the default action</h1>";
    response_set_status_code(res, 200);
    response_set_body(res, strdup(response));
    response_set_body_len(res, strlen(response));
}

route *create_route(char *path, route_action action)
{
    route *new_route = malloc(sizeof(route));

    if (!new_route) {
        return NULL;
    }

    new_route->path = path;
    new_route->action = action;

    return new_route;
}

void register_routes()
{
    routes = dllist_new();

    dllist_push(routes, create_route("/", default_action));
    dllist_push(routes, create_route("/test", test_action));
}

void destroy_routes()
{
    DLLIST_FOREACH(routes) {
        if (cur->value) {
            free(cur->value);
        }
    }

    dllist_destroy(routes);
}

bool router_handle_request(request *req, response *res)
{
    bool route_found = false;
    register_routes();

    DLLIST_FOREACH(routes) {
        if (strcmp(req->uri, ((route*)cur->value)->path) == 0) {
            ((route*)cur->value)->action(req, res);
            route_found = true;
            break;
        }
    }

    destroy_routes();

    return route_found;
}
