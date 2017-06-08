(function(){

  var sudokuFileInput;
  var dropZone;
  var solutionRow;
  var solutionGrid;

  // Initialize handlers when DOM is loaded
  document.addEventListener("DOMContentLoaded", function() {
    // important DOM elements
    sudokuFileInput = document.getElementById("sudokuFile");
    dropZone = document.getElementById("dropZone");
    solutionRow = document.getElementById("solution-row");
    solutionGrid = document.getElementById("solution-grid");

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
          console.log("... file[" + i + "].name = " + f.name);
        }
      }
    } else {
      // Use DataTransfer interface to access the file(s)
      for (var i=0; i < dt.files.length; i++) {
        console.log("... file[" + i + "].name = " + dt.files[i].name);
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
      if (solutionRow.className.indexOf(" loading") < 0) {
        solutionRow.className += " loading";
      }

      reader.addEventListener("load", function () {
        document.getElementById('puzzle-src').src = reader.result;
      }, false);
      reader.readAsDataURL(file);

      xhr.open("POST", uri, true);
      xhr.setRequestHeader("X-Requested-With", "XMLHttpRequest");
      xhr.onreadystatechange = function() {
          if (xhr.readyState == 4) {
              var data = JSON.parse(xhr.responseText);
              if (xhr.status == 200) {

                  // insert the solution in the DOM
                  while (solutionGrid.lastChild) { solutionGrid.removeChild(solutionGrid.lastChild); }
                  var i = 0;
                  for (var cell in data.Values) {
                    var val = data.Values[cell];
                    var cellDiv = document.createElement("div");
                    cellDiv.appendChild(document.createTextNode(val.value));
                    cellDiv.className = "solution-cell";
                    solutionGrid.appendChild(cellDiv);
                    
                    //if (i < 9) { cellDiv.className += " top"; }

                    i++;
                  }

              } else {
                  document.getElementById("solution").value = data.Error;
              }

              solutionRow.className = solutionRow.className.replace(" loading", "");
          }
      };
      fd.append('sudokuFile', file);
      // Initiate a multipart/form-data upload
      xhr.send(fd);
  }


}())



var isAdvancedUpload = function() {
    var div = document.createElement('div');
    return (('draggable' in div) || ('ondragstart' in div && 'ondrop' in div)) && 'FormData' in window && 'FileReader' in window;
}();

