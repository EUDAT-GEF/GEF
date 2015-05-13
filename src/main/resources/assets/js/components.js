/** @jsx React.DOM */
(function() {
"use strict";

var PT = React.PropTypes;
var ReactCSSTransitionGroup = React.addons.ReactCSSTransitionGroup;
var ReactTransitionGroup = React.addons.TransitionGroup;

window.MyReact = {};

///////////////////////////////////////////////////////////////////////////////
// Slides

var JQuerySlide = React.createClass({displayName: "JQuerySlide",
	componentWillEnter: function(callback){
		var el = jQuery(this.getDOMNode());
		el.css("display", "none");
		el.slideDown(500, callback);
		$el.slideDown(function(){
			callback();
		});
	},
	componentWillLeave: function(callback){
		var $el = jQuery(this.getDOMNode());
		$el.slideUp(function(){
			callback();
		});
	},
	render: function(){
		return this.transferPropsTo(this.props.component({style: {display: 'none'}}));
	}
});
window.MyReact.JQuerySlide = JQuerySlide;

var JQueryFade = React.createClass({displayName: "JQueryFade",
	componentWillEnter: function(callback){
		var el = jQuery(this.getDOMNode());
		el.css("display", "none");
		el.fadeIn(500, callback);
	},
	componentWillLeave: function(callback){
		jQuery(this.getDOMNode()).fadeOut(500, callback);
	},
	render: function() {
		return this.props.children;
	}
});
window.MyReact.JQueryFade = JQueryFade;

///////////////////////////////////////////////////////////////////////////////
// Error Pane

window.MyReact.ErrorPane = React.createClass({displayName: "ErrorPane",
	propTypes: {
		errorMessages: PT.array.isRequired,
	},

	renderErrorMessage: function(errorMessage, index) {
		return errorMessage ?
			React.createElement(JQueryFade, {key: index}, 
				React.createElement("div", {key: index, className: "errorMessage"}, errorMessage)
			) :
			false;
	},

	render: function() {
		return	React.createElement("div", {className: "container errorDiv"}, 
					React.createElement("div", {className: "row errorRow"}, 
						React.createElement(ReactTransitionGroup, {component: "div"}, 
							this.props.errorMessages.map(this.renderErrorMessage)
						)
					)
				);
	}
});

///////////////////////////////////////////////////////////////////////////////
// Modal


window.MyReact.Modal = React.createClass({displayName: "Modal",
	propTypes: {
		title: PT.string.isRequired,
	},
	componentDidMount: function() {
		$(this.getDOMNode()).modal({background: true, keyboard: true, show: false});
	},
	componentWillUnmount: function() {
		$(this.getDOMNode()).off('hidden');
	},
	handleClick: function(e) {
		e.stopPropagation();
	},
	render: function() {
		return (
			React.createElement("div", {onClick: this.handleClick, className: "modal fade", role: "dialog", "aria-hidden": "true"}, 
				React.createElement("div", {className: "modal-dialog"}, 
					React.createElement("div", {className: "modal-content"}, 
						React.createElement("div", {className: "modal-header"}, 
							React.createElement("button", {type: "button", className: "close", "data-dismiss": "modal"}, 
								React.createElement("span", {"aria-hidden": "true"}, "×"), 
								React.createElement("span", {className: "sr-only"}, "Close")
							), 
							React.createElement("h2", {className: "modal-title"}, this.props.title)
						), 
						React.createElement("div", {className: "modal-body"}, 
							this.props.children
						), 
						React.createElement("div", {className: "modal-footer"}, 
							React.createElement("button", {type: "button", className: "btn btn-default", "data-dismiss": "modal"}, "Close")
						)
					)
				)
			)
		);
	}
});

///////////////////////////////////////////////////////////////////////////////
// Files

window.MyReact.Files = React.createClass({displayName: "Files",
	propTypes: {
		apiURL: PT.string.isRequired,
	},

	getInitialState: function() {
		return {
			uploading: false,
			progress: 0,
			files: [],
		}
	},

	handleBrowse: function(event) {
		this.refs.fileinput.getDOMNode().click();
	},

	handleAdd: function(event) {
		var files = this.state.files;
		var fs = event.target.files;
		for (var i = 0, length = fs.length; i < length; ++i) {
			files.push(fs[i]);
			console.log("added", fs[i]);
		}
		this.setState({files:files});

	},

	handleSubmit: function(e) {
		e.preventDefault();

		var fs = this.state.files;
		if (!fs || !fs.length) {
			return;
		}

		var fd = new FormData();
		for (var i = 0, length = fs.length; i < length; ++i) {
			fd.append("file", fs[i]);
		}
		var xhr = new XMLHttpRequest();
		xhr.upload.addEventListener("progress", this.uploadProgress, false);
		xhr.addEventListener("load", this.uploadComplete, false);
		xhr.addEventListener("error", this.uploadFailed, false);
		xhr.addEventListener("abort", this.uploadCanceled, false);
		this.setState({uploading:true});
		xhr.open("POST", this.props.apiURL);
		xhr.send(fd);
	},

	uploadProgress: function(){},
	uploadComplete: function(){},
	uploadFailed: function(){},
	uploadCanceled: function(){},

	renderFile: function(f){
		return (
			React.createElement("div", {key: f.name, style: {margin:"5px 0"}}, 
				React.createElement("input", {type: "text", className: "input-large", readOnly: "readonly", 
					 style: {width:'100%'}, value: f.name})
			)
		);
	},

	render: function() {
		var display = this.state.uploading ? 'block' : 'none';
		return (
			React.createElement("div", {style: {width:300}}, 
				React.createElement("form", {className: "form-horizontal", style: {margin:0}, 
						encType: "multipart/form-data", onSubmit: this.handleSubmit}, 
					React.createElement("div", null, 
						React.createElement("input", {ref: "fileinput", name: "file", type: "file", multiple: true, onChange: this.handleAdd, 
							style: {width:0, height:0, margin:0}}), 
						React.createElement("button", {type: "button", className: "btn", onClick: this.handleBrowse}, 
							React.createElement("i", {className: "glyphicon glyphicon-folder-open"}), " Browse"
						)
					), 
					 this.state.files.map(this.renderFile), 
					 this.state.files.length ?
					React.createElement("div", {style: {width:'100%'}}, 
						React.createElement("div", {style: {width:'50%', float:'left', display: display}}, 
							React.createElement("div", {className: "bar", style: {width: this.state.progress+'%'}}), 
							React.createElement("div", {className: "percent"}, this.state.progress, "%")
						), 
						React.createElement("button", {type: "submit", className: "btn btn-primary", style: {width:'50%', float:'right'}}, 
							React.createElement("i", {className: "glyphicon glyphicon-upload"}), " Upload"
						)
					) : false
					
				)
			)
		);
	},
});

})();
