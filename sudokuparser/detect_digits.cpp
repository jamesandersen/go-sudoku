#include <iostream>
#include <cmath>
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

//#define VERBOSE

namespace Sudoku {

    // Expand rects around digit contours by this amount
    const int DIGIT_PADDING = 3;

    // Ignore any contour rect smaller than this on any side
    const int MIN_DIGIT_PIXELS = 20;

    // Scale up puzzle images smaller than this
    const int MIN_PUZZLE_SIZE = 325;

    // Scale down puzzle images larger than this
    const int MAX_PUZZLE_SIZE = 900;

    // Resize digits to this size when exporting to train SVM
    const int EXPORT_DIGIT_SIZE = 28;

    const float MIN_GRID_PCT = 0.3;

    const int CANNY_THRESHOLD = 65;

    Size findCorners(vector<Point>, Point2f[]);

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
    * Attempt to extract and warp sudoku grid
    */
    void extractGrid(const Mat& img, Mat& dst, vector<float>& gridPoints, float& scale) {
        Mat src_gray;
        blur( img, src_gray, Size(3,3) );
        Mat canny_output;
        vector<vector<Point> > contours;
        vector<Vec4i> hierarchy;

        int largest_area=0;
        int largest_contour_index=0;
        vector<Point> largestContour;
        Rect bounding_rect;

        Mat denoised;
        fastNlMeansDenoising(src_gray, denoised, 10);

        adaptiveThreshold(~denoised, src_gray, 255, ADAPTIVE_THRESH_GAUSSIAN_C, THRESH_BINARY, 11, -2);
        #ifdef VERBOSE
        //imshow("denoised", src_gray);
        imwrite("artifact_01_denoised.png", src_gray);
        #endif
        Canny( src_gray, canny_output, CANNY_THRESHOLD, CANNY_THRESHOLD * 2, 3 );
        #ifdef VERBOSE
        imwrite("artifact_02_canny.png", canny_output);
        #endif
        findContours( canny_output, contours, hierarchy, RETR_EXTERNAL, CHAIN_APPROX_SIMPLE, Point(0, 0) );
        Mat drawing = Mat::zeros( canny_output.size(), CV_8UC3 );
        
        for( size_t i = 0; i< contours.size(); i++ )
        {   
            double a=contourArea( contours[i],false);  //  Find the area of contour
            if(a>largest_area){
                largest_area=a;
                largest_contour_index=i;                //Store the index of largest contour
                largestContour = contours[i];
                bounding_rect=boundingRect(contours[i]); // Find the bounding rectangle for biggest contour
            }

            #ifdef VERBOSE
            Scalar color = Scalar( rng.uniform(0, 255), rng.uniform(0,255), rng.uniform(0,255) );
            drawContours( drawing, contours, i, color, 2, 8, hierarchy, 0, Point() );
            #endif
        }

        #ifdef VERBOSE
        cout << contours.size() << " canny contours found." << endl;
        imwrite("artifact_03_contours.png", drawing);
        cout << "largest contour index: " << largest_contour_index << endl;
        cout << "largest contour area: " << largest_area << endl;
        #endif

        float area = img.cols * img.rows;
        if (largest_area < area * MIN_GRID_PCT) {
            cout << "largest contour area is only " << (largest_area / area) * 100 << "% of source; aborting grid extraction" << endl;
            img.copyTo(dst);
            return;
        }

        // findCorners
        if (largest_contour_index < contours.size()) {
            Point2f corners[4];
            Point2f flatCorners[4];
            Size sz = findCorners(largestContour, corners);

            
            // draw contour quadrangle
            for( int j = 0; j < 4; j++ ) {
                gridPoints.push_back(corners[j].x * scale);
                gridPoints.push_back(corners[j].y * scale);
                #ifdef VERBOSE
                line( drawing, corners[j], corners[(j+1)%4], Scalar(0,0,255), 2, 8 );
                #endif
            }
            #ifdef VERBOSE
            // draw largest 
            int lenContour = largestContour.size();
            for( int j = 0; j < lenContour; j++ ) {
                line( drawing, largestContour[j], largestContour[(j+1)%lenContour], Scalar(0,0,255), 2, 8 );
            }
            imwrite("artifact_04_quadrangle.png", drawing);
            #endif
            flatCorners[0] = Point2f(0, 0);
            flatCorners[1] = Point2f(sz.width, 0);
            flatCorners[2] = Point2f(sz.width, sz.height);
            flatCorners[3] = Point2f(0, sz.height);
            Mat lambda = getPerspectiveTransform(corners, flatCorners);

            Mat output;
            warpPerspective(denoised, output, lambda, sz);

            GaussianBlur(output, dst, Size(0, 0), 3);
            addWeighted(output, 1.5, dst, -0.5, 0, dst);
            #ifdef VERBOSE
            //imshow("Warped", dst);
            imwrite("artifact_05_warped.png", dst);
            #endif
        } else {
            throw runtime_error("No grid contour found");
        }
    }

    /**
    * Find corners within the largest contour for use in performing warp transform
    */
    Size findCorners(const vector<Point> largestContour, Point2f corners[]) {
        float dist;
        float maxDist[4] = {0, 0, 0, 0};

        Moments M = moments(largestContour, true);
        float cX = float(M.m10 / M.m00);
	    float cY = float(M.m01 / M.m00);

        for (int i = 0; i < 4; i++) {
            maxDist[i] = 0.0f;
            corners[i] = Point2f(cX, cY);
        }
        
        Point2f center(cX, cY);
        
        // find the most distant point in the contour within each quadrant
        for (int i = 0; i < largestContour.size(); i++) {
            Point2f p(largestContour[i].x, largestContour[i].y);
            dist = sqrt(pow(p.x - center.x, 2) + pow(p.y - center.y, 2));
            // cout << "(" << p.x << "," << p.y << ") is " << dist << " from (" << cX << "," << cY << ")" << endl;
            if (p.x < cX && p.y < cY && maxDist[0] < dist) {
                // top left
                maxDist[0] = dist;
                corners[0] = p;
            } else if (p.x > cX && p.y < cY && maxDist[1] < dist) {
                // top right
                maxDist[1] = dist;
                corners[1] = p;
            } else if (p.x > cX && p.y > cY && maxDist[2] < dist) {
                // bottom right
                maxDist[2] = dist;
                corners[2] = p;
            } else if (p.x < cX && p.y > cY && maxDist[3] < dist) {
                // bottom left
                maxDist[3] = dist;
                corners[3] = p;
            }
        }

        // determine the dimensions that we should warp this contour to
        float widthTop = sqrt(pow(corners[0].x - corners[1].x, 2) + pow(corners[0].y - corners[1].y, 2));
        float widthBottom = sqrt(pow(corners[2].x - corners[3].x, 2) + pow(corners[2].y - corners[3].y, 2));

        float heightLeft = sqrt(pow(corners[0].x - corners[3].x, 2) + pow(corners[0].y - corners[3].y, 2));
        float heightRight = sqrt(pow(corners[1].x - corners[2].x, 2) + pow(corners[1].y - corners[2].y, 2));

#ifdef VERBOSE
        cout << "TopLeft: (" << corners[0].x << "," << corners[0].y << ")" << endl;
        cout << "TopRight: (" << corners[1].x << "," << corners[1].y << ")" << endl;
        cout << "BottomRight: (" << corners[2].x << "," << corners[2].y << ")" << endl;
        cout << "BottomLeft: (" << corners[3].x << "," << corners[3].y << ")" << endl;

        cout << "widthTop: " << widthTop << endl << "widthBottom: " << widthBottom << endl;
        cout << "heightLeft: " << heightLeft << endl << "heightRight: " << heightRight << endl;
#endif

        return Size(max(widthTop, widthBottom), max(heightLeft, heightRight));

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
    vector<Rect> FindDigitRects(const Mat& raw, Mat& cleaned, vector<float>& gridPoints, float &scale) {
        // Check if image is loaded fine
        if(!raw.data)
            cerr << "Problem loading image!!!" << endl;

        Mat src;
        raw.copyTo(src);

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
        
        // make sure image is a reasonable size
        if(gray.rows > MAX_PUZZLE_SIZE || gray.cols > MAX_PUZZLE_SIZE) {
            scale = max(gray.rows, gray.cols) / float(MAX_PUZZLE_SIZE);
            resize(gray, gray, Size(gray.cols / scale, gray.rows / scale), 0, 0, CV_INTER_AREA);
        } else if (gray.rows < MIN_PUZZLE_SIZE || gray.cols < MIN_PUZZLE_SIZE) {
            scale = min(gray.rows, gray.cols) / float(MAX_PUZZLE_SIZE);
            resize(gray, gray, Size(gray.cols / scale, gray.rows / scale), 0, 0, CV_INTER_CUBIC);
        }

        Mat grid = Mat::zeros( gray.size(), gray.type() );
        extractGrid(gray, grid, gridPoints, scale);
        
        // Apply adaptiveThreshold at the bitwise_not of gray, notice the ~ symbol
        Mat bw;
        adaptiveThreshold(~grid, bw, 255, ADAPTIVE_THRESH_GAUSSIAN_C, THRESH_BINARY, 11, -2);

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

        #ifdef VERBOSE
        imwrite("artifact_06_cleaned.png", clean);
        //imshow("cleaned", clean);
        #endif
        return digits;
    }
}

