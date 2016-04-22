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
// 5. gc for jobs older than a few days
//


import React from 'react/lib/ReactWithAddons';
import ReactDOM from 'react-dom';
const ReactCSSTransitionGroup = React.addons.CSSTransitionGroup;
import jQuery from 'jquery';
import {ErrorPane, FileAddButton, Files} from './components.jsx';

var VERSION = "0.4.0";
var PT = React.PropTypes;

window.MyGEF = window.MyGEF || {};


if (!window.requestAnimationFrame) {
    window.requestAnimationFrame = (fn) => {
        setTimeout(fn, 16);
    }
}


var apiNames = {
	datasets: "/gef/api/datasets",
	builds:   "/gef/api/builds",
	services: "/gef/api/images",
	jobs: "/gef/api/jobs",
};

function setState(state) {
	var t = this;
	if (t && t != window && t.setState) {
		t.setState(state);
	}
}


var Main = React.createClass({
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

		var errs = this.state.errorMessages.slice();
		errs.push(err);
		this.setState({errorMessages: errs});

		setTimeout(function() {
			var errs = this.state.errorMessages.slice();
			errs.shift();
			this.setState({errorMessages: errs});
		}.bind(this), 10000);
	},

	ajax: function(ajaxObject) {
		if (!ajaxObject.error) {
			ajaxObject.error = function(jqXHR, textStatus, error) {
				if (jqXHR.readyState === 0) {
					this.error("Network error, please check your internet connection");
				} else if (jqXHR.responseText) {
					this.error(jqXHR.responseText + " ("+error+")");
				} else  {
					this.error(error + " ("+textStatus+")");
				}
				console.log("ajax error, jqXHR: ", jqXHR);
			}.bind(this);
		}
		// console.log("ajax", ajaxObject);
		jQuery.ajax(ajaxObject);
	},

	buildService: function() {
		return (
			<BuildService error={this.error} ajax={this.ajax} />
		);
	},

	executeService: function() {
		return (
			<ExecuteService error={this.error} ajax={this.ajax} />
		);
	},

	browseJobs: function() {
		return (
			<BrowseJobs error={this.error} ajax={this.ajax} />
		);
	},

	browseDatasets: function() {
		return (
			<BrowseDatasets error={this.error} ajax={this.ajax} />
		);
	},


	renderToolListItem: function(pageFn, title) {
		var klass = "list-group-item " + (pageFn === this.state.page ? "active":"");
		return (
			<a href="#" className={klass} onClick={setState.bind(this, {page:pageFn})}>
				{title}
			</a>
		);
	},

	render: function() {
		return	(
			<div>
				<ErrorPane errorMessages={this.state.errorMessages} />
				<div className="container-fluid">
					<div className="row">
						<div className="col-xs-12 col-sm-2 col-md-2">
							<div className="list-group">
								{this.renderToolListItem(this.buildService, "Build Service")}
								{this.renderToolListItem(this.executeService, "Execute Service")}
								{this.renderToolListItem(this.browseJobs, "Browse Jobs")}
							</div>
							<div className="list-group">
								{this.renderToolListItem(this.browseDatasets, "Browse Datasets")}
							</div>
						</div>
						<div className="col-xs-12 col-sm-10 col-md-10">
							{ this.state.page ? this.state.page() : false }
						</div>
					</div>
				</div>
			</div>
		);
	}
});


///////////////////////////////////////////////////////////////////////////////


var BuildService = React.createClass({
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
			<div>
				<p>Created gef service</p>
				<InspectService service={this.state.created} execute={this.props.ajax}/>
			</div>
		);
	},

	renderFiles: function() {
		return (
			<div>
				<p>Please select and upload the Dockerfile, together
				with other files which are part of the container</p>
				<Files getApiURL={this.getURL} cancel={function(){}} error={this.error} done={this.onCreated} />
				{this.state.error ? <p style={{color:'red'}}>{this.state.error}</p> : false}
			</div>
		);
	},

	render: function() {
		return (
			<div>
				<h3> Build Service </h3>
				{this.state.created ? this.renderCreated() : this.renderFiles() }
 			</div>
		);
	},
});

///////////////////////////////////////////////////////////////////////////////

var ExecuteService = React.createClass({
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

	execute: function(service) {
		this.props.ajax({
			type: "POST",
			url: apiNames.jobs,
			data: { imageID: service.ID },
			success: function(json, textStatus, jqXHR) {
				if (!this.isMounted()) {
					return;
				}
				if (!json.Location) {
					this.props.error("Didn't get json location from server");
					return;
				}
				// window.location.assign(window.location.origin+"/"+ json.Location);
				// var buildURL = apiNames.builds + "/" + json.Location;
				// this.setState({buildURL: buildURL});
				// console.log("create new service url :", buildURL);
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
			<div className="row table-head">
				<div className="col-xs-12 col-sm-4">Name</div>
				<div className="col-xs-12 col-sm-8">ID</div>
			</div>
		);
	},

	renderService: function(service) {
		return (
			<div className="row" key={service.ID} onClick={this.showService.bind(this, service.ID)}>
				<div className="col-xs-12 col-sm-4"><i className="glyphicon glyphicon-transfer"/> {service.Name}</div>
				<div className="col-xs-12 col-sm-8">{service.ID}</div>
			</div>
		);
	},

	render: function() {

		return (
			<div className="execute-service-page">
				<h3> Execute Service </h3>
				{ this.state.selected ? <InspectService service={this.state.selected.Service} execute={this.execute} /> : false}
				<div style={{height:"1em"}}></div>
				<h4>All services</h4>
				{ this.renderHeads() }
				<div className="services-table">
					{ this.state.services.map(this.renderService) }
				</div>
			</div>
		);
	}
});

///////////////////////////////////////////////////////////////////////////////
var InspectService = React.createClass({
	props: {
		service: PT.object.isRequired,
		execute: PT.func.isRequired,
	},

	renderRow: function(tag, value) {
		return (
			<div className="row">
				<div className="col-xs-12 col-sm-3" style={{fontWeight:700}}>{tag}</div>
				<div className="col-xs-12 col-sm-9">{value}</div>
			</div>
		);
	},

	renderIO: function(isInput, io) {
		return (
			<div className="row">
				<div className="col-xs-12 col-sm-3" style={{fontWeight:500}}>{io.Name}</div>
				<div className="col-xs-12 col-sm-5">
					{ isInput
						? <input type="text" style={{width:'100%'}} value="/tempZone/home/rods/GEF/datasets/set1"/>
						: <input type="checkbox" checked="checked"/> }
				</div>
				<div className="col-xs-12 col-sm-4">{io.Path}</div>
			</div>
		);
	},

	renderInput: function(io) {
		return (
			<div>
				<div className="row">
					<div className="col-xs-12 col-sm-3" style={{fontWeight:700}}>Input</div>
					<div className="col-xs-12 col-sm-5" style={{fontWeight:700}}>Map to B2SAFE PID</div>
					<div className="col-xs-12 col-sm-4" style={{fontWeight:700}}>Internal location</div>
				</div>
				{io.map(this.renderIO.bind(this, true))}
			</div>
		);
	},
	renderOutput: function(io) {
		return (
			<div>
				<div className="row">
					<div className="col-xs-12 col-sm-3" style={{fontWeight:700}}>Output</div>
					<div className="col-xs-12 col-sm-5" style={{fontWeight:700}}>Copy To B2DROP</div>
					<div className="col-xs-12 col-sm-4" style={{fontWeight:700}}>Internal location</div>
				</div>
				{io.map(this.renderIO.bind(this, false))}
			</div>
		);
	},

	render: function() {
		var service = this.props.service;
		console.log(service);
		return (
			<div className="">
				<div style={{height:"1em"}}></div>
				<h4>Selected service</h4>
				{this.renderRow("Name", service.Name)}
				{this.renderRow("ID", service.ID)}
				{this.renderRow("Description", service.Description)}
				{this.renderRow("Version", service.Version)}
				<div style={{height:"1em"}}></div>
				{this.renderInput(service.Input)}
				<div style={{height:"1em"}}></div>
				{this.renderOutput(service.Output)}
				<div style={{height:"1em"}}></div>
				<div className="row">
					<div className="col-xs-3"/>
					<button className="btn btn-primary" style={{width:300}}
						onClick={this.props.execute.bind(this.service)}>Execute</button>
				</div>
			</div>
		);
	},
});

///////////////////////////////////////////////////////////////////////////////

var BrowseJobs = React.createClass({
	props: {
		error: PT.func.isRequired,
		ajax: PT.func.isRequired,
	},

	getInitialState: function() {
		return {
			jobs: [],
			selected: null,
		};
	},

	componentDidMount: function() {
		this.props.ajax({
			url: apiNames.jobs,
			success: function(json, textStatus, jqXHR) {
				if (!this.isMounted()) {
					return;
				}
				if (!json.Jobs) {
					this.props.error("Didn't get Jobs from server");
					return;
				}
				this.setState({jobs: json.Jobs});
				console.log('jobs', json);
			}.bind(this),
		});
	},

	showJob: function(jobId) {
		this.props.ajax({
			url: apiNames.jobs+"/"+jobId,
			success: function(json, textStatus, jqXHR) {
				if (!this.isMounted()) {
					return;
				}
				if (!json.Job) {
					this.props.error("Didn't get Job json from server");
					return;
				}
				this.setState({selected: json});
				console.log('job', json);
			}.bind(this),
		});
	},

	renderHeads: function() {
		return (
			<div className="row table-head">
				<div className="col-xs-12 col-sm-4">Job ID</div>
				<div className="col-xs-12 col-sm-4">Service Name</div>
				<div className="col-xs-12 col-sm-4">Status</div>
			</div>
		);
	},

	renderJob: function(job) {
		var sOver = {overflow:'scroll'};
		return (
			<div className="row" key={job.ID} onClick={this.showJob.bind(this, job.ID)}>
				<div className="col-xs-12 col-sm-4" style={sOver}>{job.ID}</div>
				<div className="col-xs-12 col-sm-4" style={sOver}><i className="glyphicon glyphicon-transfer"/> {job.ServiceName}</div>
				<div className="col-xs-12 col-sm-4" style={sOver}>{job.Status}</div>
			</div>
		);
	},

	render: function() {

		return (
			<div className="list-jobs-page">
				<h3> Browse Jobs </h3>
				<div style={{height:"1em"}}></div>
				<h4>All jobs</h4>
				{ this.renderHeads() }
				<div className="jobs-table">
					{ this.state.jobs.map(this.renderJob) }
				</div>
			</div>
		);
	}
});

///////////////////////////////////////////////////////////////////////////////
var InspectJob = React.createClass({
	props: {
		job: PT.object.isRequired,
	},

	renderKV: function(k,v) {
		return <div><dt>{k}</dt><dd>{v}</dd></div>
	},

	renderValue: function(value) {
		if (typeof value === 'object') {
			var arr = [];
			for (var k in value) {
				if (value.hasOwnProperty(k)) {
					arr.push(this.renderKV(k, value[k]));
				}
			}
			return (<dl className="dl-horizontal"> { arr } </dl>);
		} else {
			return value;
		}
	},

	renderRow: function(tag, value) {
		return (
			<div className="row">
				<div className="col-xs-12 col-sm-3" style={{fontWeight:700}}>{tag}</div>
				<div className="col-xs-12 col-sm-9">{this.renderValue(value)}</div>
			</div>
		);
	},

	render: function() {
		var job = this.props.job;
		return (
			<div className="">
				<div style={{height:"1em"}}></div>
				<h4>Selected job</h4>
				{this.renderRow("ID", job.ID)}
				{this.renderRow("Name", job.Service.Name)}
				{this.renderRow("Service ID", job.Service.ID)}
				{this.renderRow("Description", job.Service.Description)}
				{this.renderRow("Version", job.Service.Version)}
				<div style={{height:"1em"}}></div>
				{this.renderRow("Status", job.Status)}
			</div>
		);
	},
});
///////////////////////////////////////////////////////////////////////////////

var BrowseDatasets = React.createClass({
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
			<div className="row table-head">
				<div className="col-xs-12 col-sm-5 col-md-5" >ID</div>
				<div className="col-xs-12 col-sm-2 col-md-2" style={{textAlign:'right'}}>Size</div>
				<div className="col-xs-12 col-sm-5 col-md-5" style={{textAlign:'right'}}>Date</div>
			</div>
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
			<div className="row" key={name+indent} onClick={fn}>
				<div className="col-xs-12 col-sm-5 col-md-5">
					<div style={indentStyle}>
						<i className={icon}/> {name}
					</div>
				</div>
				<div className="col-xs-12 col-sm-2 col-md-2" style={{textAlign:'right'}}>{sz[0]} {sz[1]}</div>
				<div className="col-xs-12 col-sm-5 col-md-5" style={{textAlign:'right'}}>{new Date(date).toLocaleString()}</div>
			</div>
		);
	},

	renderColl: function(indent, coll) {
		return (
			<div>
				{ this.renderRow(indent, coll.expand ? "open":"close",
					coll.name, coll.size, coll.date, this.toggleExpand.bind(this, coll)) }
				{coll.expand ?
					<div>
						{ dataset.entry.colls.map(this.renderColl.bind(this, indent+1)) }
						{ dataset.entry.files.map(this.renderFile.bind(this, indent+1)) }
					</div>
				: false}
			</div>
		);
	},

	renderFile: function(indent, file) {
		return this.renderRow(indent, "file", file.name, file.size, file.date, function(){});
	},

	renderDataset: function(dataset) {
		return (
			<div>
				{ this.renderRow(0, dataset.expand ? "open":"close",
					dataset.id, dataset.entry.size, dataset.entry.date, this.toggleExpand.bind(this, dataset)) }
				{dataset.expand ?
					<div>
						{ dataset.entry.colls.map(this.renderColl.bind(this, 1)) }
						{ dataset.entry.files.map(this.renderFile.bind(this, 1)) }
					</div>
				: false}
			</div>
		);
	},

	render: function() {
		return (
			<div className="dataset-page">
				<CreateDataset error={this.props.error} ajax={this.props.ajax} />
				<h3> Browse Datasets </h3>
				{ this.renderHeads() }
				<div className="dataset-table">
					{ this.state.datasets.map(this.renderDataset) }
				</div>
			</div>
		);
	}
});

///////////////////////////////////////////////////////////////////////////////

var CreateDataset = React.createClass({
	props: {
		error: PT.func.isRequired,
		ajax: PT.func.isRequired,
	},

	render: function() {
		return (
			<div>
				<h3> Create Dataset </h3>
				<p>Please select and upload all the files in your dataset</p>
				<Files apiURL={apiNames.datasets} error={this.props.error}
						cancel={function(){}} />
			</div>
		);
	},
});

///////////////////////////////////////////////////////////////////////////////

var Footer = React.createClass({
	about: function(e) {
		main.about();
		e.preventDefault();
		e.stopPropagation();
	},

	render: function() {
		return	(
			<div className="container">
				<div className="row">
					<div className="col-xs-12 col-sm-6 col-md-6">
						<p>	<img width="45" height="31" src="images/flag-ce.jpg" style={{float:'left', marginRight:10}}/>
							EUDAT receives funding from the European Union’s Horizon 2020 research
							and innovation programme under grant agreement No. 654065.&nbsp;
							<a href="#">Legal Notice</a>.
						</p>
					</div>
					<div className="col-xs-12 col-sm-6 col-md-6 text-right">
						<ul className="list-inline pull-right" style={{marginLeft:20}}>
							<li><span style={{color:'#173b93', fontWeight:'500'}}> GEF v.{VERSION}</span></li>
						</ul>
						<ul className="list-inline pull-right">
							<li><a target="_blank" href="http://eudat.eu/what-eudat">About EUDAT</a></li>
							<li><a href="https://github.com/GEFx">Go to GitHub</a></li>
							<li><a href="mailto:emanuel.dima@uni-tuebingen.de">Contact</a></li>
						</ul>
					</div>
				</div>
			</div>
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

window.MyGEF.main = ReactDOM.render(<Main />,  document.getElementById('page'));
window.MyGEF.footer = ReactDOM.render(<Footer />, document.getElementById('footer') );
