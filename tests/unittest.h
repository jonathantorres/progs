#include <stdlib.h>
#include <stdio.h>

typedef char *(*test_func)();

int num_of_tests = 0;
int num_of_assertions = 0;

inline static void run_test(test_func func)
{
    num_of_tests++;

    // execute the function, if it returns a string,
    // this means that the test failed.
    // Exit the program with the failure message
    char *test_result = func();

    if (test_result != NULL) {
        fprintf(stderr, "[TEST FAIL]: %s\n", test_result);
        exit(EXIT_FAILURE);
    }
}

inline static void assert(int expr, char *msg)
{
    num_of_assertions++;

    if (!expr) {
        fprintf(stderr, "[ASSERTION FAIL]: %s\n", msg);
        exit(EXIT_FAILURE);
    }
}

inline static void start_tests(char *test_name)
{
    printf("Running %s...", test_name);
}

inline static void end_tests()
{
    printf("Passed. %d tests and %d assertions.\n", num_of_tests, num_of_assertions);
}
