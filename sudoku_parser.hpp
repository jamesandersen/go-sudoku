#ifndef _SUDOKU_PARSER_HPP_
#define _SUDOKU_PARSER_HPP_

#include <string>

using namespace std;

int internalTest();
void internalTest2(char * filename);
string internalParseSudoku(char * encodedImageData, bool saveOutput);

#endif