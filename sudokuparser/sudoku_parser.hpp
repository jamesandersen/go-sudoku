#ifndef _SUDOKU_PARSER_HPP_
#define _SUDOKU_PARSER_HPP_

#include <string>

using namespace std;

extern const char *SVM_MODEL_ENV_VAR_NAME;

const string internalParseSudoku(const char * encodedImageData, int length, float * gridPoints, bool saveOutput);

string internalTrainSudoku(const char * trainConfigFile);

#endif