/** @jsx React.DOM */
(function() {
"use strict";

var PT = React.PropTypes;
var ReactCSSTransitionGroup = React.addons.ReactCSSTransitionGroup;
var ReactTransitionGroup = React.addons.TransitionGroup;

window.MyReact = window.MyReact || {};

///////////////////////////////////////////////////////////////////////////////
// Slides

var JQuerySlide = window.MyReact.JQuerySlide = React.createClass({displayName: "JQuerySlide",
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

var JQueryFade = window.MyReact.JQueryFade = React.createClass({displayName: "JQueryFade",
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

var FileAddButton = window.MyReact.FileAddButton = React.createClass({displayName: "FileAddButton",
	props: {
		caption: PT.string.isRequired,
		fileAddHandler: PT.func.isRequired,
		multiple: PT.bool,
	},

	doBrowse: function(event) {
		this.refs.__fileinput.getDOMNode().click();
	},

	doAddFiles: function(event) {
		this.props.fileAddHandler(event.target.files);
	},

	render: function() {
		var noshow = {width:0, height:0, margin:0, border:'none'};
		var input =	React.createElement("input", {ref: "__fileinput", name: "file", type: "file", style: noshow, 
						onChange: this.doAddFiles});
		if (this.props.multiple) {
			input =	React.createElement("input", {ref: "__fileinput", name: "file", type: "file", style: noshow, 
						onChange: this.doAddFiles, multiple: true});
		}
		return (
			React.createElement("div", null, 
				input, 
				React.createElement("button", {type: "button", className: "btn btn-default", style: {width:"100%"}, onClick: this.doBrowse}, 
					this.props.caption || "Browse", " ")
			)
		);
	},
});


window.MyReact.Files = React.createClass({displayName: "Files",
	propTypes: {
		getApiURL: PT.func.isRequired,
		cancel: PT.func.isRequired,
		error: PT.func.isRequired,
		done: PT.func.isRequired,
	},

	getInitialState: function() {
		return {
			uploadPercent: -1, // -1 (not started), 0-100, 101 (done)
			files: [],
		}
	},

	handleAdd: function(fs) {
		console.log("adding files", fs);
		var files = this.state.files;
		for (var i = 0, length = fs.length; i < length; ++i) {
			files.push(fs[i]);
			// console.log("added", fs[i]);
		}
		this.setState({files:files});
	},

	handleSubmit: function(e) {
		// console.log("submitting", e);
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
		this.setState({uploadPercent:0});

		xhr.open("POST", this.props.getApiURL());
		xhr.onreadystatechange = function () {
			if (XMLHttpRequest.DONE != xhr.readyState) {
				return;
			}
			var codeh = xhr.status/100;
			if (codeh === 2) {
				this.uploadComplete(xhr);
			} else if (codeh === 4 || codeh === 5) {
				this.uploadFailed(xhr);
			}
		}.bind(this);
		xhr.send(fd);
	},

	handleCancel: function(e) {
		// console.log("cancel", e);
		e.preventDefault();
		this.setState(this.getInitialState());
		if (this.props.cancel) {
			this.props.cancel();
		}
	},

	handleRemove: function(f) {
		var files = this.state.files.filter(function(x){return x !== f;});
		this.setState({files:files});
	},

	uploadProgress: function(e){
		// console.log('progress', e);
		var p = e.loaded / e.total * 100;
		this.setState({uploadPercent:p});
	},
	uploadCanceled: function(e){
		// console.log('canceled', e);
		this.setState({uploadPercent:-1});
	},
	uploadFailed: function(xhr){
		// console.log('failed', xhr);
		this.setState({uploadPercent:101});
		this.props.error(xhr.responseText);
	},
	uploadComplete: function(xhr){
		// console.log('complete', xhr);
		this.setState({uploadPercent:101});
		this.props.done(xhr.responseText);
	},

	renderFile: function(f){
		return (
			React.createElement("div", {key: f.name, className: "row", style: {margin:'5px 0px'}}, 
				React.createElement("div", {className: "col-md-12"}, 
					React.createElement("div", {className: "input-group"}, 
						React.createElement("input", {type: "text", className: "input-large", readOnly: "readonly", 
							style: {width:'100%', lineHeight:'26px', paddingLeft:10}, value: f.name}), 
						React.createElement("span", {className: "input-group-btn"}, 
							React.createElement("button", {type: "button", className: "btn btn-warning btn-sm", onClick: this.handleRemove.bind(this, f)}, 
								React.createElement("i", {className: "glyphicon glyphicon-remove"})
							)
						)
					)
				)
			)
		);
	},

	renderProgress: function() {
		if (this.state.uploadPercent < 0) {
			return false;
		}
		var percent = this.state.uploadPercent;
		var widthPercent = {width: percent+'%'};
		return (
			React.createElement("div", {className: "row", style: {margin:'5px 0px'}}, 
				React.createElement("div", {className: "col-md-12"}, 
					React.createElement("div", {className: "progress", style: {margin:'10px 0'}}, 
						React.createElement("div", {className: "progress-bar progress-bar-info progress-bar-striped", 
								role: "progressbar", "aria-valuenow": percent, 
								"aria-valuemin": "0", "aria-valuemax": "100", style: widthPercent})
					)
				)
			)
		);
	},

	renderBottomRow: function() {
		return (
			React.createElement("div", {className: "row", style: {margin:'5px 0px'}}, 
				React.createElement("div", {className: "col-md-3"}, 
					React.createElement(FileAddButton, {multiple: true, caption: "Add files", fileAddHandler: this.handleAdd})
				), 

				 this.state.files.length ?
					React.createElement("div", {className: "col-md-3 col-md-offset-3"}, 
						React.createElement("button", {className: "btn btn-default", style: {width:'100%'}, onClick: this.handleCancel}, 
							"Cancel"
						)
					) : false, 
				 this.state.files.length ?
					React.createElement("div", {className: "col-md-3"}, 
						React.createElement("button", {type: "submit", className: "btn btn-primary", style: {width:'100%'}}, 
							React.createElement("i", {className: "glyphicon glyphicon-upload"}), " Upload"
						)
					) : false
			)
		);
	},

	render: function() {
		return (
			React.createElement("div", {className: "file-upload-form"}, 
				React.createElement("form", {className: "form-horizontal", style: {margin:0}, 
						encType: "multipart/form-data", onSubmit: this.handleSubmit}, 
					 this.state.files.length ? this.state.files.map(this.renderFile) : false, 
					 this.state.files.length ? this.renderProgress() : false, 
					 this.renderBottomRow() 
				)
			)
		);
	},
});

})();