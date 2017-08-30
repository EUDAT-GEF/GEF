import React from 'react';
import PropTypes from 'prop-types';
import {Row, Col} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'
import Service from './Service'

const ServiceRow = ({service}) => (
    <LinkContainer to={`/services/${service.ID}`}>
        <Row style={{marginTop:'0.5em', marginBottom:'0.5em'}}>
            <Col xs={12} sm={3} md={3}>{service.Name}</Col>
            <Col xs={12} sm={9} md={9}>{service.Description}</Col>
        </Row>
    </LinkContainer>
);

const ImageRow = ({image}) => {
    const style={color:'#aaa'}
    return (
        <LinkContainer to={`/services/${image.ID}`}>
            <Row style={{marginTop:'0.5em', marginBottom:'0.5em'}}>
                <Col xs={12} sm={3} md={3} style={style}>{image.Name || image.RepoTag}</Col>
                <Col xs={12} sm={9} md={9} style={style}>{image.Description}</Col>
            </Row>
        </LinkContainer>
    )
};

const Header = () => (
    <div className="row table-head">
        <div className="col-xs-12 col-sm-3">Name</div>
        <div className="col-xs-12 col-sm-9">Description</div>
    </div>
);

class Services extends React.Component {
    constructor(props) {
        super(props);
        this.fetchService = this.props.fetchService.bind(this);
        this.removeService = this.props.removeService.bind(this);
        this.fetchServices = this.props.fetchServices.bind(this);
        this.handleSubmit = this.props.handleSubmit.bind(this);
        this.handleUpdate = this.props.handleUpdate.bind(this);
        this.handleAddIO = this.props.handleAddIO.bind(this);
        this.handleRemoveIO = this.props.handleRemoveIO.bind(this);
    }

    componentDidMount() {
        this.props.fetchServices();
    }

    render() {
        if (this.props.services) {
            return (
                <div>
                    <h3>Browse Services</h3>
                    <h4>All Services</h4>
                    <Header/>
                    { this.props.services.map((service) => {
                        if (service.ID === this.props.params.id) {
                            return <Service key={service.ID}
                                            service={service}
                                            services={this.props.services}
                                            fetchService={this.fetchService}
                                            removeService={this.removeService}
                                            fetchServices={this.fetchServices}
                                            selectedService={this.props.selectedService}
                                            handleSubmit={this.handleSubmit}
                                            handleUpdate={this.handleUpdate}
                                            handleAddIO={this.handleAddIO}
                                            handleRemoveIO={this.handleRemoveIO}
                                            volumes={this.props.volumes}/>;
                        } else {
                            const isGef = service.Input && service.Output;
                            return isGef ? <ServiceRow key={service.ID} service={service}/> : <ImageRow key={service.ID} image={service}/>;
                        }
                    })}
                    { this.props.services.map((service) => {
                        //const isGef = service.Input && service.Output;
                        //return isGef ? false : <ImageRow key={service.ID} image={service}/>;
                    })}
                </div>
            )
        }
        else {
            return (
                <div><h4>No services found</h4></div>
            )
        }

    }

}

Services.propTypes = {
    fetchServices: PropTypes.func.isRequired,
    fetchService: PropTypes.func.isRequired,
    removeService: PropTypes.func.isRequired,
    services: PropTypes.array, // can be null
    selectedService: PropTypes.object.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    handleUpdate: PropTypes.func.isRequired,
    handleAddIO: PropTypes.func.isRequired,
    handleRemoveIO: PropTypes.func.isRequired,
    volumes: PropTypes.array.isRequired,
};

export default Services;
