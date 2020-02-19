#ifndef _SERVER_H
#define _SERVER_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdbool.h>
#include <errno.h>
#include <fcntl.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/select.h>
#include <netinet/in.h>

#include "request.h"
#include "response.h"
#include "router.h"
#include "static_file.h"

#define DEFAULT_PORT 9090
#define RECV_BUFF_LEN 1000000

#endif
