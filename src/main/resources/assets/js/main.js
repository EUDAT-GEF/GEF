/** @jsx React.DOM */
(function() {
"use strict";

var VERSION = "0.3.6";
var PT = React.PropTypes;
var ErrorPane = window.MyReact.ErrorPane;
var Files = window.MyReact.Files;

window.MyGEF = window.MyGEF || {};

function setState(state) {
	if (this && this != window) {
		console.log(this, state);
		this.setState(state);
	}
}

var Main = React.createClass({displayName: "Main",
	getInitialState: function () {
		return {
			page: this.datasets,
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

	datasets: function() {
		return (
			React.createElement(Datasets, {error: this.error, ajax: this.ajax})
		);
	},

	workflows: function() {
		return (
			React.createElement(Workflows, {error: this.error, ajax: this.ajax})
		);
	},

	jobs: function() {
		return (
			React.createElement(Jobs, {error: this.error, ajax: this.ajax})
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
								this.renderToolListItem(this.datasets, "Datasets"), 
								this.renderToolListItem(this.workflows, "Workflows"), 
								this.renderToolListItem(this.jobs, "Jobs")
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

var Datasets = React.createClass({displayName: "Datasets",
	props: {
		error: PT.func.isRequired,
		ajax: PT.func.isRequired,
	},

	getInitialState: function() {
		return {
			addNewPaneOpen: false,
			datasets: [],
		};
	},

	componentDidMount: function() {
		this.props.ajax({
			url: 'api/datasets',
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

	renderAddNew: function() {
		return (
			React.createElement("div", {className: "well"}, 
				React.createElement("h4", null, " Add new dataset "), 
				React.createElement("p", null, "Please select and upload all the files in your dataset"), 
				React.createElement(Files, {apiURL: "api/datasets", error: this.props.error, 
					cancel: setState.bind(this, {addNewPaneOpen:false})})
			)
		);
	},

	humanSize: function(sz) {
		if (sz < 1024) {
			return [sz,"B"];
		} else if (sz < 1024 * 1024) {
			return [(sz/1024).toFixed(1), "KiB"];
		} else if (sz < 1024 * 1024 * 1024) {
			return [(sz/(1024*1024)).toFixed(1), "MiB"];
		} else if (sz < 1024 * 1024 * 1024 * 1024) {
			return [(sz/(1024*1024*1024)).toFixed(1), "GiB"];
		} else {
			return [(sz/(1024*1024*1024*1024)).toFixed(1), "TiB"];
		}
	},

	renderDataset: function(dataset) {
		var sz = this.humanSize(dataset.size);
		return (
			React.createElement("tr", {key: dataset.id}, 
				React.createElement("td", null, dataset.id), 
				React.createElement("td", null, dataset.name), 
				React.createElement("td", {style: {textAlign:'right'}}, sz[0]), 
				React.createElement("td", {style: {textAlign:'left'}}, sz[1]), 
				React.createElement("td", {style: {textAlign:'right'}}, new Date(dataset.date).toLocaleString())
			)
		);
	},

	render: function() {
		return (
			React.createElement("div", {className: "dataset-page"}, 
				React.createElement("h3", null, " Datasets "), 
				 this.state.addNewPaneOpen ?
					this.renderAddNew() :
					React.createElement("div", {className: "row"}, 
						React.createElement("div", {className: "col-md-2 col-md-offset-10"}, 
							React.createElement("button", {type: "button", className: "btn btn-default", 
								onClick: setState.bind(this, {addNewPaneOpen:true})}, " Add new dataset ")
						)
					), 
				
				React.createElement("table", {className: "table table-condensed table-hover"}, 
					React.createElement("thead", null, 
						React.createElement("tr", null, 
							React.createElement("th", null, "Id"), 
							React.createElement("th", null, "Name"), 
							React.createElement("th", {style: {textAlign:'right'}}, "Size"), 
							React.createElement("th", {style: {textAlign:'left'}}), 
							React.createElement("th", {style: {textAlign:'right'}}, "Date")
						)
					), 
					React.createElement("tbody", null, 
						 this.state.datasets.map(this.renderDataset) 
					)
				)
			)
		);
	}
});

var Workflows = React.createClass({displayName: "Workflows",
	render: function() {
		return (
			React.createElement("div", null, 
				React.createElement("h3", null, " Workflows ")
			)
		);
	},
});

var Jobs = React.createClass({displayName: "Jobs",
	render: function() {
		return (
			React.createElement("div", null, 
				React.createElement("h3", null, " Jobs ")
			)
		);
	},
});


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
