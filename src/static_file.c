#include "static_file.h"

typedef struct {
    char *extension;
    char *file_name;
    char *static_contents;
} static_file;

typedef struct {
    char *extension;
    char *type_str;
    int file_type;
} content_type_t;

#define FILE_EXT_LEN 10
#define FILE_TYPE_BINARY 0
#define FILE_TYPE_TEXT 1
#define CONTENT_TYPES_LEN 28

static content_type_t content_types[CONTENT_TYPES_LEN] = {
    { "html", "text/html", FILE_TYPE_TEXT },
    { "htm", "text/html", FILE_TYPE_TEXT },
    { "css", "text/css", FILE_TYPE_TEXT },
    { "md", "text/markdown", FILE_TYPE_TEXT },
    { "txt", "text/plain", FILE_TYPE_TEXT },
    { "xml", "text/xml", FILE_TYPE_TEXT },
    { "js", "application/javascript", FILE_TYPE_TEXT },
    { "json", "application/json", FILE_TYPE_TEXT },
    { "pdf", "application/pdf", FILE_TYPE_BINARY },
    { "zip", "application/zip", FILE_TYPE_BINARY },
    { "bmp", "image/bmp", FILE_TYPE_BINARY },
    { "gif", "image/gif", FILE_TYPE_BINARY },
    { "jpg", "image/jpeg", FILE_TYPE_BINARY },
    { "jpeg", "image/jpeg", FILE_TYPE_BINARY },
    { "ico", "image/x-icon", FILE_TYPE_BINARY },
    { "png", "image/png", FILE_TYPE_BINARY },
    { "tiff", "image/tiff", FILE_TYPE_BINARY },
    { "svg", "image/svg", FILE_TYPE_TEXT },
    { "mp4", "audio/mp4", FILE_TYPE_BINARY },
    { "mp4", "video/mp4", FILE_TYPE_BINARY },
    { "mpeg", "audio/mpeg", FILE_TYPE_BINARY },
    { "mpeg", "video/mpeg", FILE_TYPE_BINARY },
    { "ogg", "audio/ogg", FILE_TYPE_BINARY },
    { "ogg", "video/ogg", FILE_TYPE_BINARY },
    { "quicktime", "video/quicktime", FILE_TYPE_BINARY },
    { "ttf", "font/ttf", FILE_TYPE_BINARY },
    { "woff", "font/woff", FILE_TYPE_BINARY },
    { "woff2", "font/woff2", FILE_TYPE_BINARY },
};

size_t _get_binary_file_contents_len(int fd)
{
    size_t contents_len = (size_t)lseek(fd, (off_t)0, SEEK_END);

    lseek(fd, (off_t)0, SEEK_SET);

    return contents_len;
}

size_t _get_text_file_contents_len(FILE *fp)
{
    size_t contents_len = 0;
    int c;
    while ((c = fgetc(fp)) != EOF) {
        contents_len++;
    }
    return contents_len;
}

char *_set_binary_file_contents(int fd, size_t contents_len)
{
    char *contents = NULL;
    long bytes_r;
    contents = malloc(contents_len);
    if (!contents) {
        return NULL;
    }
    memset(contents, 0, contents_len);
    bytes_r = read(fd, contents, contents_len);

    if (bytes_r < 0) {
        return NULL;
    }
    return contents;
}

char *_set_text_file_contents(FILE *fp, size_t contents_len)
{
    char *contents = NULL;
    char *contents_p = NULL;

    contents = malloc(contents_len + 1);
    if (!contents) {
        return NULL;
    }
    memset(contents, 0, contents_len + 1);
    contents_p = contents;
    rewind(fp);
    int c;
    while ((c = fgetc(fp)) != EOF) {
        *contents_p = c;
        contents_p++;
    }
    *contents_p = '\0';
    return contents;
}

char *_get_content_len_str(size_t contents_len)
{
    int content_len_str_len = 12;
    char *content_len_str = malloc(content_len_str_len);
    if (!content_len_str) {
        return NULL;
    }
    memset(content_len_str, 0, content_len_str_len);
    sprintf(content_len_str, "%ld", contents_len);
    return content_len_str;
}

char *_get_file_path(request *req)
{
    char *path = ""; // TODO: This should come from a configuration setting
    int file_to_serve_len = strlen(req->uri);
    // int file_to_serve_len = strlen(req->uri) + 1;
    // int file_to_serve_len = (strlen(path) + strlen(req->uri)) + 1;
    char *file_to_serve = malloc(file_to_serve_len);
    if (!file_to_serve) {
        return NULL;
    }
    memset(file_to_serve, 0, file_to_serve_len);
    file_to_serve[0] = '\0';
    strcat(file_to_serve, path);
    strcat(file_to_serve, req->uri);

    return file_to_serve;
}

char *_get_file_ext(char *uri)
{
    size_t uri_len = strlen(uri);
    size_t dot_loc = 0;
    char *uri_p = uri;
    char *file_ext = malloc(FILE_EXT_LEN);

    if (!file_ext) {
        return NULL;
    }
    memset(file_ext, 0, FILE_EXT_LEN);
    while (*uri_p != '\0') {
        if (*(uri_p - 1) == '.') {
            break;
        }
        dot_loc++;
        uri_p++;
    }
    // we couldn't find the file extension, use txt by default
    if (dot_loc == 0) {
        free(file_ext);
        return strdup("txt");
    }
    memcpy(file_ext, uri_p, uri_len - dot_loc);
    file_ext[(uri_len - dot_loc) + 1] = '\0';

    return file_ext;
}

content_type_t *_get_content_type_struct(char *file_ext)
{
    content_type_t *found_content_type = NULL;

    for (int i = 0; i < CONTENT_TYPES_LEN; i++) {
        if (strcmp(file_ext, content_types[i].extension) == 0) {
            found_content_type = &content_types[i];
        }
    }

    return found_content_type;
}

void _serve_file(request *req, response *res, char *contents, size_t contents_len)
{
    char *file_ext = _get_file_ext(req->uri);
    content_type_t *cur_content_type = _get_content_type_struct(file_ext);
    char *content_len_str = _get_content_len_str(contents_len);

    response_set_status_code(res, 200);
    response_set_header(res, "Content-Length", strdup(content_len_str));
    response_set_header(res, "Content-Type", strdup(cur_content_type->type_str));
    response_set_body(res, contents);
    response_set_body_len(res, contents_len);

    free(file_ext);
    free(content_len_str);
}

bool static_file_serve(request *req, response *res)
{
    char *file_to_serve = _get_file_path(req);
    if (!file_to_serve) {
        return false;
    }

    char *file_ext = _get_file_ext(req->uri);
    content_type_t *cur_content_type = _get_content_type_struct(file_ext);
    char *contents = NULL;
    size_t contents_len;

    if (cur_content_type != NULL && cur_content_type->file_type == FILE_TYPE_BINARY) {
        int fd = open(file_to_serve, O_RDONLY);
        if (fd == -1) {
            perror("file opening failed");
            return false;
        }
        contents_len = _get_binary_file_contents_len(fd);
        if (contents_len == 0) {
            return false;
        }
        contents = _set_binary_file_contents(fd, contents_len);
        if (!contents) {
            return false;
        }
    } else if (cur_content_type != NULL && cur_content_type->file_type == FILE_TYPE_TEXT) {
        // it's a text file
        FILE *fp = fopen(file_to_serve, "r");
        if (!fp) {
            perror("file opening failed");
            printf("%s\n", file_to_serve);
            return false;
        }
        contents_len = _get_text_file_contents_len(fp);
        if (contents_len == 0) {
            fclose(fp);
            return false;
        }
        contents = _set_text_file_contents(fp, contents_len);
        if (!contents) {
            fclose(fp);
            return false;
        }
        fclose(fp);
    } else {
        return false;
    }

    _serve_file(req, res, contents, contents_len);

    free(file_ext);
    free(file_to_serve);
    return true;
}
