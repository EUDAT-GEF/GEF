import React, {PropTypes} from 'react';
import bows from 'bows';
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
        <Row style={{marginTop:'0.5em', marginBottom:'0.5em'}}>
            <Col xs={12} sm={3} md={3} style={style}>{image.Name || image.RepoTag}</Col>
            <Col xs={12} sm={9} md={9} style={style}>{image.Description}</Col>
        </Row>
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
        this.handleSubmit = this.props.handleSubmit.bind(this);
    }

    componentDidMount() {
        this.props.fetchServices();
    }

    render() {
        console.log("The id of selected service is:", this.props.params.id);

        return (
            <div>
                <h3>Browse Services</h3>
                <h4>All Services</h4>
                <Header/>
                { this.props.services.map((service) => {
                    if (service.ID === this.props.params.id) {
                        return <Service key={service.ID} service={service} fetchService={this.fetchService} selectedService={this.props.selectedService} handleSubmit={this.handleSubmit} volumes={this.props.volumes}/>;
                    } else {
                        const isGef = service.Input.length > 0 && service.Output.length > 0;
                        return isGef ? <ServiceRow key={service.ID} service={service} /> : false;
                    }
                })}
                { this.props.services.map((service) => {
                    const isGef = service.Input.length > 0 && service.Output.length > 0;
                    return isGef ? false : <ImageRow key={service.ID} image={service} />;
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
