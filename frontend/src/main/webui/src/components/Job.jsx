import React, {PropTypes} from 'react';
import { Row, Col, Grid, Table } from 'react-bootstrap';
import { toPairs } from '../utils/utils';
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import actions from '../actions/actions';
import FileTree from './FileTree'
import ConsoleOutput from './ConsoleOutput'

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

const JobRow = ({tag, value, style}) => (
    <Row style={style}>
        <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>{tag}</Col>
        <Col xs={12} sm={3} md={3} ><Value value={value}/></Col>
    </Row>
);

class Job extends React.Component {
    constructor(props) {
        super(props);
    }

    handleInspectInputVolume() {
        this.props.actions.inspectVolume(this.props.job.InputVolume)
    }

    handleInspectOutputVolume() {
        this.props.actions.inspectVolume(this.props.job.OutputVolume)
    }

    handleConsoleOutput() {
        this.props.actions.consoleOutputFetch(this.props.job.ID)
    }

    componentDidMount() {
        this.props.actions.inspectVolume(); // send an empty volumeID when a new box is drown
        this.props.actions.consoleOutputFetch();
    }

    render() {
        let job = this.props.job;
        let service = this.props.service;
        let title = this.props.title;
        return (
            <div style={{border: "1px solid black"}}>
                <h4> Selected job</h4>
                <JobRow tag="ID" value={job.ID}/>
                <JobRow tag="Name" value={title}/>
                <JobRow tag="Input" value={job.Input}/>
                <JobRow tag="Service ID" value={job.ServiceID}/>
                <JobRow tag="Service Description" value={service ? service.Description : false}/>
                <JobRow tag="Service Version" value={service ? service.Version : false}/>
                <JobRow style={{marginTop:'1em'}} tag="Status" value={job.State.Status}/>
                <JobRow style={{marginTop:'1em'}} tag="Error" value={job.State.Error ? job.State.Error : false}/>
                <Row style={{marginTop:'1em'}}>
                    <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>Input Volume</Col>
                    <Col xs={12} sm={9} md={9} >
                        <button type="submit" className="btn btn-default" onClick={this.handleInspectInputVolume.bind(this)}>Inspect</button>
                    </Col>
                </Row>
                <Row style={{marginTop:'1em'}}>
                    <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>Output Volume</Col>
                    <Col xs={12} sm={9} md={9} >
                        <button type="submit" className="btn btn-default" onClick={this.handleInspectOutputVolume.bind(this)}>Inspect</button>
                    </Col>
                </Row>
                <Row style={{marginTop:'1em'}}>
                    <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>Console output</Col>
                    <Col xs={12} sm={9} md={9} >
                        <button type="submit" className="btn btn-default" onClick={this.handleConsoleOutput.bind(this)}>Show</button>
                    </Col>
                </Row>
                <FileTree/>
                <ConsoleOutput/>
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