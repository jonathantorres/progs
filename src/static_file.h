#ifndef _STATIC_FILE_H
#define _STATIC_FILE_H

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include "request.h"
#include "response.h"

bool static_file_serve(request *req, response *res);

#endif
