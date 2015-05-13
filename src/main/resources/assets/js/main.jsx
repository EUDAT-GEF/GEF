/** @jsx React.DOM */
(function() {
"use strict";

var VERSION = "0.3-beta-5";
var PT = React.PropTypes;
var ErrorPane = window.MyReact.ErrorPane;
var Files = window.MyReact.Files;

window.MyGEF = window.MyGEF || {};

var Main = React.createClass({
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
			<div className={classname}>
				<ul className="nav navbar-nav">
					<li className={this.state.navbarPageFn === this.renderMain ? "active":""}>
						<a className="link" tabIndex="-1"
							onClick={this.setNavbarPageFn.bind(this, this.renderMain)}>Main</a>
					</li>
				</ul>
				<ul className="nav navbar-nav navbar-right">
					<li>
						<a href="login" tabIndex="-1">
							<span className="glyphicon glyphicon-user"/>
						</a>
					</li>
				</ul>
			</div>
		);
	},

	renderMain: function() {
		var progress=0;
		return (
			<div>
				<div className="row">
					<h3>Add new dataset</h3>
					<p>Please select and upload all the files in your dataset</p>
					<Files apiURL="api/datasets" error={this.error}/>
				</div>
			</div>
		);
	},

	render: function() {
		return	(
			<div>
				<div className="navbar navbar-default navbar-static-top" role="navigation">
					<div className="container">
						<div className="navbar-header">
							<button type="button" className="navbar-toggle" onClick={this.toggleCollapse}>
								<span className="sr-only">Toggle navigation</span>
								<span className="icon-bar"></span>
								<span className="icon-bar"></span>
								<span className="icon-bar"></span>
							</button>
							<a className="navbar-brand" href="#" tabIndex="-1"><header>GEF</header></a>
						</div>
						{this.renderCollapsible()}
					</div>
				</div>

				<ErrorPane errorMessages={this.state.errorMessages} />

				<div id="push">
					<div className="container">
						{this.state.navbarPageFn()}
					</div>
					<div className="top-gap" />
				</div>
			</div>
		);
	}
});

var Footer = React.createClass({
	about: function(e) {
		main.about();
		e.preventDefault();
		e.stopPropagation();
	},

	render: function() {
		return	(
			<div className="container" style={{borderTop:"1px solid #ddd", paddingTop:5}}>
				<div className="row">
					<div className="col-md-2 col-md-offset-10">
						<a title="about" href="#" onClick={this.about}>
							<span className="glyphicon glyphicon-info-sign"></span>
							<span> v.{VERSION}</span>
						</a>
					</div>
				</div>
			</div>
		);
	}
});

var main = React.render(<Main />,  document.getElementById('body'));
React.render(<Footer />, document.getElementById('footer') );
window.MyGEF.main = main;

})();
