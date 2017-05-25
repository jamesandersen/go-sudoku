#ifndef _SUDOKU_PARSER_HPP_
#define _SUDOKU_PARSER_HPP_

#include <string>

using namespace std;

extern const char *SVM_MODEL_ENV_VAR_NAME;

string internalParseSudoku(const char * encodedImageData, int length, bool saveOutput);

#endif