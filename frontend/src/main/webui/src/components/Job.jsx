import React, {PropTypes} from 'react';
import { Row, Col, Grid, Table, Button, Modal, OverlayTrigger } from 'react-bootstrap';
import { toPairs } from '../utils/utils';
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import actions from '../actions/actions';
import FileTree from './FileTree'

const Value = ({value}) => {
    if (typeof value === 'object') {
        toPairs(value).map(({k, v}) =>
            (
                <div><dt>{k}</dt><dd>{v}</dd></div>
            ))
    } else {
        return <div>{value}</div>;
    }
};

const JobRow = ({tag, value}) => (
    <Row>
        <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>{tag}</Col>
        <Col xs={12} sm={9} md={9} ><Value value={value}/></Col>
    </Row>
);

const JobStatusIndicator = ({code, tag, message}) => {
    const inProgressColor = {
        color: '#f45d00'
    }
    const errorColor = {
        color: '#ff0000'
    }
    const successColor = {
        color: '#337ab7'
    }
    const progressAnimation = <img src="/images/progress-animation.gif" />;
    let currentProgress;
    let messageColor;

    if (code < 0) {
        currentProgress = progressAnimation;
        messageColor = inProgressColor;
    } else if (code == 0) {
        messageColor = successColor;
    } else {
        messageColor = errorColor;
    }
    return (
        <Row>
            <Col xs={12} sm={3} md={3} style={{fontWeight: 700}}>{tag}</Col>
            <Col xs={12} sm={9} md={9} style={messageColor}>{currentProgress} {message}</Col>
        </Row>
    )
};

let stateUpdateTimer;

class Job extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            showModal: false,
            buttonPressed: -1,
        };
    }

    handleModalClose() {
        this.setState({ showModal: false });
    }

    handleModalOpen() {
        this.setState({ showModal: true });
    }

    handleInspectInputVolume() {
        if (this.props.job.State.Code > -1) {
            this.setState({ buttonPressed: 1 });
            this.props.actions.inspectVolume(this.props.job.InputVolume);
            this.handleModalOpen();
        }
    }

    handleInspectOutputVolume() {
        if (this.props.job.State.Code > -1) {
            this.setState({ buttonPressed: 2 });
            this.props.actions.inspectVolume(this.props.job.OutputVolume);
            this.handleModalOpen();
        }
    }

    handleConsoleOutput() {
        this.setState({ buttonPressed: 0 });
        this.props.actions.consoleOutputFetch(this.props.job.ID);
        this.handleModalOpen();
    }

    handleJobRemoval() {
        //this.setState({ buttonPressed: 0 });
        this.props.actions.removeJob(this.props.job.ID);
        this.props.actions.fetchJobs();
    }

    tick() {
        this.props.actions.fetchJobs();
        this.props.actions.consoleOutputFetch(this.props.job.ID);
    }

    handleClick(e) {
        e.stopPropagation();
    }

    componentDidMount() {
        this.props.actions.inspectVolume(); // send an empty volumeID when a new box is drown
        this.props.actions.consoleOutputFetch();
        if (this.props.job.State.Code < 0) {
            stateUpdateTimer = setInterval(this.tick.bind(this), 1000);
        }
    }

    renderModalWindow(title, body) {
        return (
            <div>
                <Modal show={this.state.showModal} onHide={this.handleModalClose.bind(this)}>
                    <Modal.Header closeButton>
                        <Modal.Title>{title}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        {body}
                    </Modal.Body>
                    <Modal.Footer>
                        <Button onClick={this.handleModalClose.bind(this)}>Close</Button>
                    </Modal.Footer>
                </Modal>
            </div>
        )
    }

    render() {
        let job = this.props.job;
        let service = this.props.service;
        let title = this.props.title;
        let buttonClass = "btn btn-default disabled";

        if (!this.state.progressIndicator) {
            this.state.progressIndicator = " ";
        } else {
            if (this.state.progressIndicator.length>4) {
                this.state.progressIndicator = " ";
            } else {
                this.state.progressIndicator += ".";
            }
        }

        if (job.State.Code > -1) {
            clearInterval(stateUpdateTimer);
            buttonClass = "btn btn-default";
            this.state.progressIndicator = "";
        }

        let modalTitle= "";
        let modalBody= "";
        if ((this.props.task.ServiceExecution) && (this.state.buttonPressed == 0)) {
            modalTitle = "Service container console output";
            modalBody = <span><pre>{this.props.task.ServiceExecution.ConsoleOutput}</pre></span>
        }

        if ((this.props.selectedVolume.volumeContent) && (this.state.buttonPressed > 0)) {
            if (this.state.buttonPressed == 1) {
                modalTitle = "Input volume inspection";
            } else {
                modalTitle = "Output volume inspection";
            }

            modalBody = <FileTree/>
        }
        let errorMessage;
        if (job.State.Error) {
            errorMessage = <JobRow tag="Error" value={job.State.Error}/>
        }

        return (
            <div className="panel panel-default">
                <div className="panel-body">
                    <div style={{margin: "1em"}}>
                        <JobRow tag="ID" value={job.ID}/>
                        <JobRow tag="Name" value={title}/>
                        <JobRow tag="Input" value={job.Input}/>
                        <JobRow tag="Service ID" value={job.ServiceID}/>
                        <JobRow tag="Service Description" value={service ? service.Description : false}/>
                        <JobRow tag="Service Version" value={service ? service.Version : false}/>
                        <JobStatusIndicator code={job.State.Code} tag="Status" message={job.State.Status+this.state.progressIndicator}/>
                        {errorMessage}

                        <Row style={{marginTop:'2em', marginBottom:'1em'}}>
                            <Col xs={12} sm={2} md={2}></Col>
                            <Col xs={12} sm={8} md={8}>
                                <div className="text-center">
                                    <div className="btn-group" role="group" aria-label="toolbar">
                                        <button type="button" className="btn btn-default" onClick={this.handleConsoleOutput.bind(this)}>
                                            <span className="glyphicon glyphicon-console" aria-hidden="true"></span> Console Output
                                        </button>

                                        <button type="button" className={buttonClass} onClick={this.handleInspectInputVolume.bind(this)}>
                                            <span className="glyphicon glyphicon-arrow-down" aria-hidden="true"></span> Input Volume
                                        </button>

                                        <button type="button" className={buttonClass} onClick={this.handleInspectOutputVolume.bind(this)}>
                                            <span className="glyphicon glyphicon-arrow-up" aria-hidden="true"></span> Output Volume
                                        </button>
                                        <button type="button" className={buttonClass} onClick={this.handleInspectOutputVolume.bind(this)}>
                                            <span className="glyphicon glyphicon-pause" aria-hidden="true"></span> Pause the Job
                                        </button>
                                        <button type="button" className="btn btn-default" onClick={this.handleJobRemoval.bind(this)}>
                                            <span className="glyphicon glyphicon-trash" aria-hidden="true"></span> Remove the Job
                                        </button>
                                    </div>
                                </div>
                            </Col>
                            <Col xs={12} sm={2} md={2}></Col>
                        </Row>
                        {this.renderModalWindow(modalTitle, modalBody)}
                    </div>
                </div>
            </div>
        )
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

Job.propTypes = {
    job: PropTypes.object.isRequired
};

export default connect(mapStateToProps, mapDispatchToProps)(Job);