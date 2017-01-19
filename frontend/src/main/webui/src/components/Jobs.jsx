import React, {PropTypes} from 'react';
import bows from 'bows';
import _ from 'lodash';
import {Row, Col} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'
import Job from './Job'

const log = bows('Jobs');

const JobRow = ({job, title}) => (
    <LinkContainer to={`/jobs/${job.ID}`}>
        <Row>
            <Col xs={12} sm={4} md={4}>{title}</Col>
            <Col xs={12} sm={4} md={4}>{job.State.Status}</Col>
        </Row>
    </LinkContainer>
);

const Header = () => (
    <div className="row table-head">
        <div className="col-xs-12 col-sm-4">Job</div>
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
                    let title = "Job from ";
                    let serviceName = job.Service.Name;
                    if (serviceName.length == 0) {
                        serviceName = "Unknown service";
                    }
                    title = title + serviceName;

                    if (job.ID === this.props.params.id) {
                        return <Job key={job.ID} job={job} title={title}/>
                    } else {
                        return <JobRow key={job.ID} job={job} title={title}/>
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
