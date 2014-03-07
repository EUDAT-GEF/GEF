"use strict";
ko.bindingHandlers.imageSource = {
	update : function(element, valueAccessor) {
		var value = ko.utils.unwrapObservable(valueAccessor());
		jQuery(element).attr("src", value);
	}
};

ko.bindingHandlers.fadeVisible = {
	init : function(element, valueAccessor) {
		var value = valueAccessor();
		jQuery(element).toggle(ko.utils.unwrapObservable(value));
	},
	update : function(element, valueAccessor) {
		var value = valueAccessor();
		if (ko.utils.unwrapObservable(value)) {
			jQuery(element).fadeIn(100);
		} else {
			jQuery(element).fadeOut(100);
		}
	}
};

ko.bindingHandlers.slideVisible = {
	init : function(element, valueAccessor) {
		var value = valueAccessor();
		jQuery(element).toggle(ko.utils.unwrapObservable(value));
	},
	update : function(element, valueAccessor) {
		var value = valueAccessor();
		if (ko.utils.unwrapObservable(value)) {
			jQuery(element).slideDown(200);
		} else {
			jQuery(element).slideUp(200);
		}
	}
};
