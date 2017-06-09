#ifndef  DETECT_DIGITS_INC 
#define  DETECT_DIGITS_INC

#include <string>
#include <opencv2/opencv.hpp>

namespace Sudoku {
    // Resize digits to this size when exporting to train SVM
    const int EXPORT_DIGIT_SIZE = 28;

    void extractDigits(char* file);
    std::vector<cv::Rect> FindDigitRects(const cv::Mat& img, cv::Mat& cleaned, std::vector<float>& gridPoints, float &scale);
}

#endif