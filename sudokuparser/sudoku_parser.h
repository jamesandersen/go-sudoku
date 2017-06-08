#ifndef _SUDOKU_PARSER_H_
#define _SUDOKU_PARSER_H_

#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif
    // Define methods to be exposed to Go here

    extern const char *SVM_MODEL_VAR;

    void ParseSudoku(const char * encodedImageData, int length, float * gridPoints, bool saveOutput, char * parsed);

    const char* TrainSudoku(const char * trainConfigFile);

#ifdef __cplusplus
}
#endif

#endif