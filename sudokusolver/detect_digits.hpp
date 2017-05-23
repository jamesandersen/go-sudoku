#ifndef  DETECT_DIGITS_INC 
#define  DETECT_DIGITS_INC

#include <string>
#include <opencv2/opencv.hpp>

namespace Sudoku {
    void extractDigits(char* file);
    std::vector<cv::Rect> FindDigitRects(const cv::Mat& img, cv::Mat& cleaned);
}

#endif