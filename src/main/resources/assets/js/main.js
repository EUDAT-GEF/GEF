/** @jsx React.DOM */
(function() {
"use strict";

var VERSION = "0.3.6";
var PT = React.PropTypes;
var ErrorPane = window.MyReact.ErrorPane;
var Files = window.MyReact.Files;

window.MyGEF = window.MyGEF || {};

var apiRootName = "/gef/api";
var apiNames = {
	datasets: apiRootName+"/datasets",
};

function setState(state) {
	if (this && this != window && this.setState) {
		this.setState(state);
	}
}

var Main = React.createClass({displayName: "Main",
	getInitialState: function () {
		return {
			page: this.browseDatasets,
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

	browseDatasets: function() {
		return (
			React.createElement(BrowseDatasets, {error: this.error, ajax: this.ajax})
		);
	},

	executeService: function() {
		return (
			React.createElement(ExecuteService, {error: this.error, ajax: this.ajax})
		);
	},

	runningJobs: function() {
		return (
			React.createElement(RunningJobs, {error: this.error, ajax: this.ajax})
		);
	},

	createDataset: function() {
		return (
			React.createElement(CreateDataset, {error: this.error, ajax: this.ajax})
		);
	},

	createService: function() {
		return (
			React.createElement(CreateService, {error: this.error, ajax: this.ajax})
		);
	},

	renderToolListItem: function(pageFn, title) {
		var klass = "list-group-item " + (pageFn === this.state.page ? "active":"");
		return (
			React.createElement("a", {href: "#", className: klass, onClick: setState.bind(this, {page:pageFn})}, 
				title
			)
		);
	},

	render: function() {
		return	(
			React.createElement("div", null, 
				React.createElement(ErrorPane, {errorMessages: this.state.errorMessages}), 
				React.createElement("div", {className: "container"}, 
					React.createElement("div", {className: "row"}, 
						React.createElement("div", {className: "col-xs-12 col-sm-2 col-md-2"}, 
							React.createElement("div", {className: "list-group"}, 
								this.renderToolListItem(this.createService, "Create Service"), 
								this.renderToolListItem(this.executeService, "Execute Service"), 
								this.renderToolListItem(this.runningJobs, "Browse Jobs")
							), 
							React.createElement("div", {className: "list-group"}, 
								this.renderToolListItem(this.browseDatasets, "Browse Datasets")
							)
						), 
						React.createElement("div", {className: "col-xs-12 col-sm-10 col-md-10"}, 
							 this.state.page ? this.state.page() : false
						)
					)
				)
			)
		);
	}
});

///////////////////////////////////////////////////////////////////////////////

function humanSize(sz) {
	if (sz < 1024) {
		return [sz,"B  "];
	} else if (sz < 1024 * 1024) {
		return [(sz/1024).toFixed(1), "KiB"];
	} else if (sz < 1024 * 1024 * 1024) {
		return [(sz/(1024*1024)).toFixed(1), "MiB"];
	} else if (sz < 1024 * 1024 * 1024 * 1024) {
		return [(sz/(1024*1024*1024)).toFixed(1), "GiB"];
	} else {
		return [(sz/(1024*1024*1024*1024)).toFixed(1), "TiB"];
	}
}

///////////////////////////////////////////////////////////////////////////////

var CreateService = React.createClass({displayName: "CreateService",
	render: function() {
		return (
			React.createElement("div", null, 
				React.createElement("h3", null, " Create Service "), 
				React.createElement("ul", null, 
					React.createElement("li", null, "Select base image"), 
					React.createElement("li", null, "Upload files"), 
					React.createElement("li", null, "Define inputs and outputs"), 
					React.createElement("li", null, "Execute command"), 
					React.createElement("li", null, "Test data"), 
					React.createElement("li", null, "Create")
				)
			)
		);
	},
});

///////////////////////////////////////////////////////////////////////////////

var ExecuteService = React.createClass({displayName: "ExecuteService",
	props: {
		error: PT.func.isRequired,
		ajax: PT.func.isRequired,
	},

	getInitialState: function() {
		return {
		};
	},

	render: function() {
		return (
			React.createElement("div", {className: "execute-service-page"}, 
				React.createElement("h3", null, " Execute Service ")
			)
		);
	}
});

///////////////////////////////////////////////////////////////////////////////

var RunningJobs = React.createClass({displayName: "RunningJobs",
	render: function() {
		return (
			React.createElement("div", null, 
				React.createElement("h3", null, " RunningJobs ")
			)
		);
	},
});

///////////////////////////////////////////////////////////////////////////////

var CreateDataset = React.createClass({displayName: "CreateDataset",
	props: {
		error: PT.func.isRequired,
		ajax: PT.func.isRequired,
	},

	render: function() {
		return (
			React.createElement("div", null, 
				React.createElement("h3", null, " Create Dataset "), 
				React.createElement("p", null, "Please select and upload all the files in your dataset"), 
				React.createElement(Files, {apiURL: apiNames.datasets, error: this.props.error, 
						cancel: function(){}})
			)
		);
	},
});

///////////////////////////////////////////////////////////////////////////////

var BrowseDatasets = React.createClass({displayName: "BrowseDatasets",
	props: {
		error: PT.func.isRequired,
 		ajax: PT.func.isRequired,
	},

	getInitialState: function() {
		return {
			datasets: [],
		};
	},

	componentDidMount: function() {
		this.props.ajax({
			url: apiNames.datasets,
			success: function(json, textStatus, jqXHR) {
				if (!this.isMounted()) {
					return;
				}
				if (!json.datasets) {
					this.props.error("Didn't get json datasets from server");
					return;
				}
				this.setState({datasets: json.datasets});
			}.bind(this),
		});
	},

	renderHeads: function(dataset) {
		return (
			React.createElement("div", {className: "row table-head"}, 
				React.createElement("div", {className: "col-xs-12 col-sm-5 col-md-5"}, "ID"), 
				React.createElement("div", {className: "col-xs-12 col-sm-2 col-md-2", style: {textAlign:'right'}}, "Size"), 
				React.createElement("div", {className: "col-xs-12 col-sm-5 col-md-5", style: {textAlign:'right'}}, "Date")
			)
		);
	},

	renderDataset: function(dataset) {
		var sz = humanSize(dataset.entry.size);
		return (
			React.createElement("div", {className: "row"}, 
				React.createElement("div", {key: dataset.id}, 
					React.createElement("div", {className: "col-xs-12 col-sm-5 col-md-5"}, dataset.id), 
					React.createElement("div", {className: "col-xs-12 col-sm-2 col-md-2", style: {textAlign:'right'}}, sz[0], " ", sz[1]), 
					React.createElement("div", {className: "col-xs-12 col-sm-5 col-md-5", style: {textAlign:'right'}}, new Date(dataset.entry.date).toLocaleString())
				)
			)
		);
	},

	render: function() {
		return (
			React.createElement("div", {className: "dataset-page"}, 
				React.createElement("h3", null, " Browse Datasets "), 
				 this.renderHeads(), 
				React.createElement("div", {className: "dataset-table"}, 
					 this.state.datasets.map(this.renderDataset) 
				)
			)
		);
	}
});

///////////////////////////////////////////////////////////////////////////////

var Footer = React.createClass({displayName: "Footer",
	about: function(e) {
		main.about();
		e.preventDefault();
		e.stopPropagation();
	},

	render: function() {
		return	(
			React.createElement("div", {className: "container"}, 
				React.createElement("div", {className: "row"}, 
					React.createElement("div", {className: "col-xs-12 col-sm-6 col-md-6"}, 
						React.createElement("p", null, " ", React.createElement("img", {width: "45", height: "31", src: "images/flag-ce.jpg", style: {float:'left', marginRight:10}}), 
							"EUDAT receives funding from the European Union’s Horizon 2020 research" + ' ' +
							"and innovation programme under grant agreement No. 654065. ", 
							React.createElement("a", {href: "#"}, "Legal Notice"), "."
						)
					), 
					React.createElement("div", {className: "col-xs-12 col-sm-6 col-md-6 text-right"}, 
						React.createElement("ul", {className: "list-inline pull-right", style: {marginLeft:20}}, 
							React.createElement("li", null, React.createElement("span", {style: {color:'#173b93', fontWeight:'500'}}, " GEF v.", VERSION))
						), 
						React.createElement("ul", {className: "list-inline pull-right"}, 
							React.createElement("li", null, React.createElement("a", {target: "_blank", href: "http://eudat.eu/what-eudat"}, "About EUDAT")), 
							React.createElement("li", null, React.createElement("a", {href: "https://github.com/GEFx"}, "Go to GitHub")), 
							React.createElement("li", null, React.createElement("a", {href: "mailto:emanuel.dima@uni-tuebingen.de"}, "Contact"))
						)
					)
				)
			)
		);
	}
});

window.MyGEF.main = React.render(React.createElement(Main, null),  document.getElementById('page'));
window.MyGEF.footer = React.render(React.createElement(Footer, null), document.getElementById('footer') );

})();
