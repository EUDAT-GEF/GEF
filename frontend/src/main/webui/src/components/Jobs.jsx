'use strict';

import React, {PropTypes} from 'react';
import bows from 'bows';
import _ from 'lodash';
import {Row, Col} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'
import Job from './Job'

const log = bows('Jobs');

const JobRow = ({job}) => (
    <LinkContainer to={`/jobs/${job.ID}`}>
        <Row>
            <Col xs={12} sm={4} md={4}>{job.ID}</Col>
            <Col xs={12} sm={4} md={4}><i className="glyphicon glyphicon-transfer"/>{job.Service.Name}</Col>
            <Col xs={12} sm={4} md={4}>{job.State.Status}</Col>
        </Row>
    </LinkContainer>
);

const Header = () => (
    <div className="row table-head">
        <div className="col-xs-12 col-sm-4">Job ID</div>
        <div className="col-xs-12 col-sm-4">Service Name</div>
        <div className="col-xs-12 col-sm-4">Status</div>
    </div>
);

class Jobs extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.props.fetchJobs();
    }

    render() {
        return (
            <div>
                <h3>Browse Jobs</h3>
                <h4>All jobs</h4>
                <Header/>
                {_.map(this.props.jobs, (job) => {
                    if(job.ID === this.props.params.id) {
                        return <Job key={job.ID} job={job}/>
                    } else{
                        return <JobRow key={job.ID} job={job}/>
                    }
                })}
            </div>
        );
    }

}

Jobs.propTypes = {
    jobs: PropTypes.array.isRequired,
    fetchJobs: PropTypes.func.isRequired,
    jobID: PropTypes.string
};

export default Jobs;
