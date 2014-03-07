'use strict';

var gefcontrollers = angular.module('gef.controllers', []);

gefcontrollers.controller('Login', [function() {
    }]);

gefcontrollers.controller('Datasets', ['$scope', '$http', function($scope, $http) {
        $http.get('rest/datasets').success(function(data) {
            console.log("datasets: ", data);
            $scope.datasets = data;
        });
        $scope.orderProp = '-date';
    }]);

gefcontrollers.controller('Workflows', ['$scope', '$rootScope', '$http', '$location',
    function($scope, $rootScope, $http, $location) {
        $http.get('rest/workflows').success(function(data) {
            console.log("workflows: ", data);
            $scope.workflows = data;
        });
        $scope.workflowsOrderProp = '-date';

        $scope.run = function(workflowId, workflowName) {
            $rootScope.newjob = {
                workflowId: workflowId,
                workflowName: workflowName
            };
            $location.path('newjob');
        };
    }]);

gefcontrollers.controller('Jobs', ['$scope', '$http', '$routeParams', '$location',
    function($scope, $http, $routeParams, $location) {
        $scope.showFile = false;
        var jobId = $routeParams.jobId;
        console.log(jobId)
        $scope.title = jobId ? ("Job " + jobId) : "Jobs";
        $http.get('rest/jobs' + (jobId ? ("/" + jobId) : "")).success(function(data, x) {
            console.log("jobs: ", data, x);
            $scope.showFile = false;
            $scope.jobs = data;
        });
        $scope.jobsOrderProp = '-date';

        $scope.goto = function(location) {
            if (jobId) {
                $http.get('rest/jobs/' + jobId + "/" + location)
                        .success(function(data, x) {
                    $scope.file = data;
                    $scope.showFile = true;
                });
            } else {
                $location.path('jobs/' + location);
            }
        };
    }]);

gefcontrollers.controller('NewDataset', ['$scope', '$http', '$location', function($scope, $http, $location) {
        $scope.files = [];

        var fileInput = angular.element(document.getElementById("fileInput"));
        fileInput.bind('change', function(event) {
            $scope.$apply(function() {
                $scope.files = event.target.files;
            });
        });
        angular.element(document.getElementById("browse")).bind('click', function() {
            fileInput[0].click();
        });

        $scope.$watchCollection('files', function(newValue, oldValue) {
            if (newValue === oldValue) {
                return;
            }
            console.log("selected files:", $scope.files);
            document.getElementById("filename").value =
                    newValue.length === 0 ? "" :
                    newValue.length === 1 ? newValue[0].name : (newValue.length + ' files');
        });

        function uploadProgress(evt) {
            $scope.$apply(function() {
                if (evt.lengthComputable) {
                    $scope.progress = Math.round(evt.loaded * 100 / evt.total);
                } else {
                    $scope.progress = 'unable to compute';
                }
            });
        }

        function uploadComplete(evt) {
            console.log(evt.target.responseText);
            $scope.$apply(function() {
                $scope.progressVisible = false;
                $location.path('datasets');
            });
        }

        function uploadFailed() {
            console.log("There was an error attempting to upload the file.");
            $scope.$apply(function() {
                $scope.progressVisible = false;
            });
        }

        function uploadCanceled() {
            $scope.$apply(function() {
                $scope.progressVisible = false;
            });
            console.log("The upload has been canceled by the user or the browser dropped the connection.");
        }

        $scope.create = function() {
            var fd = new FormData();
            for (var i = 0, length = $scope.files.length; i < length; ++i) {
                fd.append("file", $scope.files[i]);
            }
            var xhr = new XMLHttpRequest();
            xhr.upload.addEventListener("progress", uploadProgress, false);
            xhr.addEventListener("load", uploadComplete, false);
            xhr.addEventListener("error", uploadFailed, false);
            xhr.addEventListener("abort", uploadCanceled, false);
            xhr.open("POST", "rest/datasets");
            $scope.progressVisible = true;
            xhr.send(fd);
        };
    }]);

gefcontrollers.controller('NewWorkflow', ['$scope', '$http', '$location', function($scope, $http, $location) {
        $scope.files = [];

        var fileInput = angular.element(document.getElementById("fileInput"));
        fileInput.bind('change', function(event) {
            $scope.$apply(function() {
                $scope.files = event.target.files;
            });
        });
        angular.element(document.getElementById("browse")).bind('click', function() {
            fileInput[0].click();
        });

        $scope.$watchCollection('files', function(newValue, oldValue) {
            if (newValue === oldValue) {
                return;
            }
            console.log("selected files:", $scope.files);
            document.getElementById("filename").value =
                    newValue.length === 0 ? "" :
                    newValue.length === 1 ? newValue[0].name : (newValue.length + ' files');
        });

        function uploadProgress(evt) {
            $scope.$apply(function() {
                if (evt.lengthComputable) {
                    $scope.progress = Math.round(evt.loaded * 100 / evt.total);
                } else {
                    $scope.progress = 'unable to compute';
                }
            });
        }

        function uploadComplete(evt) {
            console.log(evt.target.responseText);
            $scope.$apply(function() {
                $scope.progressVisible = false;
                $location.path('workflows');
            });
        }

        function uploadFailed() {
            console.log("There was an error attempting to upload the file.");
            $scope.$apply(function() {
                $scope.progressVisible = false;
            });
        }

        function uploadCanceled() {
            $scope.$apply(function() {
                $scope.progressVisible = false;
            });
            console.log("The upload has been canceled by the user or the browser dropped the connection.");
        }

        $scope.create = function() {
            var fd = new FormData();
            for (var i = 0, length = $scope.files.length; i < length; ++i) {
                fd.append("file", $scope.files[i]);
            }
            var xhr = new XMLHttpRequest();
            xhr.upload.addEventListener("progress", uploadProgress, false);
            xhr.addEventListener("load", uploadComplete, false);
            xhr.addEventListener("error", uploadFailed, false);
            xhr.addEventListener("abort", uploadCanceled, false);
            xhr.open("POST", "rest/workflows");
            $scope.progressVisible = true;
            xhr.send(fd);
        };
    }]);


gefcontrollers.controller('NewJob', ['$scope', '$rootScope', '$http', '$location',
    function($scope, $rootScope, $http, $location) {
        $scope.workflowId = $rootScope.newjob.workflowId;
        $scope.workflowName = $rootScope.newjob.workflowName;
        $scope.datasetId = "";
        $scope.orderProp = '-date';
        $scope.doShowDatasets = false;
        $scope.runError = "";
        $http.get('rest/datasets').success(function(data) {
            console.log("datasets: ", data);
            $scope.datasets = data;
        });

        $scope.showDatasets = function() {
            $scope.doShowDatasets = true;
        };
        $scope.hideDatasets = function() {
            $scope.doShowDatasets = false;
        };
        $scope.select = function(did, dname) {
            $scope.datasetId = did;
            $scope.datasetName = dname;
            $scope.doShowDatasets = false;
            $scope.runError = "";
        };

        function succeeded(ev) {
            console.log("succeeded", ev);
            $scope.$apply(function() {
                if (ev && ev.target && ev.target.status
                        && (200 <= ev.target.status && ev.target.status < 300)) {
                    $location.path('jobs/' + ev.target.responseText);
                } else {
                    $scope.runError = "Execution request failed: " + ev.target.statusText + ": " + ev.target.responseText;
                }
            });
        }
        function failed(ev) {
            console.log("failed", ev);
            $scope.$apply(function() {
                $scope.runError = "Execution request failed: " + ev.target.statusText + ": " + ev.target.responseText;
            });
        }

        $scope.runJob = function() {
            console.log("try runJob", $scope.workflowId, $scope.datasetId);
            //if ($scope.workflowId && $scope.datasetId) {
            var fd = new FormData();
            fd.append("workflowPid", $scope.workflowId);
            fd.append("datasetPid", $scope.datasetId);
            var xhr = new XMLHttpRequest();
            xhr.addEventListener("load", succeeded, false);
            xhr.addEventListener("error", failed, false);
            xhr.addEventListener("abort", failed, false);
            xhr.open("POST", "rest/jobs");
            xhr.send(fd);
            $scope.runError = "Please wait...";
//            }
//            else {
//				$scope.runError = "Please fill in the fields";
//			}
        };

    }]);


gefcontrollers.controller('About', [function() {
    }]);
