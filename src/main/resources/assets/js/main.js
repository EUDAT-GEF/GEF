/** @jsx React.DOM */
(function() {
"use strict";

var VERSION = "0.3-beta-5";
var PT = React.PropTypes;
var ErrorPane = window.MyReact.ErrorPane;
var Files = window.MyReact.Files;

window.MyGEF = window.MyGEF || {};

var Main = React.createClass({displayName: "Main",
	getInitialState: function () {
		return {
			navbarCollapse: false,
			navbarPageFn: this.renderMain,
			errorMessages: [],
		};
	},

	error: function(errObj) {
		var err = "";
		if (typeof errObj === 'string' || errObj instanceof String) {
			err = errObj;
		} else if (typeof errObj === 'object' && errObj.statusText) {
			console.log("ERROR: jqXHR = ", errObj);
			err = errObj.statusText;
		} else {
			return;
		}

		var that = this;
		var errs = this.state.errorMessages.slice();
		errs.push(err);
		this.setState({errorMessages: errs});

		setTimeout(function() {
			var errs = that.state.errorMessages.slice();
			errs.shift();
			that.setState({errorMessages: errs});
		}, 10000);
	},

	ajax: function(ajaxObject) {
		var that = this;
		if (!ajaxObject.error) {
			ajaxObject.error = function(jqXHR, textStatus, error) {
				if (jqXHR.readyState === 0) {
					that.error("Network error, please check your internet connection");
				} else if (jqXHR.responseText) {
					that.error(jqXHR.responseText + " ("+error+")");
				} else  {
					that.error(error + " ("+textStatus+")");
				}
				console.log("ajax error, jqXHR: ", jqXHR);
			};
		}
		// console.log("ajax", ajaxObject);
		jQuery.ajax(ajaxObject);
	},

	toggleCollapse: function() {
		this.setState({navbarCollapse: !this.state.navbarCollapse});
	},

	setNavbarPageFn: function(pageFn) {
		this.setState({navbarPageFn:pageFn});
	},

	renderCollapsible: function() {
		var classname = "navbar-collapse collapse " + (this.state.navbarCollapse?"in":"");
		return (
			React.createElement("div", {className: classname}, 
				React.createElement("ul", {className: "nav navbar-nav"}, 
					React.createElement("li", {className: this.state.navbarPageFn === this.renderMain ? "active":""}, 
						React.createElement("a", {className: "link", tabIndex: "-1", 
							onClick: this.setNavbarPageFn.bind(this, this.renderMain)}, "Main")
					)
				), 
				React.createElement("ul", {className: "nav navbar-nav navbar-right"}, 
					React.createElement("li", null, 
						React.createElement("a", {href: "login", tabIndex: "-1"}, 
							React.createElement("span", {className: "glyphicon glyphicon-user"})
						)
					)
				)
			)
		);
	},

	renderMain: function() {
		var progress=0;
		return (
			React.createElement("div", null, 
				React.createElement("div", {className: "row"}, 
					React.createElement("h3", null, "Add new dataset"), 
					React.createElement("p", null, "Please select and upload all the files in your dataset"), 
					React.createElement(Files, {apiURL: "api/datasets", error: this.error})
				)
			)
		);
	},

	render: function() {
		return	(
			React.createElement("div", null, 
				React.createElement("div", {className: "navbar navbar-default navbar-static-top", role: "navigation"}, 
					React.createElement("div", {className: "container"}, 
						React.createElement("div", {className: "navbar-header"}, 
							React.createElement("button", {type: "button", className: "navbar-toggle", onClick: this.toggleCollapse}, 
								React.createElement("span", {className: "sr-only"}, "Toggle navigation"), 
								React.createElement("span", {className: "icon-bar"}), 
								React.createElement("span", {className: "icon-bar"}), 
								React.createElement("span", {className: "icon-bar"})
							), 
							React.createElement("a", {className: "navbar-brand", href: "#", tabIndex: "-1"}, React.createElement("header", null, "GEF"))
						), 
						this.renderCollapsible()
					)
				), 

				React.createElement(ErrorPane, {errorMessages: this.state.errorMessages}), 

				React.createElement("div", {id: "push"}, 
					React.createElement("div", {className: "container"}, 
						this.state.navbarPageFn()
					), 
					React.createElement("div", {className: "top-gap"})
				)
			)
		);
	}
});

var Footer = React.createClass({displayName: "Footer",
	about: function(e) {
		main.about();
		e.preventDefault();
		e.stopPropagation();
	},

	render: function() {
		return	(
			React.createElement("div", {className: "container", style: {borderTop:"1px solid #ddd", paddingTop:5}}, 
				React.createElement("div", {className: "row"}, 
					React.createElement("div", {className: "col-md-2 col-md-offset-10"}, 
						React.createElement("a", {title: "about", href: "#", onClick: this.about}, 
							React.createElement("span", {className: "glyphicon glyphicon-info-sign"}), 
							React.createElement("span", null, " v.", VERSION)
						)
					)
				)
			)
		);
	}
});

var main = React.render(React.createElement(Main, null),  document.getElementById('body'));
React.render(React.createElement(Footer, null), document.getElementById('footer') );
window.MyGEF.main = main;

})();
