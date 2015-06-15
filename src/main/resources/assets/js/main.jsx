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

var Main = React.createClass({
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
			<Datasets error={this.error} ajax={this.ajax} />
		);
	},

	workflows: function() {
		return (
			<Workflows error={this.error} ajax={this.ajax} />
		);
	},

	jobs: function() {
		return (
			<Jobs error={this.error} ajax={this.ajax} />
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
								{this.renderToolListItem(this.datasets, "Datasets")}
								{this.renderToolListItem(this.workflows, "Workflows")}
								{this.renderToolListItem(this.jobs, "Jobs")}
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

var Datasets = React.createClass({
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
			<div className="well">
				<h4> Add new dataset </h4>
				<p>Please select and upload all the files in your dataset</p>
				<Files apiURL="api/datasets" error={this.props.error}
					cancel={setState.bind(this, {addNewPaneOpen:false})} />
			</div>
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
			<tr key={dataset.id}>
				<td>{dataset.id}</td>
				<td>{dataset.name}</td>
				<td style={{textAlign:'right'}}>{sz[0]}</td>
				<td style={{textAlign:'left'}}>{sz[1]}</td>
				<td style={{textAlign:'right'}}>{new Date(dataset.date).toLocaleString()}</td>
			</tr>
		);
	},

	render: function() {
		return (
			<div className="dataset-page">
				<h3> Datasets </h3>
				{ this.state.addNewPaneOpen ?
					this.renderAddNew() :
					<div className="row">
						<div className="col-md-2 col-md-offset-10">
							<button type="button" className="btn btn-default"
								onClick={setState.bind(this, {addNewPaneOpen:true})}> Add new dataset </button>
						</div>
					</div>
				}
				<table className="table table-condensed table-hover">
					<thead>
						<tr>
							<th>Id</th>
							<th>Name</th>
							<th style={{textAlign:'right'}}>Size</th>
							<th style={{textAlign:'left'}}></th>
							<th style={{textAlign:'right'}}>Date</th>
						</tr>
					</thead>
					<tbody>
						{ this.state.datasets.map(this.renderDataset) }
					</tbody>
				</table>
			</div>
		);
	}
});

var Workflows = React.createClass({
	render: function() {
		return (
			<div>
				<h3> Workflows </h3>
			</div>
		);
	},
});

var Jobs = React.createClass({
	render: function() {
		return (
			<div>
				<h3> Jobs </h3>
			</div>
		);
	},
});


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
