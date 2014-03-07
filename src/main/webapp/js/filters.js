'use strict';

/* Filters */

var geffilters = angular.module('gef.filters', []);
geffilters.filter('interpolate', ['version', function(version) {
		return function(text) {
			return String(text).replace(/\%VERSION\%/mg, version);
		};
	}]);
