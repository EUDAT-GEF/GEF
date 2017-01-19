import React, {PropTypes} from 'react';
import bows from 'bows';
import _ from 'lodash';
import {Row, Col} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'
import Service from './Service'

const log = bows('Servcies');

const ServiceRow = ({service}) => (
    <LinkContainer to={`/services/${service.ID}`}>
        <Row>
            <Col xs={12} sm={4} md={4}>{service.Name}</Col>
            <Col xs={12} sm={4} md={4}>{service.Description}</Col>
        </Row>
    </LinkContainer>
);

const Header = () => (
    <div className="row table-head">
        <div className="col-xs-12 col-sm-4">Name</div>
        <div className="col-xs-12 col-sm-4">Description</div>
    </div>
);

class Services extends React.Component {
    constructor(props) {
        super(props);
        this.fetchService = this.props.fetchService.bind(this);
        this.handleSubmit = this.props.handleSubmit.bind(this);
    }

    componentDidMount() {
        this.props.fetchServices();
    }

    render() {
        log("The id of selected service is:", this.props.params.id);

        return (
            <div>
                <h3>Browse Services</h3>
                <h4>All Services</h4>
                <Header/>
                {_.map(this.props.services, (service) => {
                    if (service.ID === this.props.params.id) {
                        return <Service key={service.ID} service={service} fetchService={this.fetchService} selectedService={this.props.selectedService} handleSubmit={this.handleSubmit} volumes={this.props.volumes}/>;
                    } else {
                        if ((service.Input.length > 0) && (service.Output.length > 0)) { // Show only GEF services
                            return <ServiceRow key={service.ID} service={service} />;
                        }
                    }
                })}
            </div>
        );
    }

}

Services.propTypes = {
    fetchServices: PropTypes.func.isRequired,
    fetchService: PropTypes.func.isRequired,
    services: PropTypes.array.isRequired,
    selectedService: PropTypes.object.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    volumes: PropTypes.array.isRequired,
};

export default Services;
