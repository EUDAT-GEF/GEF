'use strict';
import React, {PropTypes} from 'react';
import bows from 'bows';
import _ from 'lodash';
import {Row, Col} from 'react-bootstrap';

const log = bows('BrowseJobs');

const ServiceRow = ({service}) => (
    <Row>
        <Col xs={12} sm={4} md={4}>{service.ID}</Col>
        <Col xs={12} sm={4} md={4}><i className="glyphicon glyphicon-transfer"/>{job.Service.Name}</Col>
        <Col xs={12} sm={4} md={4}>{job.State.Status}</Col>
    </Row>
);

const Header = () => (
    <div className="row table-head">
        <div className="col-xs-12 col-sm-4">Name</div>
        <div className="col-xs-12 col-sm-4">ID</div>
    </div>
);

class BrowseServices extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.props.fetchJobs();
    }

    render() {
        _.map(this.props.jobs, (job) => {
            log("job: ", job );
        });
        return (
            <div>
                <h3>Browse Jobs</h3>
                <h4>All jobs</h4>
                <Header/>
                {_.map(this.props.jobs, (job) => (<JobRow job={job}/>))}
            </div>
        );
    }

}

BrowseServices.propTypes = {

};

export default BrowseServices;

