#include "../detect_digits.hpp"
#include "../identify_digits.hpp"
#include "../sudoku_parser.hpp"

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
      }

  } else {
      cout << "No argument" << endl;
  }
}