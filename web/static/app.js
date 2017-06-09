(function(){

  var sudokuFileInput;
  var dropZone;
  var puzzleRow;
  var solutionGrid;
  var puzzleImg;
  var puzzleAnnotation;

  // Initialize handlers when DOM is loaded
  document.addEventListener("DOMContentLoaded", function() {
    // important DOM elements
    sudokuFileInput = document.getElementById("sudokuFile");
    dropZone = document.getElementById("dropZone");
    puzzleRow = document.getElementById("puzzle-row");
    solutionGrid = document.getElementById("solution-grid");
    puzzleImg = document.getElementById('puzzle-src');
    puzzleAnnotation = document.getElementById('puzzle-annotation');

    // Wire up file input change handler
    sudokuFileInput.addEventListener("change", function(evt) {
      if (sudokuFileInput.files.length > 0) {
        sendFile(sudokuFileInput.files[0]);
      }
    });

    // Wire up drag/drop events
    dropZone.addEventListener("drop", drop_handler);
    dropZone.addEventListener("dragenter", dragenter_handler);
    dropZone.addEventListener("dragover", dragover_handler);
    dropZone.addEventListener("dragend", dragend_handler);
  });

  function drop_handler(ev) {
    console.log("Drop");
    dropZone.className = dropZone.className.replace(' drag-active', ''); 
    ev.preventDefault();
    // If dropped items aren't files, reject them
    var dt = ev.dataTransfer;
    if (dt.items) {
      // Use DataTransferItemList interface to access the file(s)
      for (var i=0; i < dt.items.length; i++) {
        if (dt.items[i].kind == "file") {
          var f = dt.items[i].getAsFile();
          sendFile(f);
          console.log("Uploading  " + f.name + "...");
        }
      }
    } else {
      // Use DataTransfer interface to access the file(s)
      for (var i=0; i < dt.files.length; i++) {
          sendFile(f);
          console.log("Uploading  " + f.name + "...");
      }  
    }
  }

  function dragenter_handler(ev) { 
    if (dropZone.className.indexOf('drag-active') == -1) {
      dropZone.className += ' drag-active';
    }
  }

  function dragover_handler(ev) {
    // Prevent default select and drag behavior
    ev.preventDefault();
  }

  function dragend_handler(ev) {
    console.log("dragEnd");
    dropZone.className = dropZone.className.replace(' drag-active', ''); 
    // Remove all of the drag data
    var dt = ev.dataTransfer;
    if (dt.items) {
      // Use DataTransferItemList interface to remove the drag data
      for (var i = 0; i < dt.items.length; i++) {
        dt.items.remove(i);
      }
    } else {
      // Use DataTransfer interface to remove the drag data
      ev.dataTransfer.clearData();
    }
  }

  function sendFile(file) {
      var uri = "/";
      var xhr = new XMLHttpRequest();
      var fd = new FormData();
      var reader  = new FileReader();

      dropZone.className = '';
      if (puzzleRow.className.indexOf(" loading") < 0) {
        puzzleRow.className += " loading";
      }

      // clear out solution, if any
      while (solutionGrid.lastChild) { solutionGrid.removeChild(solutionGrid.lastChild); }

      reader.addEventListener("load", function () {
        puzzleImg.src = reader.result;
      }, false);
      reader.readAsDataURL(file);

      xhr.open("POST", uri, true);
      xhr.setRequestHeader("X-Requested-With", "XMLHttpRequest");
      xhr.onreadystatechange = function() {
          if (xhr.readyState == 4) {
              var data = JSON.parse(xhr.responseText);
              if (xhr.status == 200) {
                setSolution(data);
              } else {
                  document.getElementById("solution").value = data.Error;
              }

              puzzleRow.className = puzzleRow.className.replace(" loading", "");
          }
      };
      fd.append('sudokuFile', file);
      // Initiate a multipart/form-data upload
      xhr.send(fd);
  }

  function setSolution(data) {
    // insert the solution in the DOM
    var i = 0;
    for (var cell in data.Values) {
      var val = data.Values[cell];
      var cellDiv = document.createElement("div");
      if (val.source === 1) {
        cellDiv.appendChild(document.createTextNode(val.value));
      }
      cellDiv.className = "solution-cell";
      if (val.source === 0) { cellDiv.className += " parsed"; }
      solutionGrid.appendChild(cellDiv);
      
      //if (i < 9) { cellDiv.className += " top"; }

      i++;
    }

    puzzleAnnotation.src = getImageCanvas(data.Points);
    var transform = new PerspectiveTransform(
      solutionGrid, 
      solutionGrid.clientWidth, 
      solutionGrid.clientHeight, 
      true);

    for (var i = 0; i < 4; i++) {
      var targetPoint = i === 0 ? transform.topLeft : i === 1 ? transform.topRight : i === 2 ? transform.bottomRight : transform.bottomLeft;
      targetPoint.x = (data.Points[i].x / puzzleImg.naturalWidth) * solutionGrid.clientWidth;
      targetPoint.y = (data.Points[i].y / puzzleImg.naturalHeight) * solutionGrid.clientHeight;
    }
    if (transform.checkError()==0){
        transform.update();
    }
  }

  function getImageCanvas(puzzleCorners) {
    var canvas = document.createElement("canvas");
    canvas.width = puzzleImg.naturalWidth;
    canvas.height = puzzleImg.naturalHeight;
    var ctx = canvas.getContext("2d");
    ctx.fillStyle = "#FF0000";
    for(var i = 0; i < puzzleCorners.length; i++) {
      ctx.beginPath();
      ctx.arc(puzzleCorners[i].x, puzzleCorners[i].y, 10, 0, Math.PI * 2, true); // Outer circle
      ctx.fill();
    }
    return canvas.toDataURL();
  }

}())
