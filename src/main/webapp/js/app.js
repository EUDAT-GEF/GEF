'use strict';

var gef = angular.module('gef', ['ui.bootstrap', 'gef.filters', 'gef.services', 'gef.directives', 'gef.controllers']);
gef.config(['$routeProvider', function($routeProvider) {
        $routeProvider.when('/login', {templateUrl: 'partials/login.html', controller: 'Login'});
        $routeProvider.when('/datasets', {templateUrl: 'partials/datasets.html', controller: 'Datasets'});
        $routeProvider.when('/newdataset', {templateUrl: 'partials/newdataset.html', controller: 'NewDataset'});
        $routeProvider.when('/workflows', {templateUrl: 'partials/workflows.html', controller: 'Workflows'});
        $routeProvider.when('/newworkflow', {templateUrl: 'partials/newworkflow.html', controller: 'NewWorkflow'});
        $routeProvider.when('/jobs', {templateUrl: 'partials/jobs.html', controller: 'Jobs'});
        $routeProvider.when('/jobs/:jobId', {templateUrl: 'partials/jobs.html', controller: 'Jobs'});
        $routeProvider.when('/newjob', {templateUrl: 'partials/newjob.html', controller: 'NewJob'});
        $routeProvider.when('/about', {templateUrl: 'partials/about.html', controller: 'About'});
        $routeProvider.otherwise({redirectTo: 'datasets'});
    }]);
