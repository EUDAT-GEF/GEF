'use strict';

import React, {PropTypes} from 'react';
import {Row, Col, Grid} from 'react-bootstrap';
import Radium from 'radium';
import _ from 'lodash';


const styles = {
    jobRowStyle: {
        fontWeight:700
    },
    jobStyle: {
        height: "1em"
    }
};

const jobRowStyle = {
    fontWeight:700
};

const Value = ({value}) => {
    if (typeof value === 'object') {
        _.toPairs(value).map(({k, v}) =>
            (
                 <div><dt>{k}</dt><dd>{v}</dd></div>
            ))
    } else {
        return value;
    }
};

const JobRow = ({tag, value}) => (
    <Row>
        <Col xs={12} sm={3} md={3} style={styles.jobRowStyle}>{tag}</Col>
        <Col xs={12} sm={3} md={3} ><Value value={value}/></Col>
    </Row>
);


const InspectJob = ({job}) => {
    return (
        <div>
            <div style={styles.jobStyle}></div>
            <h4> Selected job</h4>
            <JobRow tag="ID" value={job.ID}/>
            <JobRow tag="Name" value={job.Service.Name}/>
            <JobRow tag="Service ID" value={job.Service.ID}/>
            <JobRow tag="Description" value={job.Service.Description}/>
            <JobRow tag="Version" value={job.Service.Version}/>
            <div style={styles.jobStyle}></div>
            <JobRow tag="Status" value={job.Status}/>
        </div>
    )

};

InspectJob.propTypes = {
    job: PropTypes.object.isRequired
};

export default Radium(InspectJob)

