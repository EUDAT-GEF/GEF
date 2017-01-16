/**
 * Created by wqiu on 17/08/16.
 */
import React, {PropTypes} from 'react';
import axios from 'axios';
import apiNames from '../utils/GefAPI';
import bows from 'bows';
import _ from 'lodash';
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

const volume = (volume) => {
    return (
        <option key={volume.ID} value={volume.ID}> {volume.ID} </option>
    )
};

const getOptions  = (volumes) => {
    log("volumes", volumes);
    return volumes.map(volume)
};

const IOMap = (volumes, {VolumeID, Name, Path}) => (
    <div>
        {/*<label>{ID}</label>*/}
        <div>{Name}</div>
        <div>Internal Path: {Path} </div>
        <div>
            <Field name={VolumeID} component="select">
                <option></option>
                {getOptions(volumes)}
            </Field>
        </div>
    </div>
);

const JobCreatorForm = (props) => {
    console.log(props);
    const { handleSubmit, pristine, reset, submitting, volumes, service } = props;
    console.log("PRISTINE = " + pristine);
    console.log("SUBMITTING = " + submitting);
    return (
    <form onSubmit={handleSubmit}>
        <div>
            <label>PID</label>
            <div>
                <Field name="pid" component="input" type="text" placeholder="PID"/>
            </div>
        </div>
        <div>
            <button type="submit" onClick={handleSubmit} disabled={pristine || submitting}>Submit</button>
        </div>
    </form>
    )
};

const JobCreator = reduxForm(
    {form: 'JobCreator'}
)(JobCreatorForm);

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
                <div style={{border: "1px solid black"}}>
                    <div style={{height: "1em"}}></div>
                    {tagValueRow("Name", Name)}
                    {tagValueRow("ID", ID)}
                    {tagValueRow("Description", Description)}
                    {tagValueRow("Version", Version)}
                    <JobCreator volumes={this.props.volumes} handleSubmit={this.handleSubmit} service={this.props.selectedService.Service}/>
                    <div style={{height: "1em"}}></div>
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
