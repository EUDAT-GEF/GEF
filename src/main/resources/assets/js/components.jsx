/** @jsx React.DOM */
(function() {
"use strict";

var PT = React.PropTypes;
var ReactCSSTransitionGroup = React.addons.ReactCSSTransitionGroup;
var ReactTransitionGroup = React.addons.TransitionGroup;

window.MyReact = {};

///////////////////////////////////////////////////////////////////////////////
// Slides

var JQuerySlide = React.createClass({
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

var JQueryFade = React.createClass({
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

window.MyReact.ErrorPane = React.createClass({
	propTypes: {
		errorMessages: PT.array.isRequired,
	},

	renderErrorMessage: function(errorMessage, index) {
		return errorMessage ?
			<JQueryFade key={index}>
				<div key={index} className="errorMessage">{errorMessage}</div>
			</JQueryFade> :
			false;
	},

	render: function() {
		return	<div className="container errorDiv">
					<div className="row errorRow">
						<ReactTransitionGroup component="div">
							{this.props.errorMessages.map(this.renderErrorMessage)}
						</ReactTransitionGroup>
					</div>
				</div>;
	}
});

///////////////////////////////////////////////////////////////////////////////
// Modal


window.MyReact.Modal = React.createClass({
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
			<div onClick={this.handleClick} className="modal fade" role="dialog" aria-hidden="true">
				<div className="modal-dialog">
					<div className="modal-content">
						<div className="modal-header">
							<button type="button" className="close" data-dismiss="modal">
								<span aria-hidden="true">&times;</span>
								<span className="sr-only">Close</span>
							</button>
							<h2 className="modal-title">{this.props.title}</h2>
						</div>
						<div className="modal-body">
							{this.props.children}
						</div>
						<div className="modal-footer">
							<button type="button" className="btn btn-default" data-dismiss="modal">Close</button>
						</div>
					</div>
				</div>
			</div>
		);
	}
});

///////////////////////////////////////////////////////////////////////////////
// Files

window.MyReact.Files = React.createClass({
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
			<div key={f.name} style={{margin:"5px 0"}}>
				<input type="text" className="input-large" readOnly="readonly"
					 style={{width:'100%'}} value={f.name}/>
			</div>
		);
	},

	render: function() {
		var display = this.state.uploading ? 'block' : 'none';
		return (
			<div style={{width:300}}>
				<form className="form-horizontal" style={{margin:0}}
						encType="multipart/form-data" onSubmit={this.handleSubmit}>
					<div>
						<input ref='fileinput' name="file" type="file" multiple onChange={this.handleAdd}
							style={{width:0, height:0, margin:0}} />
						<button type="button" className="btn" onClick={this.handleBrowse}>
							<i className="glyphicon glyphicon-folder-open"/> Browse
						</button>
					</div>
					{ this.state.files.map(this.renderFile) }
					{ this.state.files.length ?
					<div style={{width:'100%'}}>
						<div style={{width:'50%', float:'left', display: display}}>
							<div className="bar" style={{width: this.state.progress+'%'}} />
							<div className="percent">{this.state.progress}%</div>
						</div>
						<button type="submit" className="btn btn-primary" style={{width:'50%', float:'right'}}>
							<i className="glyphicon glyphicon-upload"/> Upload
						</button>
					</div> : false
					}
				</form>
			</div>
		);
	},
});

})();
