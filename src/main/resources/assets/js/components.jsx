/** @jsx React.DOM */
(function() {
"use strict";

var PT = React.PropTypes;
var ReactCSSTransitionGroup = React.addons.ReactCSSTransitionGroup;
var ReactTransitionGroup = React.addons.TransitionGroup;

window.MyReact = window.MyReact || {};

///////////////////////////////////////////////////////////////////////////////
// Slides

var JQuerySlide = window.MyReact.JQuerySlide = React.createClass({
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

var JQueryFade = window.MyReact.JQueryFade = React.createClass({
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
						<div className="col-md-offset-4 col-md-8">
							<ReactTransitionGroup component="div">
								{this.props.errorMessages.map(this.renderErrorMessage)}
							</ReactTransitionGroup>
						</div>
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

var FileAddButton = window.MyReact.FileAddButton = React.createClass({
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
		var input =	<input ref='__fileinput' name="file" type="file" style={noshow}
						onChange={this.doAddFiles} />;
		if (this.props.multiple) {
			input =	<input ref='__fileinput' name="file" type="file" style={noshow}
						onChange={this.doAddFiles} multiple />;
		}
		return (
			<div>
				{input}
				<button type="button" className="btn btn-default" style={{width:"100%"}} onClick={this.doBrowse}>
					{this.props.caption || "Browse"} </button>
			</div>
		);
	},
});


window.MyReact.Files = React.createClass({
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
			<div key={f.name} className="row" style={{margin:'5px 0px'}}>
				<div className="col-md-12">
					<div className="input-group">
						<input type="text" className="input-large" readOnly="readonly"
							style={{width:'100%', lineHeight:'26px', paddingLeft:10}} value={f.name}/>
						<span className="input-group-btn">
							<button type="button" className="btn btn-warning btn-sm" onClick={this.handleRemove.bind(this, f)}>
								<i className="glyphicon glyphicon-remove"/>
							</button>
						</span>
					</div>
				</div>
			</div>
		);
	},

	renderProgress: function() {
		if (this.state.uploadPercent < 0) {
			return false;
		}
		var percent = this.state.uploadPercent;
		var widthPercent = {width: percent+'%'};
		return (
			<div className="row" style={{margin:'5px 0px'}}>
				<div className="col-md-12">
					<div className="progress" style={{margin:'10px 0'}}>
						<div className="progress-bar progress-bar-info progress-bar-striped"
								role="progressbar" aria-valuenow={percent}
								aria-valuemin="0" aria-valuemax="100" style={widthPercent} />
					</div>
				</div>
			</div>
		);
	},

	renderBottomRow: function() {
		return (
			<div className="row" style={{margin:'5px 0px'}}>
				<div className="col-md-3">
					<FileAddButton multiple={true} caption="Add files" fileAddHandler={this.handleAdd} />
				</div>

				{ this.state.files.length ?
					<div className="col-md-3 col-md-offset-3">
						<button className="btn btn-default" style={{width:'100%'}} onClick={this.handleCancel}>
							Cancel
						</button>
					</div> : false }
				{ this.state.files.length ?
					<div className="col-md-3">
						<button type="submit" className="btn btn-primary" style={{width:'100%'}}>
							<i className="glyphicon glyphicon-upload"/> Upload
						</button>
					</div> : false }
			</div>
		);
	},

	render: function() {
		return (
			<div className="file-upload-form">
				<form className="form-horizontal" style={{margin:0}}
						encType="multipart/form-data" onSubmit={this.handleSubmit}>
					{ this.state.files.length ? this.state.files.map(this.renderFile) : false }
					{ this.state.files.length ? this.renderProgress() : false }
					{ this.renderBottomRow() }
				</form>
			</div>
		);
	},
});

})();
