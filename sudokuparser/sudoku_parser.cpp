#include "detect_digits.hpp"
#include "identify_digits.hpp"
#include "sudoku_parser.hpp"

#include <string>
#include <tuple>
#include <iostream>
#include <fstream>
#include <iterator>
#include <opencv2/opencv.hpp>

using namespace Sudoku;
using namespace std;
using namespace cv;

const char *SVM_MODEL_ENV_VAR_NAME = "GO_SUDOKU_SVM_MODEL";

const string internalParseSudoku(const char * encImgData, int length, float * gridPoints, bool saveOutput) {
    std::vector<char> encodedImageData(encImgData, encImgData + length);

    Mat sudokuBoard = imdecode(encodedImageData, CV_LOAD_IMAGE_ANYDEPTH);
    cout << "channels: " << sudokuBoard.channels() << " type: " << sudokuBoard.type() << endl;
    if (sudokuBoard.type() == 2) {
        sudokuBoard.convertTo(sudokuBoard, CV_8U, 0.00390625);
    }
    
    Mat cleanedBoard;
    vector<float> gPoints;
    float scale = 1.0;
    vector<Rect> digits = FindDigitRects(sudokuBoard, cleanedBoard, gPoints, scale);


    map<string, int> digitMap;

    if (digits.size() > 0) {
        // get the bounding box of all digits
        Rect allDigits = digits[0];
        for( size_t i = 0; i< digits.size(); i++ ) { allDigits |= digits[i]; }
        if (gPoints.size() == 0) {
            gridPoints[0] = allDigits.x * scale;
            gridPoints[1] = allDigits.y * scale;
            gridPoints[2] = (allDigits.x + allDigits.width) * scale;
            gridPoints[3] = allDigits.y * scale;
            gridPoints[4] = allDigits.br().x * scale;
            gridPoints[5] = allDigits.br().y * scale;
            gridPoints[6] = allDigits.x * scale;
            gridPoints[7] = (allDigits.y + allDigits.height) * scale;
        }

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
        }
        
        rectangle( digitBounds, allDigits, teal, 1, 8, 0 );
        #ifdef VERBOSE
        imwrite("artifact_07_detected.png", digitBounds);
        #endif
    }

    string puzzle = "";
    for (int y = 0; y < 9; y++ ) {
        for (int x = 0; x < 9; x++) {
            string key = string(1, "ABCDEFGHI"[y]) + to_string(x + 1);
            auto search = digitMap.find(key);
            if(search != digitMap.end()) {
                puzzle += to_string(search->second);
            } else {
                puzzle += ".";
            }
            
        }
    }
    
    cout << length << " byte puzzle parsed as " << puzzle << endl;

    // set the grid corners
    if (gPoints.size() == 8) {
        copy(gPoints.begin(), gPoints.end(), gridPoints);
    }

    return puzzle;
}

// https://stackoverflow.com/a/9676623/385152
template<typename T>
vector<T> 
split(const T & str, const T & delimiters) {
    vector<T> v;
    typename T::size_type start = 0;
    auto pos = str.find_first_of(delimiters, start);
    while(pos != T::npos) {
        if(pos != start) // ignore empty tokens
            v.emplace_back(str, start, pos - start);
        start = pos + 1;
        pos = str.find_first_of(delimiters, start);
    }
    if(start < str.length()) // ignore trailing delimiter
        v.emplace_back(str, start, str.length() - start); // add what's left of the string
    return v;
}

map<string, string> parseTrainConfig(const char * trainConfigFile) {
    map<string, string> trainFiles;
    ifstream file(trainConfigFile);
    string line;
    while(getline(file, line)) {
        vector<string> v = split<string>(line, ",");
        if (v.size() == 2) {
            trainFiles[v[0]] = v[1];
        } else {
            cout << "Discarding line with unexpected number of tokens: " << line << endl;
        }
    }
    return trainFiles;
}

map<int, vector<Mat> > labelDigits(Mat &clean, vector<Rect> digits, string labels) {
    map<int, vector<Mat> > labeled;
    
    if (labels.length() != 81) {
        throw runtime_error("Invalid label string!");
    }

    if (digits.size() > 0) {
        Rect allDigits = digits[0];
        for( size_t i = 0; i < digits.size(); i++ ) { allDigits |= digits[i]; }

        double cellWidth = allDigits.width / 9.0;
        double cellHeight = allDigits.height / 9.0;

        for( size_t i = 0; i< digits.size(); i++ )
        {
            Point center = (digits[i].br() + digits[i].tl())*0.5;
            int row = int(floor((center.y - allDigits.y) / cellHeight));
            int col = int(floor((center.x - allDigits.x) / cellWidth));

            string labelStr = labels.substr((row * 9) + col, 1);
            if (labelStr != ".") {
                int label = atoi(labelStr.c_str());

                // Extract the digit
                Mat digit = Mat(clean, digits[i]);
                resize(digit, digit, Size(EXPORT_DIGIT_SIZE, EXPORT_DIGIT_SIZE), 0, 0, CV_INTER_AREA);
                // despeckle
                fastNlMeansDenoising(digit, digit, 50.0, 5, clean.cols / 10);

                labeled[label].push_back(digit);
            }

            // draw random colored rect
            //Scalar color = Scalar( rng.uniform(0, 255), rng.uniform(0,255), rng.uniform(0,255) );
            //rectangle( digitBounds, digits[i], color, 1, 8, 0 );


            //replace(filename, "samples/", "");
            //string outFile = "digits/" + string(1, rowChar) + to_string(col) + "_" + filename;
            //imwrite( outFile, digit );
        }
    }
    

    return labeled;
}

string internalTrainSudoku(const char * trainConfigFile) {
    map<string, string> trainFiles = parseTrainConfig(trainConfigFile);

    map<int, vector<Mat>> allLabeledDigits;
    for( const pair<string, string> element : trainFiles )
    {
        try {
            // read sample image and find digits
            auto sudokuBoard = imread(element.first, CV_LOAD_IMAGE_ANYDEPTH);
            Mat cleanedBoard;
            vector<float> gridPoints; // not used for training
            float scale = 1.0;
            auto digits = FindDigitRects(sudokuBoard, cleanedBoard, gridPoints, scale);

            // extract digit images with labels
            auto labeledDigits = labelDigits(cleanedBoard, digits, element.second);
            for( const pair<int, vector<Mat> > el : labeledDigits )
            {
                for(Mat digit : el.second)
                    allLabeledDigits[el.first].push_back(digit);
            }
            //imshow(file, cleanedBoard);
            cout << "captured " << labeledDigits.size() << " digits from " << element.first << endl;
        } catch (const std::exception& e) {
            cout << "Exception occurred while processing " << element.first << ": " << e.what() << endl;
        } catch (...) {
            cout << "Exception occurred while processing " << element.first << endl;
        }
    }

    // determine the size of the training set
    int maxCols = 0;
    for( const pair<int, vector<Mat> > el : allLabeledDigits )
    {
        maxCols = max(maxCols, static_cast<int>(el.second.size()) * EXPORT_DIGIT_SIZE);
    }
    Mat combined = Mat::zeros(Size(maxCols, static_cast<int>(allLabeledDigits.size()) * EXPORT_DIGIT_SIZE), CV_8UC1);
    for( const pair<int, vector<Mat> > el: allLabeledDigits )
    {
        int x = 0;
        int y = (el.first - 1) * EXPORT_DIGIT_SIZE;
        for (Mat digit : el.second) {
            digit.copyTo(combined(Rect(x, y, EXPORT_DIGIT_SIZE, EXPORT_DIGIT_SIZE)));
            x += EXPORT_DIGIT_SIZE;
        }
    }

    // write out the training data image
    imwrite("combined.png", combined);

    // Ensure model file location is set
    if (!getenv(SVM_MODEL_ENV_VAR_NAME)) {
        setenv(SVM_MODEL_ENV_VAR_NAME, "model4.yml", 0);
    }

    // Train the SVM
    return TrainSVM("combined.png", EXPORT_DIGIT_SIZE);
}

