/** @jsx React.DOM */
(function() {
"use strict";

var PT = React.PropTypes;


// var fileInput = angular.element(document.getElementById("fileInput"));

// fileInput.bind('change', function(event) {
//     $scope.$apply(function() {
//         $scope.files = event.target.files;
//     });
// });
// angular.element(document.getElementById("browse")).bind('click', function() {
//     fileInput[0].click();
// });

// $scope.$watchCollection('files', function(newValue, oldValue) {
//     if (newValue === oldValue) {
//         return;
//     }
//     console.log("selected files:", $scope.files);
//     document.getElementById("filename").value =
//             newValue.length === 0 ? "" :
//             newValue.length === 1 ? newValue[0].name : (newValue.length + ' files');
// });

// function uploadProgress(evt) {
//     $scope.$apply(function() {
//         if (evt.lengthComputable) {
//             $scope.progress = Math.round(evt.loaded * 100 / evt.total);
//         } else {
//             $scope.progress = 'unable to compute';
//         }
//     });
// }

// function uploadComplete(evt) {
//     console.log(evt.target.responseText);
//     $scope.$apply(function() {
//         $scope.progressVisible = false;
//         $location.path('datasets');
//     });
// }

// function uploadFailed() {
//     console.log("There was an error attempting to upload the file.");
//     $scope.$apply(function() {
//         $scope.progressVisible = false;
//     });
// }

// function uploadCanceled() {
//     $scope.$apply(function() {
//         $scope.progressVisible = false;
//     });
//     console.log("The upload has been canceled by the user or the browser dropped the connection.");
// }





})();
