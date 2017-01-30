import React, {PropTypes} from 'react';
import bows from 'bows';
import {Row, Col} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'
import Job from './Job'

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
        this.props.fetchServices();
    }

    render() {
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
    }
}

Jobs.propTypes = {
    jobs: PropTypes.array.isRequired,
    fetchJobs: PropTypes.func.isRequired,
    services: PropTypes.array.isRequired,
    fetchServices: PropTypes.func.isRequired,
    jobID: PropTypes.string
};

export default Jobs;
