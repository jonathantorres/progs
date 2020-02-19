#include <stdio.h>

// This is how I'd like to use this thing
// This won't compile!

void index_action(request *req, response *res)
{
    // add headers
    // get request data
    // set response data
    // return response string
}

int main(void)
{
    server *s = server_init();
    // use a router to define your routes
    router *r = router_new();
    request *req = request_init();
    response *res = response_init();

    // GET routes
    router_add_get_route(r, "/", index_action);
    router_add_get_route(r, "/contact", contact_action);
    router_add_get_route(r, "/user/{id}", user_action); // with a required param
    router_add_get_route(r, "/invoice/{id?}", invoice_action); // with a default param (uses a default value)

    // POST routes
    router_add_post_route(r, "/create", create_action);
    router_add_post_route(r, "/add/{id}", add_action);

    // run server
    s.run();
    return 0;
}
