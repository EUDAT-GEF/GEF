/**
 * Created by wqiu on 17/08/16.
 */
import React, {PropTypes} from 'react';
import axios from 'axios';
import bows from 'bows';
import {Row, Col, Grid} from 'react-bootstrap';
import {Field, reduxForm} from 'redux-form';
// this is a detailed view of a service, user will be able to execute service in this view


const log = bows("Service");

const tagValueRow  = (tag, value) => (
    <Row>
        <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>{tag}</Col>
        <Col xs={12} sm={9} md={9} >{value}</Col>
    </Row>
);

const JobCreatorForm = (props) => {
    const { handleSubmit, pristine, reset, submitting, service } = props;
    const inputStyle = {
        height: '34px',
        padding: '6px 12px',
        fontSize: '14px',
        lineHeight: '1.42857143',
        color: '#555',
        backgroundColor: '#fff',
        backgroundImage: 'none',
        border: '1px solid #ccc',
        borderRadius: '4px',
    }
    return (
        <form onSubmit={handleSubmit}>
            <Row>
                <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>
                    PID or URL
                </Col>
                <Col xs={12} sm={9} md={9} >
                    <div className="input-group">
                        <Field name="pid" component="input" type="text" placeholder="Put your PID or URL"
                               style={inputStyle} className="form-control"/>
                        <span className="input-group-btn">
                            <button type="submit" className="btn btn-default" onClick={handleSubmit} disabled={pristine || submitting}>
                                <span className="glyphicon glyphicon-play" aria-hidden="true"></span> Submit
                            </button>
                        </span>
                    </div>
                </Col>
            </Row>
        </form>
    )
};

const JobCreator = reduxForm({form: 'JobCreator'} )(JobCreatorForm);

class Service extends React.Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.props.handleSubmit.bind(this);
    }

    componentDidMount() {
        this.props.fetchService(this.props.service.ID);
    }

    render() {
        if(! this.props.selectedService.Service) {
            return (<div>loading</div>)
        } else {
            log("selectedService:", this.props.selectedService);
            const {ID, Name, Description, Version} = this.props.selectedService.Service;
            return (

                <div className="panel panel-default">
                    <div className="panel-body">
                        <div style={{margin: "1em"}}>
                            <div style={{height: "1em"}}></div>
                            {tagValueRow("Name", Name)}
                            {tagValueRow("ID", ID)}
                            {tagValueRow("Description", Description)}
                            {tagValueRow("Version", Version)}
                            <JobCreator handleSubmit={this.handleSubmit} service={this.props.selectedService.Service}/>
                            <div style={{height: "1em"}}></div>
                        </div>
                    </div>
                </div>

            )
        }
    }
}


Service.propTypes = {
    service: PropTypes.object.isRequired,
    fetchService: PropTypes.func.isRequired,
    selectedService: PropTypes.object.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    volumes: PropTypes.array.isRequired,
};

export default Service;