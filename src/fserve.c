#include "fserve.h"

int bind_and_listen(int port);
void set_default_response_headers(response *res);
void *create_request_buffer(int conn_fd);
void serve_404(response *res);

int main(int argc, char *argv[])
{
    int server_fd, conn_fd;
    int port = DEFAULT_PORT;
    char *dir = NULL;

    for (int i = 1; i < argc - 1; i++) {
        if (strcmp("-p", argv[i]) == 0) {
            port = atoi(argv[i+1]);
        }

        if (strcmp("-d", argv[i]) == 0) {
            dir = strdup(argv[i+1]);
        }
    }

    request *req = NULL;
    response *res = NULL;

    if ((server_fd = bind_and_listen(port)) < 0) {
        exit(1);
    }

    printf("Server running on port %d", port);
    if (dir) {
        printf(" and path '%s'", dir);
    }
    printf("\n");

    while (true) {
        char *req_buff = NULL;

        if ((conn_fd = accept(server_fd, (struct sockaddr*) NULL, NULL)) < 0) {
            perror("server: accept()");
            continue;
        }

        if ((req_buff = create_request_buffer(conn_fd)) == NULL) {
            close(conn_fd);
            continue;
        }

        // create the request object with the request buffer
        if ((req = request_create(req_buff)) == NULL) {
            close(conn_fd);
            free(req_buff);
            continue;
        }
        free(req_buff);

        // initialize the response object
        if ((res = response_create()) == NULL) {
            close(conn_fd);
            continue;
        }

        bool file_found = false;
        file_found = static_file_serve(req, res, dir);

        if (!file_found) {
            serve_404(res);
        }
        set_default_response_headers(res);

        char *http_status_line = response_get_start_line(res);
        char *response_headers = response_get_headers(res);
        char *response_body = response_get_body(res);

        if (send(conn_fd, http_status_line, strlen(http_status_line), 0) < 0) {
            perror("server: send() http status line");
            close(conn_fd);
            continue;
        }
        if (send(conn_fd, response_headers, strlen(response_headers), 0) < 0) {
            perror("server: send() http response headers");
            close(conn_fd);
            continue;
        }
        if (send(conn_fd, response_body, res->body_len, 0) < 0) {
            perror("server: send() http response body");
            close(conn_fd);
            continue;
        }
        close(conn_fd);
        request_destroy(req);
        response_destroy(res);
    }
    close(server_fd);
    free(dir);
    return 0;
}

void *create_request_buffer(int conn_fd)
{
    int bytes_r = 0;
    char *recv_buff = NULL;

    recv_buff = malloc(RECV_BUFF_LEN);
    if (recv_buff == NULL) {
        fprintf(stderr, "no memory!\n");
        return NULL;
    }
    memset(recv_buff, 0, RECV_BUFF_LEN);

    bytes_r = recv(conn_fd, recv_buff, RECV_BUFF_LEN, 0);
    if (bytes_r < 0) {
        perror("Error: recv()");
        free(recv_buff);
        return NULL;
    }
    recv_buff[bytes_r-1] = '\0';

    return recv_buff;
}

int bind_and_listen(int port)
{
    int server_fd;
    struct sockaddr_in server_addr;

    // TODO: Use getaddrinfo() to get information
    // on the ip addresses and whatnot

    if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
        perror("server: socket()");
        return -1;
    }

    memset(&server_addr, 0, sizeof(struct sockaddr_in));

    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(port);

    // re-use the address when stopping and restarting the server
    int reuse_addr = 1;
    if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &reuse_addr, sizeof(int)) < 0) {
        perror("server: setsockopt()");
        close(server_fd);
        return -1;
    }

    if (bind(server_fd, (struct sockaddr*) &server_addr, sizeof(struct sockaddr_in)) < 0) {
        perror("server: bind()");
        close(server_fd);
        return -1;
    }

    if (listen(server_fd, SOMAXCONN) < 0) {
        perror("server: listen()");
        close(server_fd);
        return -1;
    }

    return server_fd;
}

void set_default_response_headers(response *res)
{
    response_set_header(res, "Server", strdup("fserve v0.0.1"));
    response_set_header(res, "Connection", strdup("close"));
}

void serve_404(response *res)
{
    char *not_found_msg = "Not found\n";
    char not_found_len_str[3] = {0, 0, 0};
    snprintf(not_found_len_str, 3, "%d", (int)strlen(not_found_msg));
    response_set_status_code(res, 404);
    response_set_body(res, strdup(not_found_msg));
    response_set_body_len(res, strlen(not_found_msg));
    response_set_header(res, "Content-Length", strdup(not_found_len_str));
    response_set_header(res, "Content-Type", strdup("text/html"));
}
