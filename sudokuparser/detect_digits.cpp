#include <iostream>
#include <map>
#include <string>
#include <iterator>
#include <functional>
#include <opencv2/opencv.hpp>

using namespace std;
using namespace cv;
RNG rng(12345);

using LineTestFn = function<bool(Rect&, Mat&)>;
using ExpandRectFn = function<Rect(Rect&, Mat&)>;

namespace Sudoku {

    // Expand rects around digit contours by this amount
    const int DIGIT_PADDING = 3;

    // Ignore any contour rect smaller than this on any side
    const int MIN_DIGIT_PIXELS = 15;

    // Scale up puzzle images smaller than this
    const int MIN_PUZZLE_SIZE = 250;

    // Scale down puzzle images larger than this
    const int MAX_PUZZLE_SIZE = 800;

    // Resize digits to this size when exporting to train SVM
    const int EXPORT_DIGIT_SIZE = 28;

    /**
    * Detect numeric digits in a sudoku grid in img Mat and return Rect instances where they are found
    */
    
    vector<Rect> findDigits(const Mat& img) {
        vector<Rect> digits;
        Rect imgRect = Rect(0, 0, img.cols, img.rows);

        // Find all contours
        vector<vector<Point> > contours;
        vector<Vec4i> hierarchy;
        //findContours( img, contours, hierarchy, RETR_TREE, CHAIN_APPROX_SIMPLE, Point(0, 0) );
        findContours( img, contours, hierarchy, RETR_EXTERNAL, CHAIN_APPROX_SIMPLE, Point(0, 0) );
        vector<vector<Point> > contours_poly( contours.size() );
        vector<Rect> boundRect( contours.size() );

        for( size_t i = 0; i < contours.size(); i++ )
        {
            approxPolyDP( Mat(contours[i]), contours_poly[i], 3, true );
            boundRect[i] = boundingRect( Mat(contours_poly[i]) );
            
            if( hierarchy[i][3] < 0 && boundRect[i].height >= MIN_DIGIT_PIXELS) { // find "root" contours
                double aspectRatio = boundRect[i].height / double(boundRect[i].width);
                // check reasonable aspect ratio for digits
                if ((aspectRatio >= 1 && aspectRatio < 3.2)) { 
                    int widthToAdd = boundRect[i].height - boundRect[i].width + (2 * DIGIT_PADDING);
                    int pointOffset = int(floor(double(widthToAdd / 2)));
                    boundRect[i] = boundRect[i] - Point(pointOffset, DIGIT_PADDING);
                    boundRect[i] = boundRect[i] + Size(widthToAdd, 2 * DIGIT_PADDING);
                    boundRect[i] &= imgRect;

                    // check white/black pixel ratio to avoid accidental noise getting picked up
                    double wbRatio = countNonZero(img(boundRect[i])) / double(boundRect[i].width * boundRect[i].height);
                    if (wbRatio > 0.1 && wbRatio < 0.4) {
                        digits.push_back(boundRect[i]); // AND boundRect and imgRect to ensure in boundary
                    }
                }
            }
        }

        return digits;
    }

    /**
    * Detect horizontal/vertical lines of a sudoku grid in img Mat and copy expanded lines to dst Mat
    */
    void extractLines(const Mat& img, Mat& dst, bool horizontal) {
        // Clone the source image
        Mat clone = img.clone();

        LineTestFn lineTest;
        ExpandRectFn expandRect;
        Size size;

        // Setup correct structure size, line test and rect expansion for horizontal vs. vertical
        if (horizontal) {
            size = Size(img.cols / 9, 1);
            lineTest = [](Rect& rect, Mat& mat) { return rect.height / double(mat.rows) < 0.05 && rect.width / double(mat.cols) > 0.111; };
            expandRect = [](Rect& rect, Mat& mat) { 
                Rect expanded = rect;
                if (expanded.y > 1) { expanded.y -= 2; }
                
                if (expanded.y + expanded.height < mat.rows) {
                    expanded.height += min(4, mat.rows - expanded.y - expanded.height);
                }
                expanded.x = 0;
                expanded.width = mat.cols;
                return expanded; 
            };
        } else {
            size = Size(1, img.rows / 9);
            lineTest = [](Rect& rect, Mat& mat) { return rect.width / double(mat.cols) < 0.05 && rect.height / double(mat.rows) > 0.111; };
            expandRect = [](Rect& rect, Mat& mat) { 
                Rect expanded = rect;
                if (expanded.x > 1) { expanded.x -= 2; }
                
                if (expanded.x + expanded.width < mat.cols) {
                    expanded.width += min(4, mat.cols - expanded.x - expanded.width);
                }
                expanded.y = 0;
                expanded.height = mat.rows;
                return expanded; 
            };
        }
        // Create structure element for extracting lines through morphology operations
        Mat structure = getStructuringElement(MORPH_RECT, size);

        // Apply morphology operations
        erode(clone, clone, structure, Point(-1, -1));
        dilate(clone, clone, structure, Point(-1, -1));

        // Find all contours
        vector<vector<Point> > contours;
        vector<Vec4i> hierarchy;
        findContours( clone, contours, hierarchy, RETR_TREE, CHAIN_APPROX_SIMPLE, Point(0, 0) );
        vector<vector<Point> > contours_poly( contours.size() );
        vector<Rect> boundRect( contours.size() );

        // Mark contours which pass line test in the destination image
        for( size_t i = 0; i < contours.size(); i++ )
        {
            approxPolyDP( Mat(contours[i]), contours_poly[i], 3, true );
            boundRect[i] = boundingRect( Mat(contours_poly[i]) );
            if (lineTest(boundRect[i], clone)) {
                Rect expanded = expandRect(boundRect[i], clone);
                dst(expanded) |= 255; // set the expanded rect to white
            }
        }
    }

    // Really? ... this needs to be done manually?
    // http://stackoverflow.com/questions/3418231/replace-part-of-a-string-with-another-string
    bool replace(std::string& str, const std::string& from, const std::string& to) {
        size_t start_pos = str.find(from);
        if(start_pos == std::string::npos)
            return false;
        str.replace(start_pos, from.length(), to);
        return true;
    }

    /**
    * Detect Sudoku board and digits in the "raw" Mat
    */
    vector<Rect> FindDigitRects(const Mat& raw, Mat& cleaned) {
        // Check if image is loaded fine
        if(!raw.data)
            cerr << "Problem loading image!!!" << endl;

        Mat src;
        raw.copyTo(src);

        // make sure image is a reasonable size
        if(src.rows > MAX_PUZZLE_SIZE || src.cols > MAX_PUZZLE_SIZE) {
            resize(src, src, Size(src.cols / 2, src.rows / 2), 0, 0, CV_INTER_AREA);
        } else if (src.rows < MIN_PUZZLE_SIZE || src.cols < MIN_PUZZLE_SIZE) {
            resize(src, src, Size(src.cols * 2, src.rows * 2), 0, 0, CV_INTER_CUBIC);
        }

        //imshow("src: " + filename, src);
        // Transform source image to gray if it is not
        Mat gray;
        if (src.channels() == 3)
        {
            cvtColor(src, gray, CV_BGR2GRAY);
        }
        else
        {
            gray = src;
        }
        
        // Apply adaptiveThreshold at the bitwise_not of gray, notice the ~ symbol
        Mat bw;
        adaptiveThreshold(~gray, bw, 255, ADAPTIVE_THRESH_GAUSSIAN_C, THRESH_BINARY, 11, -2);

        Mat lines = Mat::zeros( bw.size(), bw.type() );
        extractLines(bw, lines, true);
        extractLines(bw, lines, false);
        
        //imshow("lines: " + filename, lines);

        // subtract grid lines from the black/white image
        // so they don't interfere with digit detection
        Mat clean = bw - lines;
        blur(clean, clean, Size(1, 1));

        // find digits
        vector<Rect> digits = findDigits(clean);

        clean.copyTo(cleaned);

        return digits;
    }
}

