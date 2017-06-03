The OpenCV API for loading an SVM model does not, at present support passing in the data.  Instead, it reads from a file path.

To simplify deployment, [go-bindata](https://github.com/jteeuwen/go-bindata) is used to embed the SVM model file in a go source file as follows:
```
go-bindata -o svm_model.go -pkg sudokuparser data/
```