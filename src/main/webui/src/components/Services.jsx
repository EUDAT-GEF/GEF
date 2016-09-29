'use strict';
import React, {PropTypes} from 'react';
import bows from 'bows';
import _ from 'lodash';
import {Row, Col} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'

const log = bows('BrowseServcies');

const ServiceRow = ({service}) => (
    <LinkContainer to="/services(/:id)">
        <Row>
            <Col xs={12} sm={4} md={4}><i className="glyphicon glyphicon-transfer"/>{service.Name}</Col>
            <Col xs={12} sm={4} md={4}>{service.ID}</Col>
        </Row>
    </LinkContainer>
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
        this.props.fetchServices();
    }

    render() {
        _.map(this.props.services, (service) => {
            log("service: ", service);
        });
        return (
            <div>
                <h3>Browse Services</h3>
                <h4>All Services</h4>
                <Header/>
                {_.map(this.props.services, (service) => (<ServiceRow service={service}/>))}
            </div>
        );
    }

}

BrowseServices.propTypes = {
    fetchServices: PropTypes.func.isRequired,
    services: PropTypes.array.isRequired
};

export default BrowseServices;

