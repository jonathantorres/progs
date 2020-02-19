VPATH := src tests bin
CFLAGS := gcc -std=gnu11 -Wall -Wextra -Isrc

all: test

server: server.c server.h request.o response.o router.o static_file.o htable.o array.o dllist.o
	$(CFLAGS) src/server.c request.o response.o router.o static_file.o htable.o array.o dllist.o -o bin/server

request_test: request_test.c request.o
	$(CFLAGS) tests/request_test.c request.o htable.o array.o -o bin/request_test
request.o: request.c request.h
	$(CFLAGS) -c src/request.c src/request.h

response.o: response.c response.h
	$(CFLAGS) -c src/response.c src/response.h

router.o: router.c router.h
	$(CFLAGS) -c src/router.c src/router.h

static_file.o: static_file.c static_file.h
	$(CFLAGS) -c src/static_file.c src/static_file.h

# Hash Table
htable_test: htable_test.c htable.o array.o
	$(CFLAGS) tests/htable_test.c htable.o array.o -o bin/htable_test
htable.o: htable.c htable.h
	$(CFLAGS) -c src/htable.c src/htable.h

# Array
array.o: array.c array.h
	$(CFLAGS) -c src/array.c src/array.h
array_test: array_test.c array.o
	$(CFLAGS) tests/array_test.c array.o -o bin/array_test

dllist.o: dllist.c dllist.h
	$(CFLAGS) -c src/dllist.c src/dllist.h

# Example test
example_test: example_test.c
	$(CFLAGS) tests/example_test.c -o bin/example_test

.PHONY: test
test: example_test
	./bin/example_test

clean:
	rm -f ./*.o src/*.h.gch
	rm -Rf ./bin && mkdir bin && touch ./bin/.gitkeep
