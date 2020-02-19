#include <string.h>
#include <stdbool.h>
#include "unittest.h"
#include "request.h"
#include "htable.h"

char *test_request_parse()
{
    char *buf = ""
        "GET / HTTP/1.1\r\n"
        "Host: localhost:9090\r\n"
        "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:65.0) Gecko/20100101 Firefox/65.0\r\n"
        "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\n"
        "Accept-Language: en-US,en;q=0.5\r\n"
        "Accept-Encoding: gzip, deflate\r\n"
        "Connection: keep-alive\r\n"
        "Upgrade-Insecure-Requests: 1\r\n"
        "Pragma: no-cache\r\n"
        "Cache-Control: no-cache\r\n\r\n"
        "\r\n";

    request *r = request_parse(buf);
    assert(strcmp(r->method, "GET") == 0, "Request Method should be GET");
    assert(strcmp(r->uri, "/") == 0, "Request Uri should be /");
    assert(strcmp((char*)htable_get(r->headers, "Host"), "localhost:9090") == 0, "Header \"Host\" not parsed correctly");
    assert(strcmp((char*)htable_get(r->headers, "Accept-Language"), "en-US,en;q=0.5") == 0, "Header \"Accept-Language\" not parsed correctly");
    assert(strcmp((char*)htable_get(r->headers, "Accept-Encoding"), "gzip, deflate") == 0, "Header \"Accept-Encoding\" not parsed correctly");
    assert(strcmp((char*)htable_get(r->headers, "Connection"), "keep-alive") == 0, "Header \"Connection\" not parsed correctly");
    assert(strcmp((char*)htable_get(r->headers, "Upgrade-Insecure-Requests"), "1") == 0, "Header \"Upgrade-Insecure-Requests\" not parsed correctly");
    assert(strcmp((char*)htable_get(r->headers, "Pragma"), "no-cache") == 0, "Header \"Pragma\" not parsed correctly");
    assert(strcmp((char*)htable_get(r->headers, "Cache-Control"), "no-cache") == 0, "Header \"Cache-Control\" not parsed correctly");
    request_destroy(r);

    return NULL;
}

int main(void)
{
    start_tests("request tests");
    run_test(test_request_parse);
    end_tests();

    return 0;
}
