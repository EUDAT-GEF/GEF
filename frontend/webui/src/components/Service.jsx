/**
 * Created by wqiu on 17/08/16.
 */
import React, {PropTypes} from 'react';
import { Row, Col, Grid, Table, Button, Modal, OverlayTrigger, FormGroup, ControlLabel } from 'react-bootstrap';
import {Field, FieldArray, reduxForm} from 'redux-form';
import validate from './ServiceMetadataValidator'

// this is a detailed view of a service, user will be able to execute service in this view


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




const InputTable = ({service}) => {
    let inCounter = 0;
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
            { service.Input.map((src) => {
                inCounter++;
                return (
                    <tr key={"src-"+inCounter+"-"+service.ID}>
                        <td>{src.ID}</td>
                        <td>{src.Name}</td>
                        <td>{src.Path}</td>
                        <td>
                            <Button type="submit" bsStyle="primary" bsSize="xsmall">
                                <span className="glyphicon glyphicon-remove" aria-hidden="true"></span> Remove
                            </Button>
                        </td>
                    </tr>
                )

            })}

            </tbody>
        </Table>
    )
};

const OutputTable = ({service}) => {

    let outCounter = 0;
    let outputs = [];
    service.Output.map((out, index) => {
        outputs.push(out);
    });
    //outputs.push({});

    const IOTableRow = ({out, index}) => {
        return (
            <tr>
                <td>{out.ID}</td>
                <td>{out.Name}</td>
                <td>{out.Path}</td>
                <td>
                    <Button type="submit" bsStyle="primary" bsSize="xsmall" onClick={() => outputs.push({})}>
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

            { service.Output.map((out, index) => {
                outCounter++;
                return (
                    <FieldArray name={`${out}.ID`} component={IOTableRow} out={out}/>
                )
            })}

            </tbody>
        </Table>

    )
};





const ServiceEditForm = (props) => {
    const { handleUpdate, handleAddOutput, service } = props;
    
    return (
        <form>
            <Row>
                <Col xs={12} sm={12} md={12} >

                    <FormGroup controlId="serviceNameGroup">
                        <ControlLabel>Name</ControlLabel>
                        <Field name="serviceName" component="input" type="text" className="form-control"/>
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
                        <div className="input-group">
                            <span className="input-group-addon">Name</span>
                            <Field name="inputSourceName" component="input" type="text" placeholder="Any name"
                                   className="form-control"/>
                            <span className="input-group-addon">Path</span>
                            <Field name="inputSourcePath" component="input" type="text" placeholder="Path in the container"
                                   className="form-control"/>
                            <span className="input-group-btn">
                                <Button type="submit" className="btn btn-default">
                                    <span className="glyphicon glyphicon-plus" aria-hidden="true"></span> Add
                                </Button>
                            </span>
                        </div>
                        <InputTable key={"in-"+service.ID} service={service}/>
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
                                <Button type="submit" className="btn btn-default" onClick={handleAddOutput}>
                                    <span className="glyphicon glyphicon-plus" aria-hidden="true"></span> Add
                                </Button>
                            </span>
                        </div>
                        <OutputTable key={"out-"+service.ID} service={service}/>


                        <Field name="outputHidden[0]" component="hidden" type="text"/>
                        <Field name="outputHidden[1]" component="hidden" type="text"/>
                        <Field name="outputHidden[2]" component="hidden" type="text"/>
                    </FormGroup>

                    <Button type="submit" className="btn btn-primary" onClick={handleUpdate}>
                        <span className="glyphicon glyphicon-floppy-disk" aria-hidden="true"></span> Save
                    </Button>





                </Col>
            </Row>
        </form>
    )
};

/*const validate = values => {
    const errors = {}
    const requiredFields = [ 'serviceName', 'serviceDescription', 'serviceVersion' ]
    requiredFields.forEach(field => {
        if (!values[ field ]) {
            errors[ field ] = 'Required'
        }
    })
    /*if (values.email && !/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i.test(values.email)) {
        errors.email = 'Invalid email address'
    }
    console.log(errors);
    return errors
}*/


const ServiceEdit = reduxForm({form: 'ServiceEdit', validate} )(ServiceEditForm);





class Service extends React.Component {



    constructor(props) {
        super(props);
        this.handleSubmit = this.props.handleSubmit.bind(this);
        this.handleUpdate = this.props.handleUpdate.bind(this);
        this.handleAddOutput = this.props.handleAddOutput.bind(this);

        this.state = {
            showModal: false,
            changedService: this.props.selectedService.Service,
        };
    }

    handleModalClose() {
        this.setState({ showModal: false });
    }

    handleModalOpen() {
        this.setState({showModal: true});
    }


    componentDidMount() {
        this.props.fetchService(this.props.service.ID);
    }

    renderModalWindow(inService) {
        let initialServiceValues = {
            serviceName: inService.Name,
            serviceDescription: inService.Description,
            serviceVersion: inService.Version,
            outputHidden: ["sometext", "another text", "third text"],
        };

        return (
            <div>
                <Modal show={this.state.showModal} onHide={this.handleModalClose.bind(this)}>
                    <Modal.Header closeButton>
                        <Modal.Title>{inService.Name}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <ServiceEdit handleUpdate={this.handleUpdate} handleAddOutput={this.handleAddOutput} service={inService} initialValues={initialServiceValues}/>
                    </Modal.Body>
                    <Modal.Footer>
                        <Button className="btn btn-primary" onClick={this.handleUpdate}>Save</Button>
                        <Button className="btn btn-primary" onClick={this.handleAddOutput}>Add output</Button>
                        <Button onClick={this.handleModalClose.bind(this)}>Close</Button>
                    </Modal.Footer>
                </Modal>
            </div>
        )
    }

    render() {







        if(! this.props.selectedService.Service) {
            return (<div>loading</div>)
        } else {
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
                            <button type="submit" className="btn btn-default" onClick={this.handleModalOpen.bind(this)}>
                                <span className="glyphicon glyphicon-edit" aria-hidden="true"></span> Edit Metadata
                            </button>

                            <div style={{height: "1em"}}></div>
                        </div>
                    </div>
                    {this.renderModalWindow(this.props.selectedService.Service)}
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
    handleUpdate: PropTypes.func.isRequired,
    handleAddOutput: PropTypes.func.isRequired,
    volumes: PropTypes.array.isRequired,
};

export default Service;