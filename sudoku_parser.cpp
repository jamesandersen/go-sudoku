#include "detect_digits.hpp"
#include "identify_digits.hpp"

#include <string>
#include <iostream>
#include <fstream>
#include <iterator>
#include <opencv2/opencv.hpp>

using namespace Sudoku;
using namespace std;
using namespace cv;

const int EXPORT_DIGIT_SIZE = 28;

string internalParseSudoku(char * encImgData, bool saveOutput) {

    int zi = sizeof(encImgData);
    std::vector<char> encodedImageData(encImgData, encImgData + zi);

    Mat sudokuBoard = imdecode(encodedImageData, CV_LOAD_IMAGE_ANYDEPTH);
    Mat cleanedBoard;
    vector<Rect> digits = FindDigitRects(sudokuBoard, cleanedBoard);
    map<string, int> digitMap;

    for (int y = 0; y < 9; y++ ) {
        for (int x = 0; x < 9; x++) {

        }
    }

    if (digits.size() > 0) {
        // get the bounding box of all digits
        Rect allDigits = digits[0];
        for( size_t i = 0; i< digits.size(); i++ ) { allDigits |= digits[i]; }

        double cellWidth = allDigits.width / 9.0;
        double cellHeight = allDigits.height / 9.0;

        Mat digitBounds = cleanedBoard.clone();
        Scalar pink = Scalar(255, 105, 180);
        Scalar teal = Scalar(20, 135, 128);
        if (saveOutput) {
            cvtColor( digitBounds, digitBounds, COLOR_GRAY2BGR );
        }

        for( size_t i = 0; i< digits.size(); i++ )
        {
            Point center = (digits[i].br() + digits[i].tl())*0.5;
            int row = int(floor((center.y - allDigits.y) / cellHeight));
            char rowChar = "ABCDEFGHI"[row];
            int col = int(floor((center.x - allDigits.x) / cellWidth));

            // save the digit
            Mat digitImg = Mat(cleanedBoard, digits[i]);
            resize(digitImg, digitImg, Size(EXPORT_DIGIT_SIZE, EXPORT_DIGIT_SIZE), 0, 0, CV_INTER_AREA);
            // despeckle
            fastNlMeansDenoising(digitImg, digitImg, 50.0, 5, cleanedBoard.cols / 10);

            int digit = IdentifyDigit(digitImg);

            if (saveOutput) {
                rectangle( digitBounds, digits[i], pink, 1, 8, 0 );
                putText(digitBounds, to_string(digit), center + Point(5, 12), FONT_HERSHEY_PLAIN, 0.8, teal);
            }
            
            digitMap[string(1, rowChar) + to_string(col + 1)] = digit;
            //cout << string(1, rowChar) + to_string(col + 1) + " identified as " + to_string(digit) << endl;
        }
        
        rectangle( digitBounds, allDigits, teal, 1, 8, 0 );
        imwrite("detected.png", digitBounds);
    }

    string puzzle = "";
    for (int y = 0; y < 9; y++ ) {
        for (int x = 0; x < 9; x++) {
            string key = string(1, "ABCDEFGHI"[y]) + to_string(x + 1);
            auto search = digitMap.find(key);
            if(search != digitMap.end()) {
                puzzle += to_string(search->second);
            }
            puzzle += ",";
        }
    }
    
    return puzzle.substr(0, puzzle.length() - 1);
}

int internalTest() { return 1; }

void internalTest2(char * filename) {
    cout << "Called by GOLANG!" << endl;
}

/*
int main(int argc, char *argv[])
{
    // check if there is more than one argument and use the second one
  //  (the first argument is the executable)
  if (argc > 3)
  {
    string file(argv[1]);
    string digitFile(argv[2]);
    string sampleBoardFile(argv[3]);

    
    ifstream testFile(argv[3], std::ios::binary);
    vector<char> fileContents((istreambuf_iterator<char>(testFile)), istreambuf_iterator<char>());

    //char* boardBytes = readFileBytes(argv[3]);

    string boardString = ParseSudoku(fileContents, true);
    cout << sampleBoardFile << " parsed as " << boardString << endl;
    //Mat digitImg = imread(digitFile);
    //int digit = IdentifyDigit(digitImg);

    ///TrainSVM(file, 28);

  } else {
      cout << "No argument" << endl;
  }

}
*/