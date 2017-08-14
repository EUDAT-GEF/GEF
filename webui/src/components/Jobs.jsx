import React from 'react';
import PropTypes from 'prop-types';
import { Row, Col, Grid, Panel, Table, Button, Glyphicon, Modal, OverlayTrigger } from 'react-bootstrap';
import {BootstrapTable, TableHeaderColumn} from 'react-bootstrap-table';
import { toPairs } from '../utils/utils';
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import * as actions from '../actions/actions';
import FileTree from './FileTree'

const selectRowProp = {
    mode: 'checkbox'
};
const inProgressColor = {
    color: '#f45d00'
};
const errorColor = {
    color: '#ff0000'
};
const successColor = {
    color: '#337ab7'
};
const progressAnimation = <img src="/images/progress-animation.gif" />;

let allJobs = [];
let jobStatusUpdateTimer;
let activeJobs;
let inactiveJobs;
let failedJobs;

class Jobs extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            showModal: false,
            buttonPressed: -1,
            timerOn: true,
        };

        this.options = {
            defaultSortName: 'created', // default sort column name
            defaultSortOrder: 'desc',  // default sort order
            onClickGroupSelected: this.onClickGroupSelected,
        };

        jobStatusUpdateTimer = setInterval(this.tick.bind(this), 1000);
    }

    componentDidMount() {
        this.props.fetchJobs();
        this.props.fetchServices();
    }

    componentWillUnmount() {
        clearInterval(jobStatusUpdateTimer);
    }

    tick() {
        this.props.fetchJobs();
        if ((this.state.timerOn) && (!this.hasJobsRunning())) {
            clearInterval(jobStatusUpdateTimer);
            this.setState({timerOn: false});
        }

        if ((!this.state.timerOn) && (this.hasJobsRunning())) {
            jobStatusUpdateTimer = setInterval(this.tick.bind(this), 1000);
            this.setState({timerOn: true});
        }
    }

    formatJobDuration(durationTime) {
        var sec_num = parseInt(durationTime, 10);
        var hours   = Math.floor(sec_num / 3600);
        var minutes = Math.floor((sec_num - (hours * 3600)) / 60);
        var seconds = sec_num - (hours * 3600) - (minutes * 60);

        if (hours   < 10) {hours   = "0"+hours;}
        if (minutes < 10) {minutes = "0"+minutes;}
        if (seconds < 10) {seconds = "0"+seconds;}
        return hours+':'+minutes+':'+seconds;
    }


    hasJobsRunning() {
        var runningJobfound = false;
        this.props.jobs.map((job) => {

            if (job.State.Code == -1) {
                console.log(job.State);
                runningJobfound = true;

            }
        });
        return runningJobfound;
    }

    buttonFormatter(cell, row, removeJob) {



        return (
            <ButtonGroup>
                <Button bsSize="xsmall"><Glyphicon glyph="console"/></Button>
                <Button bsSize="xsmall"><Glyphicon glyph="arrow-down"/></Button>
                <Button bsSize="xsmall"><Glyphicon glyph="arrow-up"/></Button>
                <Button bsSize="xsmall" onClick={() => removeJob(row.id)}><Glyphicon glyph="trash"/></Button>
            </ButtonGroup>
        );
    }

    statusFormatter(cell, row) {
        var currentProgress;
        var messageColor;

        if (row.code < 0) {
            currentProgress = progressAnimation;
            messageColor = inProgressColor;
        } else if (row.code == 0) {
            messageColor = successColor;
        } else {
            messageColor = errorColor;
        }
        return (
            <div style={messageColor}>{currentProgress} {cell}</div>
        );
    }

    getSelectedRowKeys() {
        //Here is your answer
        console.log(this.refs.table.state.selectedRowKeys)
    }


    handleJobRemoval(jobID) {
        console.log("REMOVING");
        console.log(jobID);
        this.props.actions.removeJob(jobID);
    }



    onClickGroupSelected(v1,v2) {
        console.log(v1);
        console.log(v2);
    }





    render() {

        allJobs = [];
        activeJobs = 0;
        inactiveJobs = 0;
        failedJobs = 0;




        if (this.props.jobs) {
            return (
                <div>
                    <h3>Browse Jobs</h3>
                    {this.props.jobs.map((job) => {
                        let service = null;
                        for (var i = 0; i < this.props.services.length; ++i) {
                            if (job.ServiceID == this.props.services[i].ID) {
                                service = this.props.services[i];
                                break;
                            }
                        }
                        let serviceName = (service && service.Name && service.Name.length) ? service.Name :
                            (service && service.ID && service.ID.length) ? service.ID : "unknown service";
                        let title = "Job from " + serviceName;





                        let execDuration = "";
                        if (job.State.Code == -1) {
                            let currentDate = new Date();
                            execDuration = currentDate - Date.parse(job.Created);
                            activeJobs += 1;
                        } else {
                            execDuration = Date.parse(job.Finished) - Date.parse(job.Created);
                            if (job.State.Code == 0) {
                                inactiveJobs += 1;
                            } else {
                                failedJobs += 1;
                            }
                        }

                        let createdDate = new Date(job.Created);
                        let fmtCreatedDate = createdDate.toLocaleDateString('en-GB');
                        let fmtCreatedTime = createdDate.toLocaleTimeString('en-GB');

                        allJobs.push({"title": title, "id": job.ID, "created": fmtCreatedDate + " " + fmtCreatedTime, "duration": this.formatJobDuration(execDuration/1000), "status": job.State.Status, "code": job.State.Code});



                    })}
                    {/*<Job key={job.ID} job={job} service={service} title={title}/>*/}
                    <Panel>
                        <Col sm={8}>
                            Out of {this.props.jobs.length} jobs <span style={inProgressColor}>{activeJobs} are active</span>, <span style={successColor}>{inactiveJobs} are finished successfully</span>,  <span style={errorColor}>{failedJobs} failed</span>
                        </Col>
                        <Col sm={4}>
                            <Button onClick={this.getSelectedRowKeys.bind(this)} className="btn pull-right"><Glyphicon glyph="trash"/> Remove selected jobs</Button>
                        </Col>
                    </Panel>
                    <div>
                        <BootstrapTable data={allJobs} selectRow={selectRowProp} options={this.options} ref="table">
                            <TableHeaderColumn dataField='id' isKey dataSort>ID</TableHeaderColumn>
                            <TableHeaderColumn dataField='title' dataSort>Title</TableHeaderColumn>
                            <TableHeaderColumn dataField='created' dataSort>Created</TableHeaderColumn>
                            <TableHeaderColumn dataField='duration' dataSort>Duration</TableHeaderColumn>
                            <TableHeaderColumn dataField='status' dataSort dataFormat={this.statusFormatter}>Status</TableHeaderColumn>
                            <TableHeaderColumn dataField="button" dataFormat={this.buttonFormatter} formatExtraData={this.props.actions.removeJob}>Operations</TableHeaderColumn>
                        </BootstrapTable>
                    </div>

                </div>
            );
        } else {
            return (
                <div><h4>No jobs found</h4></div>
            )
        }
    }
}

function mapStateToProps(state) {
    return state
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators(actions, dispatch)
    }
}

Jobs.propTypes = {
    jobs: PropTypes.array, // can be null
    fetchJobs: PropTypes.func.isRequired,
    services: PropTypes.array, // can be null
    fetchServices: PropTypes.func.isRequired,
    jobID: PropTypes.string
};

export default connect(mapStateToProps, mapDispatchToProps)(Jobs);