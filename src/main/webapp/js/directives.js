'use strict';

var gefdirectives = angular.module('gef.directives', []);

gefdirectives.directive('appVersion', ['version', function(version) {
		return function(scope, elm, attrs) {
			elm.text(version);
		};
	}]);
