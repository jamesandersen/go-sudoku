#ifndef _SUDOKU_PARSER_H_
#define _SUDOKU_PARSER_H_

#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif
	int test();
    void test2(char * filename);
    const char* ParseSudoku(const char * encodedImageData, int length, bool saveOutput);

#ifdef __cplusplus
}
#endif

#endif