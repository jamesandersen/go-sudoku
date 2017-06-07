#include "sudoku_parser.hpp"
#include "sudoku_parser.h"

#include <string>
#include <iostream>

// Provide implemenation of the external C API
// this file may wrap the C++ calls but cannot include C++ types such as vector 

const char *SVM_MODEL_VAR = SVM_MODEL_ENV_VAR_NAME;

const char* ParseSudoku(const char * encodedImageData, int length, bool saveOutput, char * parsed) {
    string result = internalParseSudoku(encodedImageData, length, saveOutput);
    //cout << "Testing string issue: result: " << result << endl << "c_str(): " << result.c_str() << endl;
    cout << "Returning  parsed result: " << result << endl;
    strncpy(parsed, result.c_str(), 81);
    
    return parsed;
}

const char* TrainSudoku(const char * trainConfigFile) {
    return internalTrainSudoku(trainConfigFile).c_str();
}