#include "sudoku_parser.hpp"
#include "sudoku_parser.h"

#include <string>
#include <iostream>

// Provide implemenation of the external C API
// this file may wrap the C++ calls but cannot include C++ types such as vector 

const char *SVM_MODEL_VAR = SVM_MODEL_ENV_VAR_NAME;

const char * ParseSudoku(const char * encodedImageData, int length, bool saveOutput) {
    
    return internalParseSudoku(encodedImageData, length, saveOutput).c_str();
}

const char* TrainSudoku(const char * trainConfigFile) {
    return internalTrainSudoku(trainConfigFile).c_str();
}