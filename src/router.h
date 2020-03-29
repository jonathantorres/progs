#ifndef _ROUTER_H
#define _ROUTER_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "dl_list.h"
#include "request.h"
#include "response.h"

bool router_handle_request(request *req, response *res);

#endif
