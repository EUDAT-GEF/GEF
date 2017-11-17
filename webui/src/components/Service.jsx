/**
 * Created by wqiu on 17/08/16.
 */
import React from 'react';
import PropTypes from 'prop-types';
import { Row, Col, Table, Button, Modal, Panel, FormGroup, FormControl, ControlLabel, Glyphicon } from 'react-bootstrap';
import {Field, FieldArray, reduxForm, initialize} from 'redux-form';
// this is a detailed view of a service, user will be able to execute service in this view

const renderField = ({ input, label, type, meta: { touched, error } }) => (
    <div className="form-group has-error has-feedback">
        <div>
            <input {...input} placeholder={label} type={type}/>
            {touched && error && <span>{error}</span>}
        </div>
    </div>
)

const tagValueRow  = (tag, value) => (
    <Row>
        <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>{tag}</Col>
        <Col xs={12} sm={9} md={9} >{value}</Col>
    </Row>
);

const JobCreatorForm = (props) => {
    const { handleSubmit, handleServiceEditModalOpen, handleRemoveService, pristine, reset, submitting, service } = props;
    const urlStyle = {
        padding: '4px 2px',
        fontSize: '14px',
        lineHeight: '1.42857143',
        color: '#555',
        backgroundColor: '#fff',
        backgroundImage: 'none',
        border: '1px solid #ccc',
        borderRadius: '4px',
    }
    const stringStyle = {
        height: '70px',
        padding: '4px 2px',
        fontSize: '14px',
        lineHeight: '1.42857143',
        color: '#555',
        backgroundColor: '#fff',
        backgroundImage: 'none',
        border: '1px solid #ccc',
        borderRadius: '4px',
    }
    const toolbarStyle = {
        padding: '5px',
    }

    let srcList = [];
    if (service.Input) {
        srcList = service.Input;
    }
    let inputCounter=-1;
    return (
        <form onSubmit={handleSubmit}>
            <Row>
                <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>
                    PID or URL
                </Col>
                <Col xs={12} sm={9} md={9} >
                    <div className="input-group">

                        {srcList.map((inputSrc) => {
                            inputCounter++;
                            if (inputSrc.Type.toLowerCase()=="url") {
                                return (
                                    <Field name={`pid_${inputSrc.ID}`} component="input"
                                           placeholder={`Input source #${inputCounter + 1}: insert PID or URL`}
                                           style={urlStyle} className="form-control" key={`pid_${inputSrc.ID}`}/>
                                )
                            }
                            if (inputSrc.Type.toLowerCase()=="string") {
                                return (
                                    <Field name={`pid_${inputSrc.ID}`} component="textarea"
                                           placeholder={`Input source #${inputCounter + 1}: insert your content`}
                                           style={stringStyle} className="form-control" key={`pid_${inputSrc.ID}`}/>
                                )
                            }
                        })}

                        <div className="text-center">
                            <div className="btn-group" role="group" aria-label="toolbar" style={toolbarStyle}>
                                <Button onClick={handleSubmit} disabled={pristine || submitting}><Glyphicon glyph="play"/> Start a New Job</Button>
                                <Button onClick={handleServiceEditModalOpen}><Glyphicon glyph="edit"/> Edit Metadata</Button>
                                <Button onClick={handleRemoveService}><Glyphicon glyph="trash"/> Remove the Service</Button>
                            </div>
                        </div>
                    </div>
                </Col>
            </Row>
        </form>
    )
};

const JobCreator = reduxForm({form: 'JobCreator'} )(JobCreatorForm);

const InputTable = ({service, handleRemoveIO}) => {
    let inCounter = -1;
    let inputs = [];
    let srcList = [];
    if (service.Input) {
        srcList = service.Input
    }
    srcList.map((input) => {
        inputs.push(input);
    });
    const IOTableRow = ({input, index}) => {
        return (
            <tr>
                <td>{input.ID}</td>
                <td>{input.Name}</td>
                <td>{input.Path}</td>
                <td>{input.Type}</td>
                <td>{input.FileName}</td>
                <td>
                    <Button type="submit" bsStyle="primary" bsSize="xsmall" onClick={(evt) => handleRemoveIO(true, index, evt)}>
                        <span className="glyphicon glyphicon-remove" aria-hidden="true"></span> Remove
                    </Button>
                </td>
            </tr>
        )
    };

    return (
        <Table responsive>
            <thead>
            <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Path</th>
                <th>Type</th>
                <th>File</th>
                <th></th>
            </tr>
            </thead>
            <tbody>
            { srcList.map((input) => {
                inCounter++;
                return (
                    <FieldArray name={`${input}.ID`} component={IOTableRow} input={input} key={`${input}.ID` + inCounter} index={inCounter}/>
                )
            })}
            </tbody>
        </Table>
    )
};

const OutputTable = ({service, handleRemoveIO}) => {
    let outCounter = -1;
    let outputs = [];
    let srcList = [];
    if (service.Output) {
        srcList = service.Output
    }
    srcList.map((out) => {
        outputs.push(out);
    });
    const IOTableRow = ({out, index}) => {
        return (
            <tr>
                <td>{out.ID}</td>
                <td>{out.Name}</td>
                <td>{out.Path}</td>
                <td>
                    <Button type="submit" bsStyle="primary" bsSize="xsmall" onClick={(evt) => handleRemoveIO(false, index, evt)}>
                        <span className="glyphicon glyphicon-remove" aria-hidden="true"></span> Remove
                    </Button>
                </td>
            </tr>
        )
    };

    return (
        <Table responsive>
            <thead>
            <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Path</th>
                <th></th>
            </tr>
            </thead>
            <tbody>
            { srcList.map((out) => {
                outCounter++;
                return (
                    <FieldArray name={`${out}.ID`} component={IOTableRow} out={out} key={`${out}.ID` + outCounter} index={outCounter}/>
                )
            })}
            </tbody>
        </Table>
    )
};


const ServiceEditForm = (props) => {
    const {handleAddIO, handleRemoveIO, service } = props;

    return (
        <form>
            <Row>
                <Col xs={12} sm={12} md={12} >

                    <FormGroup controlId="serviceNameGroup">
                        <ControlLabel>Name</ControlLabel>
                        <Field name="serviceName" component={renderField} type="text" className="form-control"/>
                    </FormGroup>
                    <FormGroup controlId="serviceDescriptionGroup">
                        <ControlLabel>Description</ControlLabel>
                        <div>
                            <Field name="serviceDescription" component="textarea" placeholder="Describe what the service does"/>
                        </div>
                    </FormGroup>
                    <FormGroup controlId="serviceVersionGroup">
                        <ControlLabel>Version</ControlLabel>
                        <Field name="serviceVersion" component="input" type="text" placeholder="Version of the service" className="form-control"/>
                    </FormGroup>

                    <FormGroup>
                        <ControlLabel>Inputs</ControlLabel>
                        <Panel>
                            <div className="form-group row">
                                <div className="col-sm-6">
                                    <ControlLabel>Name</ControlLabel>
                                    <Field name="inputSourceName" component="input" type="text" placeholder="Input name"
                                           className="form-control"/>
                                </div>
                                <div className="col-sm-6">
                                    <ControlLabel>Path</ControlLabel>
                                    <Field name="inputSourcePath" component="input" type="text" placeholder="Path in the container"
                                           className="form-control"/>
                                </div>
                            </div>
                            <div className="form-group row">
                                <div className="col-sm-6">
                                    <ControlLabel>Type</ControlLabel>
                                    <Field name="inputSourceType" component="select" type="text" placeholder="select"
                                           className="form-control">
                                        <option value="url">URL/PID</option>
                                        <option value="string">String data</option>
                                    </Field>
                                </div>
                                <div className="col-sm-6">
                                    <ControlLabel>File Name</ControlLabel>
                                    <Field name="inputSourceFileName" component="input" type="text" placeholder="File name"
                                           className="form-control"/>
                                </div>
                            </div>
                            <div className="form-group row text-center">
                                <Button type="submit" className="btn btn-default" onClick={(evt) => handleAddIO(true, evt)}>
                                    <span className="glyphicon glyphicon-plus" aria-hidden="true"></span> Add a new input
                                </Button>
                            </div>
                            <InputTable key={"in-"+service.ID} service={service} handleRemoveIO={handleRemoveIO}/>
                        </Panel>
                    </FormGroup>

                    <FormGroup>
                        <ControlLabel>Outputs</ControlLabel>
                        <div className="input-group">
                            <span className="input-group-addon">Name</span>
                            <Field name="outputSourceName" component="input" type="text" placeholder="Any name"
                                   className="form-control"/>
                            <span className="input-group-addon">Path</span>
                            <Field name="outputSourcePath" component="input" type="text" placeholder="Path in the container"
                                   className="form-control"/>
                            <span className="input-group-btn">
                                <Button type="submit" className="btn btn-default" onClick={(evt) => handleAddIO(false, evt)}>
                                    <span className="glyphicon glyphicon-plus" aria-hidden="true"></span> Add
                                </Button>
                            </span>
                        </div>
                        <OutputTable key={"out-"+service.ID} service={service} handleRemoveIO={handleRemoveIO}/>

                    </FormGroup>
                </Col>
            </Row>
        </form>
    )
};

const ServiceEdit = reduxForm({form: 'ServiceEdit'})(ServiceEditForm);

class Service extends React.Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.props.handleSubmit.bind(this);
        this.handleUpdate = this.props.handleUpdate.bind(this);
        this.handleAddIO = this.props.handleAddIO.bind(this);
        this.handleRemoveIO = this.props.handleRemoveIO.bind(this);

        this.state = {
            showServiceEditModal: false,
        };
    }

    handleServiceEditModalClose() {
        this.setState({ showServiceEditModal: false });
    }

    handleServiceEditModalOpen() {
        this.setState({showServiceEditModal: true});
    }

    handleRemoveService() {
        this.props.removeService(this.props.service.ID)
    }

    componentDidMount() {
        this.props.fetchService(this.props.service.ID);
    }

    render() {
        if(! this.props.selectedService.Service) {
            return (<div>loading</div>)
        } else {
            const {ID, Name, Description, Version} = this.props.selectedService.Service;
            let initialServiceValues = {
                serviceName: this.props.selectedService.Service.Name,
                serviceDescription: this.props.selectedService.Service.Description,
                serviceVersion: this.props.selectedService.Service.Version,
            };

            return (

                <div className="panel panel-default">
                    <div className="panel-body">
                        <div style={{margin: "1em"}}>
                            <div style={{height: "1em"}}></div>
                            {tagValueRow("Name", Name)}
                            {tagValueRow("ID", ID)}
                            {tagValueRow("Description", Description)}
                            {tagValueRow("Version", Version)}
                            <JobCreator handleSubmit={this.handleSubmit} handleServiceEditModalOpen={this.handleServiceEditModalOpen.bind(this)} handleRemoveService={this.handleRemoveService.bind(this)} service={this.props.selectedService.Service}/>
                        </div>
                    </div>

                    <div>
                        <Modal show={this.state.showServiceEditModal} onHide={this.handleServiceEditModalClose.bind(this)} className="metadata-modal">
                            <Modal.Header closeButton>
                                <Modal.Title>{this.props.selectedService.Service.Name}</Modal.Title>
                            </Modal.Header>
                            <Modal.Body className="metadata-modal-body">
                                <ServiceEdit initialValues={initialServiceValues} service={this.props.selectedService.Service} handleAddIO={this.handleAddIO.bind(this)} handleRemoveIO={this.handleRemoveIO.bind(this)}/>
                            </Modal.Body>
                            <Modal.Footer>
                                <Button className="btn btn-primary" onClick={this.handleUpdate}>Save</Button>
                                <Button onClick={this.handleServiceEditModalClose.bind(this)}>Close</Button>
                            </Modal.Footer>
                        </Modal>
                    </div>
                </div>
            )
        }
    }
}

Service.propTypes = {
    service: PropTypes.object.isRequired,
    fetchService: PropTypes.func.isRequired,
    removeService: PropTypes.func.isRequired,
    fetchServices: PropTypes.func.isRequired,
    selectedService: PropTypes.object.isRequired,
    services: PropTypes.object.isRequired,
    handleSubmit: PropTypes.func.isRequired,
    handleUpdate: PropTypes.func.isRequired,
    handleAddIO: PropTypes.func.isRequired,
    volumes: PropTypes.array.isRequired,
};

export default Service;