#ifndef  IDENTIFY_DIGITS_INC 
#define  IDENTIFY_DIGITS_INC 

#include <string>
#include <opencv2/opencv.hpp>

namespace Sudoku {
    void TrainSVM(std::string &pathName, int digitSize);
    int IdentifyDigit(cv::Mat &digitMat);
}

#endif