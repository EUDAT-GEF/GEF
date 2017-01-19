import React, {PropTypes} from 'react';
import { Row, Col, Grid } from 'react-bootstrap';
import { toPairs } from '../utils/utils';


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


const Job = ({job, title}) => {
    return (
        <div style={{border: "1px solid black"}}>
            <h4> Selected job</h4>
            <JobRow tag="ID" value={job.ID}/>
            <JobRow tag="Name" value={title}/>
            <JobRow tag="Service ID" value={job.Service.ID}/>
            <JobRow tag="Description" value={job.Service.Description}/>
            <JobRow tag="Version" value={job.Service.Version}/>
            <JobRow style={{marginTop:'1em'}} tag="Status" value={job.State.Status}/>
        </div>
    )

};

Job.propTypes = {
    job: PropTypes.object.isRequired
};

export default Job
