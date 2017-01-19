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
        return <div>{value}</div>;
    }
};

const JobRow = ({tag, value}) => (
    <Row>
        <Col xs={12} sm={3} md={3} style={styles.jobRowStyle}>{tag}</Col>
        <Col xs={12} sm={3} md={3} ><Value value={value}/></Col>
    </Row>
);


const Job = ({job, title}) => {
    console.log(job);
    return (
        <div style={{border: "1px solid black"}}>
            <div style={styles.jobStyle}></div>
            <h4> Selected job</h4>
            <JobRow tag="ID" value={job.ID}/>
            <JobRow tag="Name" value={title}/>
            <JobRow tag="Service ID" value={job.Service.ID}/>
            <JobRow tag="Description" value={job.Service.Description}/>
            <JobRow tag="Version" value={job.Service.Version}/>
            <div style={styles.jobStyle}></div>
            <JobRow tag="Status" value={job.State.Status}/>
        </div>
    )

};

Job.propTypes = {
    job: PropTypes.object.isRequired
};

export default Radium(Job)
