import React from 'react';
import PropTypes from 'prop-types';
import {Row, Col} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'
import Job from './Job'

const JobRow = ({job, title}) => (
    <LinkContainer to={`/jobs/${job.ID}`}>
        <Row>
            <Col xs={12} sm={3} md={3}>{title}</Col>
            <Col xs={12} sm={9} md={9}>{job.State.Status}</Col>
        </Row>
    </LinkContainer>
);

const Header = () => (
    <div className="row table-head">
        <div className="col-xs-12 col-sm-3">Job</div>
        <div className="col-xs-12 col-sm-9">Status</div>
    </div>
);

class Jobs extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.props.fetchJobs();
        this.props.fetchServices();
    }

    render() {
        if (this.props.jobs) {
            return (
                <div>
                    <h3>Browse Jobs</h3>
                    <h4>All jobs</h4>
                    <Header/>
                    { this.props.jobs.map((job) => {
                        let service = null;
                        for (var i = 0; i < this.props.services.length; ++i) {
                            if (job.ServiceID == this.props.services[i].ID) {
                                service = this.props.services[i];
                                break;
                            }
                        }
                        const serviceName = (service && service.Name && service.Name.length) ? service.Name :
                            (service && service.ID && service.ID.length) ? service.ID : "unknown service";
                        const title = "Job from " + serviceName;

                        if (job.ID === this.props.params.id) {
                            return <Job key={job.ID} job={job} service={service} title={title}/>
                        } else {
                            return <JobRow key={job.ID} job={job} title={title}/>
                        }
                    })}
                </div>
            );
        } else {
            return (
                <div><h4>No jobs found</h4></div>
            )
        }
    }
}

Jobs.propTypes = {
    jobs: PropTypes.array, // can be null
    fetchJobs: PropTypes.func.isRequired,
    services: PropTypes.array, // can be null
    fetchServices: PropTypes.func.isRequired,
    jobID: PropTypes.string
};

export default Jobs;