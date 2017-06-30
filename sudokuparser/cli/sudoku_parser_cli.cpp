#include "../detect_digits.hpp"
#include "../identify_digits.hpp"
#include "../sudoku_parser.hpp"

using namespace std;
using namespace Sudoku;

int main(int argc, char *argv[])
{
    // check if there is more than one argument and use the second one
  //  (the first argument is the executable)
  if (argc > 2)
  {
      if (string(argv[1]) == "train") {
          string trainConfigFile(argv[2]);

          string result = internalTrainSudoku(argv[2]);
          cout << "Training returned: '" << result << "'" << endl;

          //cv::waitKey(0);
          return 0;
      } else if (string(argv[1]) == "parse") {
          ifstream is (argv[2], std::ifstream::binary);
          if (is) {
                // get length of file:
                is.seekg (0, is.end);
                int length = is.tellg();
                is.seekg (0, is.beg);

                char * buffer = new char [length];
                cout << "Reading " << length << " characters... ";
                // read data as a block:
                is.read (buffer,length);
                bool fileReadSuccessfully = false;
                if (is) {
                    cout << "all characters read successfully." << endl;
                    fileReadSuccessfully = true;
                }
                else 
                {
                    cout << "error: only " << is.gcount() << " could be read";
                }
                is.close();

                if (fileReadSuccessfully) {
                    float * gridPoints;
                    string parsed = internalParseSudoku(buffer, length, gridPoints, true);
                    cout << parsed << endl;
                }

                delete[] buffer;

                if (string(argv[3]) == "wait") {
                    cv::waitKey(0);
                }
          }
      }

  } else {
      cout << "No argument" << endl;
  }
}