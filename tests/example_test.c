#include <string.h>
#include <stdbool.h>
#include "unittest.h"

char *test_example()
{
    assert(true, "Something went wrong");
    return NULL;
}

int main(void)
{
    start_tests("example tests");
    run_test(test_example);
    end_tests();

    return 0;
}
