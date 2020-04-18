VPATH := src tests bin
CFLAGS := gcc -std=gnu11 -Wall -Wextra -Isrc

all: fserve

fserve: fserve.c fserve.h request.o response.o static_file.o h_table.o array.o dl_list.o
	$(CFLAGS) src/fserve.c request.o response.o static_file.o h_table.o array.o dl_list.o -o bin/fserve

request.o: request.c request.h
	$(CFLAGS) -c src/request.c src/request.h
response.o: response.c response.h
	$(CFLAGS) -c src/response.c src/response.h
router.o: router.c router.h
	$(CFLAGS) -c src/router.c src/router.h
static_file.o: static_file.c static_file.h
	$(CFLAGS) -c src/static_file.c src/static_file.h
h_table.o: h_table.c h_table.h
	$(CFLAGS) -c src/h_table.c src/h_table.h
array.o: array.c array.h
	$(CFLAGS) -c src/array.c src/array.h
dl_list.o: dl_list.c dl_list.h
	$(CFLAGS) -c src/dl_list.c src/dl_list.h

# Tests
request_test: request_test.c request.o
	$(CFLAGS) tests/request_test.c request.o h_table.o array.o -o bin/request_test
example_test: example_test.c
	$(CFLAGS) tests/example_test.c -o bin/example_test

.PHONY: test
test: example_test
	./bin/example_test

clean:
	rm -f ./*.o src/*.h.gch
	rm -Rf ./bin && mkdir bin && touch ./bin/.gitkeep
