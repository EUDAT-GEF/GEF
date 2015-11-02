/** @jsx React.DOM */
(function() {
"use strict";

// 1. create gefservice: first screen
//	- make possible the upload of a Dockerfile with some files
//	- 	the Docker file must contain all the labels
// 	- 	the server reads the labels and displays them in UI
//	- user must accept to create the service
//	- the frontend server delegates gef-docker to build the service(docker image)
//	- 	and the final image becomes a gef service
//	- 	the gefservice is assigned a PID
//	-	the user is informed, gets the PID of the new service
// 2. list all the gefservices with their metadata
// 	- make possible to execute one of them -> switch to the run wizard
// 3. execute gefservice
//	- input a pid of a dataset
//	- select one of the gefservices from a list
//	- run -> switch to the job monitoring page
// 4. job monitoring
//	- select running/finished job
//	- the UI displays the status, stdout and stderr
//	- the server exports the results automatically to b2drop
// 5. gc for jobs older than...
//


var VERSION = "0.3.7";
var PT = React.PropTypes;
var ErrorPane = window.MyReact.ErrorPane;
var FileAddButton = window.MyReact.FileAddButton;
var Files = window.MyReact.Files;

window.MyGEF = window.MyGEF || {};

var apiNames = {
	datasets: "/gef/api/datasets",
	builds:   "/gef/api/builds",
	services: "/gef/api/images",
};

function setState(state) {
	var t = this;
	if (t && t != window && t.setState) {
		t.setState(state);
	}
}

var Main = React.createClass({displayName: "Main",
	getInitialState: function () {
		return {
			page: this.executeService,
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

	buildService: function() {
		return (
			React.createElement(BuildService, {error: this.error, ajax: this.ajax})
		);
	},

	executeService: function() {
		return (
			React.createElement(ExecuteService, {error: this.error, ajax: this.ajax})
		);
	},

	browseJobs: function() {
		return (
			React.createElement(BrowseJobs, {error: this.error, ajax: this.ajax})
		);
	},

	browseDatasets: function() {
		return (
			React.createElement(BrowseDatasets, {error: this.error, ajax: this.ajax})
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
								this.renderToolListItem(this.buildService, "Build Service"), 
								this.renderToolListItem(this.executeService, "Execute Service"), 
								this.renderToolListItem(this.browseJobs, "Browse Jobs")
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


var BuildService = React.createClass({displayName: "BuildService",
	props: {
		error: PT.func.isRequired,
		ajax: PT.func.isRequired,
	},

	getInitialState: function() {
		return {
			buildURL: null,
			error: null,
			created: null,
		};
	},

	getURL: function() {
		return this.state.buildURL;
	},

	componentDidMount: function() {
		this.props.ajax({
			type: "POST",
			url: apiNames.builds,
			success: function(json, textStatus, jqXHR) {
				if (!this.isMounted()) {
					return;
				}
				if (!json.Location) {
					this.props.error("Didn't get json location from server");
					return;
				}
				var buildURL = apiNames.builds + "/" + json.Location;
				this.setState({buildURL: buildURL});
				console.log("create new service url :", buildURL);
			}.bind(this),
		});
	},

	dockerfileAdd: function(files) {
		if (files.length === 1) {
			// console.log("dockerfile add", files);
			this.setState({dockerfile: files[0]});
		}
	},

	error: function(err) {
		this.setState({error:err});
	},

	onCreated: function(responseText) {
		var json = JSON.parse(responseText);
		this.setState({created:json.Service});
	},

	renderCreated: function() {
		return (
			React.createElement("div", null, 
				React.createElement("p", null, "Created gef service"), 
				React.createElement(InspectService, {service: this.state.created})
			)
		);
	},

	renderFiles: function() {
		return (
			React.createElement("div", null, 
				React.createElement("p", null, "Please select and upload the Dockerfile, together" + ' ' +
				"with other files which are part of the container"), 
				React.createElement(Files, {getApiURL: this.getURL, cancel: function(){}, error: this.error, done: this.onCreated}), 
				this.state.error ? React.createElement("p", {style: {color:'red'}}, this.state.error) : false
			)
		);
	},

	render: function() {
		return (
			React.createElement("div", null, 
				React.createElement("h3", null, " Build Service "), 
				this.state.created ? this.renderCreated() : this.renderFiles()
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
			services: [],
			selected: null,
		};
	},

	componentDidMount: function() {
		this.props.ajax({
			url: apiNames.services,
			success: function(json, textStatus, jqXHR) {
				if (!this.isMounted()) {
					return;
				}
				if (!json.Services) {
					this.props.error("Didn't get Services from server");
					return;
				}
				this.setState({services: json.Services});
			}.bind(this),
		});
	},

	showService: function(serviceId) {
		this.props.ajax({
			url: apiNames.services+"/"+serviceId,
			success: function(json, textStatus, jqXHR) {
				if (!this.isMounted()) {
					return;
				}
				if (!json.Service) {
					this.props.error("Didn't get Service json from server");
					return;
				}
				this.setState({selected: json});
			}.bind(this),
		});
	},

	renderHeads: function(dataset) {
		return (
			React.createElement("div", {className: "row table-head"}, 
				React.createElement("div", {className: "col-xs-12 col-sm-4"}, "Name"), 
				React.createElement("div", {className: "col-xs-12 col-sm-8"}, "ID")
			)
		);
	},

	renderService: function(service) {
		return (
			React.createElement("div", {className: "row", key: service.ID, onClick: this.showService.bind(this, service.ID)}, 
				React.createElement("div", {className: "col-xs-12 col-sm-4"}, React.createElement("i", {className: "glyphicon glyphicon-transfer"}), " ", service.Name), 
				React.createElement("div", {className: "col-xs-12 col-sm-8"}, service.ID)
			)
		);
	},

	render: function() {

		return (
			React.createElement("div", {className: "execute-service-page"}, 
				React.createElement("h3", null, " Execute Service "), 
				 this.state.selected ? React.createElement(InspectService, {service: this.state.selected.Service}) : false, 
				React.createElement("div", {style: {height:"1em"}}), 
				React.createElement("h4", null, "All services"), 
				 this.renderHeads(), 
				React.createElement("div", {className: "services-table"}, 
					 this.state.services.map(this.renderService) 
				)
			)
		);
	}
});

///////////////////////////////////////////////////////////////////////////////

var InspectService = React.createClass({displayName: "InspectService",
	props: {
		service: PT.object.isRequired,
	},

	renderIO: function(io) {
		return React.createElement("div", null, React.createElement("dt", null, io.Name), React.createElement("dd", null, io.Path))
	},

	renderValue: function(value) {
		if (typeof value === 'object') {
			return (React.createElement("dl", {className: "dl-horizontal"}, " ",  value.map(this.renderIO), " "));
		}else {
			return value;
		}
	},

	renderRow: function(tag, value) {
		return (
			React.createElement("div", {className: "row"}, 
				React.createElement("div", {className: "col-xs-12 col-sm-3", style: {fontWeight:700}}, tag), 
				React.createElement("div", {className: "col-xs-12 col-sm-9"}, this.renderValue(value))
			)
		);
	},

	render: function() {
		var service = this.props.service;
		return (
			React.createElement("div", {className: ""}, 
				React.createElement("div", {style: {height:"1em"}}), 
				React.createElement("h4", null, "Selected service"), 
				this.renderRow("Name", service.Name), 
				this.renderRow("ID", service.ID), 
				this.renderRow("Description", service.Description), 
				this.renderRow("Version", service.Version), 
				React.createElement("div", {style: {height:"1em"}}), 
				this.renderRow("Input", service.Input), 
				this.renderRow("Output", service.Output), 
				React.createElement("div", {className: "row"}, 
					React.createElement("div", {className: "col-xs-4"}), 
					React.createElement("button", {className: "btn btn-primary"}, "Execute")
				)
			)
		);
	},
});

///////////////////////////////////////////////////////////////////////////////

var BrowseJobs = React.createClass({displayName: "BrowseJobs",
	render: function() {
		return (
			React.createElement("div", null, 
				React.createElement("h3", null, " Browse Jobs ")
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
				// console.log(json.datasets);
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

	toggleExpand: function(coll) {
		coll.expand = !coll.expand;
		this.setState({datasets:this.state.datasets});
	},

	renderRow: function(indent, state, name, size, date, fn) {
		var indentStyle = {marginLeft: 20 * indent};
		var sz = humanSize(size);
		var icon = "glyphicon " + (state === 'close' ? "glyphicon-folder-close" :
			state === 'open' ? "glyphicon-folder-open" : "glyphicon-file");
		return (
			React.createElement("div", {className: "row", key: name+indent, onClick: fn}, 
				React.createElement("div", {className: "col-xs-12 col-sm-5 col-md-5"}, 
					React.createElement("div", {style: indentStyle}, 
						React.createElement("i", {className: icon}), " ", name
					)
				), 
				React.createElement("div", {className: "col-xs-12 col-sm-2 col-md-2", style: {textAlign:'right'}}, sz[0], " ", sz[1]), 
				React.createElement("div", {className: "col-xs-12 col-sm-5 col-md-5", style: {textAlign:'right'}}, new Date(date).toLocaleString())
			)
		);
	},

	renderColl: function(indent, coll) {
		return (
			React.createElement("div", null, 
				 this.renderRow(indent, coll.expand ? "open":"close",
					coll.name, coll.size, coll.date, this.toggleExpand.bind(this, coll)), 
				coll.expand ?
					React.createElement("div", null, 
						 dataset.entry.colls.map(this.renderColl.bind(this, indent+1)), 
						 dataset.entry.files.map(this.renderFile.bind(this, indent+1)) 
					)
				: false
			)
		);
	},

	renderFile: function(indent, file) {
		return this.renderRow(indent, "file", file.name, file.size, file.date, function(){});
	},

	renderDataset: function(dataset) {
		return (
			React.createElement("div", null, 
				 this.renderRow(0, dataset.expand ? "open":"close",
					dataset.id, dataset.entry.size, dataset.entry.date, this.toggleExpand.bind(this, dataset)), 
				dataset.expand ?
					React.createElement("div", null, 
						 dataset.entry.colls.map(this.renderColl.bind(this, 1)), 
						 dataset.entry.files.map(this.renderFile.bind(this, 1)) 
					)
				: false
			)
		);
	},

	render: function() {
		return (
			React.createElement("div", {className: "dataset-page"}, 
				React.createElement(CreateDataset, {error: this.props.error, ajax: this.props.ajax}), 
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


///////////////////////////////////////////////////////////////////////////////

function pairs(o) {
	var a = []
	for (var k in o) {
		if (o.hasOwnProperty(k)) {
			a.push([k, o[k]]);
		}
	}
	return a;
}


///////////////////////////////////////////////////////////////////////////////

window.MyGEF.main = React.render(React.createElement(Main, null),  document.getElementById('page'));
window.MyGEF.footer = React.render(React.createElement(Footer, null), document.getElementById('footer') );

})();
