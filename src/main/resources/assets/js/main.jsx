/** @jsx React.DOM */
(function() {
"use strict";

// 1. create gefservice: first screen
//	- make possible the upload of a Dockerfile with some files
//	- 	the Docker file must contain all the labels
// 	- 	the server reads the labels and displays them in UI
//	- user must accept to create the image
//	- the frontend server delegates gef-docker to build the image
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
	images:   "/gef/api/images",
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
			<BrowseDatasets error={this.error} ajax={this.ajax} />
		);
	},

	executeService: function() {
		return (
			<ExecuteService error={this.error} ajax={this.ajax} />
		);
	},

	runningJobs: function() {
		return (
			<RunningJobs error={this.error} ajax={this.ajax} />
		);
	},

	createDataset: function() {
		return (
			<CreateDataset error={this.error} ajax={this.ajax} />
		);
	},

	createService: function() {
		return (
			<CreateService error={this.error} ajax={this.ajax} />
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
				<div className="container">
					<div className="row">
						<div className="col-xs-12 col-sm-2 col-md-2">
							<div className="list-group">
								{this.renderToolListItem(this.createService, "Create Service")}
								{this.renderToolListItem(this.executeService, "Execute Service")}
								{this.renderToolListItem(this.runningJobs, "Browse Jobs")}
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


var CreateService = React.createClass({
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
			console.log("dockerfile add", files);
			this.setState({dockerfile: files[0]});
		}
	},

	error: function(err) {
		this.setState({error:err});
	},

	created: function(responseText) {
		var json = JSON.parse(responseText);
		console.log("created", json.Image);
		this.setState({created:json.Image});
	},

	renderCreated: function() {
		var image = this.state.created;
		return (
			<div>
				<p>Created gef service</p>
				<p>{image.ID}</p>
				<p>{image.Labels}</p>
			</div>
		);
	},

	renderFiles: function() {
		return (
			<div>
				<p>Please select and upload the Dockerfile, together
				with other files which are part of the container</p>
				<Files getApiURL={this.getURL} cancel={function(){}} error={this.error} done={this.created} />
				{this.state.error ? <p style={{color:'red'}}>{this.state.error}</p> : false}
			</div>
		);
	},

	render: function() {
		return (
			<div>
				<h3> Create Service </h3>
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
			imageIds: [],
		};
	},

	componentDidMount: function() {
		console.log("url: ", apiNames.images);
		this.props.ajax({
			url: apiNames.images,
			success: function(json, textStatus, jqXHR) {
				if (!this.isMounted()) {
					return;
				}
				if (!json.ImageIDs) {
					this.props.error("Didn't get ImageIDs from server");
					return;
				}
				this.setState({imageIds: json.ImageIDs});
				console.log("got image ids :", json.ImageIDs);
			}.bind(this),
		});
	},

	showImage: function(imageId) {
		this.props.ajax({
			url: apiNames.images+"/"+imageId,
			success: function(json, textStatus, jqXHR) {
				if (!this.isMounted()) {
					return;
				}
				if (!json.Image) {
					this.props.error("Didn't get Image from server");
					return;
				}
				// this.setState({imageIds: json.ImageIDs});
				console.log("got image data:", json);
			}.bind(this),
		});
	},

	renderImageDetails: function(image) {
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

	renderHeads: function(dataset) {
		return (
			<div className="row table-head">
				<div className="col-xs-12 col-sm-12 col-md-12">Image ID</div>
			</div>
		);
	},

	renderImageId: function(imageId) {
		return (
			<div className="row" key={imageId} onClick={this.showImage.bind(this, imageId)}>
				<div className="col-xs-12 col-sm-12 col-md-12"><i className="glyphicon glyphicon-transfer"/> {imageId}</div>
			</div>
		);
	},

	render: function() {
		return (
			<div className="execute-service-page">
				<h3> Execute Service </h3>
				{ this.renderHeads() }
				<div className="images-table">
					{ this.state.imageIds.map(this.renderImageId) }
				</div>
			</div>
		);
	}
});

///////////////////////////////////////////////////////////////////////////////

var RunningJobs = React.createClass({
	render: function() {
		return (
			<div>
				<h3> Running Jobs </h3>
			</div>
		);
	},
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
				console.log(json.datasets);
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

window.MyGEF.main = React.render(<Main />,  document.getElementById('page'));
window.MyGEF.footer = React.render(<Footer />, document.getElementById('footer') );

})();
