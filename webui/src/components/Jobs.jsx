import React from 'react';
import PropTypes from 'prop-types';
import { Row, Col, Grid, Panel, Table, Button, Glyphicon, Modal, OverlayTrigger } from 'react-bootstrap';
import {BootstrapTable, TableHeaderColumn} from 'react-bootstrap-table';
import { toPairs } from '../utils/utils';
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import * as actions from '../actions/actions';
import FileTree from './FileTree'

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
let jobStatusUpdateTimer;

function sortByDate(a, b) {
     var dateA = new Date(a), dateB = new Date(b);
     return dateA - dateB;
 }

class Jobs extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            timerOn: true,
            showModal: false,
        };
        //jobStatusUpdateTimer = setInterval(this.tick.bind(this), 1000);
    }

    componentDidMount() {
        this.props.fetchJobs();
        this.props.fetchServices();
        this.jobStatusUpdateTimer = setInterval(() => this.tick(), 1000);
    }

    componentWillUnmount() {
        clearInterval(jobStatusUpdateTimer);
    }

    tick() {
        if (this.state.timerOn) {
            this.props.fetchJobs();
        }

        if ((this.state.timerOn) && (!this.hasJobsRunning())) {
            clearInterval(this.jobStatusUpdateTimer);
            this.setState({timerOn: false});
        }

        if ((!this.state.timerOn) && (this.hasJobsRunning())) {
            //jobStatusUpdateTimer = setInterval(this.tick.bind(this), 1000);
            this.jobStatusUpdateTimer = setInterval(() => this.tick(), 1000);
            this.setState({timerOn: true});
        }
    }



    hasJobsRunning() {
        var runningJobfound = false;
        if (this.props.jobs) {
            this.props.jobs.map((job) => {
                if (job.State.Code < 0) {
                    runningJobfound = true;
                }
            });
        }
        return runningJobfound;
    }

    jobDurationFormatter(durationTime) {
        var sec_num = parseInt(durationTime, 10);
        var hours   = Math.floor(sec_num / 3600);
        var minutes = Math.floor((sec_num - (hours * 3600)) / 60);
        var seconds = sec_num - (hours * 3600) - (minutes * 60);

        if (hours   < 10) {hours   = "0"+hours;}
        if (minutes < 10) {minutes = "0"+minutes;}
        if (seconds < 10) {seconds = "0"+seconds;}
        return hours+':'+minutes+':'+seconds;
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

    createdFormatter(cell, row) {
        let createdDate = new Date(row.created);
        let fmtCreatedDate = createdDate.toLocaleDateString('en-GB');
        let fmtCreatedTime = createdDate.toLocaleTimeString('en-GB');
        return (
            <div>{fmtCreatedDate + " " + fmtCreatedTime}</div>
        );
    }

    finishedFormatter(cell, row) {
        let finishedDate = new Date(row.finished);
        let finishedDateString = "";

        let fmtFinishedDate = finishedDate.toLocaleDateString('en-GB');
        let fmtFinishedTime = finishedDate.toLocaleTimeString('en-GB');

        if (row.code < 0) {
            finishedDateString = "running";
        } else {
            finishedDateString = fmtFinishedDate + " " + fmtFinishedTime;
        }
        return (
            <div>{finishedDateString}</div>
        );
    }

    removeSelectedJobs() {
        this.props.actions.removeJobs(this.refs.table.state.selectedRowKeys);
    }

    handleInspectVolume(volumeID) {
        this.setState({ buttonPressed: 1 });
        this.props.actions.inspectVolume(volumeID);
        this.handleModalOpen();
    }

    handleModalClose() {
        this.setState({ showModal: false });
    }

    handleModalOpen() {
        this.setState({ showModal: true });
    }

    expandComponent(row) {
        let inCounter = -1;
        let outCounter = -1;
        return (
            <div>
                <div>
                    <span><pre>{(row && row.console && row.console.length) ? row.console : "No information"}</pre></span>
                </div>
                <div className="text-center">
                    <div className="btn-group" role="group" aria-label="toolbar">
                        { row.input.map((input) => {
                            inCounter++;
                            return (
                                <Button onClick={() => this.handleInspectVolume(input["VolumeID"])}><Glyphicon
                                    glyph="arrow-down"/> Input #{inCounter+1}</Button>
                            )})
                        }

                        { row.output.map((output) => {
                            outCounter++;
                            return (
                                <Button onClick={() => this.handleInspectVolume(output["VolumeID"])}><Glyphicon
                                    glyph="arrow-up"/> Output #{outCounter+1}</Button>
                            )})
                        }

                        <Button onClick={ () => this.props.actions.removeJobs([row.id])}><Glyphicon glyph="trash"/> Remove job</Button>

                    </div>
                </div>
            </div>
        );
    }

    populateTable() {
        var allJobs = [];
        var jobIdList = [];
        var activeJobs = 0;
        var inactiveJobs = 0;
        var failedJobs = 0;
        this.props.jobs.map((job) => {
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

            let jobStartTime = new Date(job.Created);
            let jobFinishTime = new Date(jobStartTime.getTime() + (1000 * job.Duration));

            if (job.State.Code < 0) {
                activeJobs += 1;
            } else {
                if (job.State.Code == 0) {
                    inactiveJobs += 1;
                } else {
                    failedJobs += 1;
                }
            }

            let ConsoleOutput = "";
            if (job.Tasks) {
                for (var t = 0; t < job.Tasks.length; ++t) {
                    if (job.Tasks[t].Name == "Service execution") {
                        ConsoleOutput = job.Tasks[t].ConsoleOutput;
                        break;
                    }
                }
            }

            let inputVolumes = job.InputVolume ? job.InputVolume : [];
            let outputVolumes = job.OutputVolume ? job.OutputVolume : [];

            if (jobIdList.indexOf(job.ID)==-1) {
                jobIdList.push(job.ID);
                allJobs.push(
                    {
                        "id": job.ID,
                        "title": title,
                        "created": job.Created,
                        "duration": job.Duration,
                        "finished": jobFinishTime,
                        "status": job.State.Status,
                        "code": job.State.Code,
                        "console": ConsoleOutput,
                        "input": inputVolumes,
                        "output": outputVolumes
                    }
                );
            }
        });

        return [allJobs, activeJobs, inactiveJobs, failedJobs];
    }


    render() {
        const options = {
            defaultSortName: 'created', // default sort column name
            defaultSortOrder: 'desc',  // default sort order
            expandBy: 'row'
        };
        const selectRow = {
            mode: 'checkbox',
            clickToSelect: true,
            clickToExpand: true
        };

        if (this.props.jobs) {
            var tableData = this.populateTable();
            var allJobs = tableData[0];
            var activeJobs = tableData[1];
            var inactiveJobs = tableData[2];
            var failedJobs = tableData[3];

            return (
                <div>
                    <h3>Browse Jobs</h3>
                    <Panel>
                        <Col sm={8}>
                            Out of {this.props.jobs.length} jobs <span style={inProgressColor}>{activeJobs} are active</span>, <span style={successColor}>{inactiveJobs} are finished successfully</span>,  <span style={errorColor}>{failedJobs} failed</span>
                        </Col>
                        <Col sm={4}>
                            <Button onClick={this.removeSelectedJobs.bind(this)} className="btn pull-right"><Glyphicon glyph="trash"/> Remove selected jobs</Button>
                        </Col>
                    </Panel>
                    <div>
                        <BootstrapTable data={allJobs} selectRow={selectRow}  expandComponent={this.expandComponent.bind(this)} expandableRow={() => {return true}} options={options} expandColumnOptions={{ expandColumnVisible: true }} ref="table" keyField='id'>
                            <TableHeaderColumn dataField='id' dataSort expandable={ true }>ID</TableHeaderColumn>
                            <TableHeaderColumn dataField='title' dataSort>Title</TableHeaderColumn>
                            <TableHeaderColumn dataField='created' dataSort dataFormat={this.createdFormatter}>Created</TableHeaderColumn>
                            <TableHeaderColumn dataField='finished' dataSort dataFormat={this.finishedFormatter}>Finished</TableHeaderColumn>
                            <TableHeaderColumn dataField='duration' dataSort dataFormat={this.jobDurationFormatter}>Duration</TableHeaderColumn>
                            <TableHeaderColumn dataField='status' dataSort dataFormat={this.statusFormatter}>Status</TableHeaderColumn>
                        </BootstrapTable>
                    </div>
                    <div>
                        <Modal show={this.state.showModal} onHide={this.handleModalClose.bind(this)}>
                            <Modal.Header closeButton>
                                <Modal.Title>Volume Inspection</Modal.Title>
                            </Modal.Header>
                            <Modal.Body>
                                <FileTree/>
                            </Modal.Body>
                            <Modal.Footer>
                                <Button onClick={this.handleModalClose.bind(this)}>Close</Button>
                            </Modal.Footer>
                        </Modal>
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
