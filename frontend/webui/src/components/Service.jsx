/**
 * Created by wqiu on 17/08/16.
 */
import React, {PropTypes} from 'react';
import bows from 'bows';
import { Row, Col, Grid, Table, Button, Modal, OverlayTrigger, FormGroup, ControlLabel, FormControl,  } from 'react-bootstrap';
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


const IOTable = ({service}) => {
    return (
        <Table responsive>
            <tbody>
                { service.Input.map((src) => {
                    return (
                        <tr>
                            <td>1</td>
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

const ServiceEditForm = (props) => {
    const { handleSubmit, pristine, reset, submitting, service } = props;

    return (
        <form onSubmit={handleSubmit}>
            <Row>
                <Col xs={12} sm={12} md={12} >

                    <FormGroup controlId="srvName">
                        <ControlLabel>Name</ControlLabel>
                        <FormControl type="text" value={props.service.Name}/>
                    </FormGroup>
                    <FormGroup controlId="srvDescription">
                        <ControlLabel>Description</ControlLabel>
                        <FormControl componentClass="textarea" value={props.service.Description}/>
                    </FormGroup>
                    <FormGroup controlId="srvVersion">
                        <ControlLabel>Version</ControlLabel>
                        <FormControl type="text" value={props.service.Version}/>
                    </FormGroup>

                    <FormGroup>
                        <ControlLabel>Inputs</ControlLabel>
                        <div className="input-group">
                            <span className="input-group-addon">Name</span>
                            <FormControl type="text" placeholder="Name"/>
                            <span className="input-group-addon">Path</span>
                            <FormControl type="text" placeholder="Path"/>

                            <span className="input-group-btn">
                                <button type="submit" className="btn btn-default">
                                    <span className="glyphicon glyphicon-plus" aria-hidden="true"></span> Add
                                </button>
                            </span>
                        </div>

                        <IOTable service={props.service}/>

                        <Table responsive>
                            <tbody>
                            <tr>
                                <td>1</td>
                                <td>Input 1</td>
                                <td>Input 1</td>
                                <td>
                                    <Button type="submit" bsStyle="primary" bsSize="xsmall">
                                        <span className="glyphicon glyphicon-remove" aria-hidden="true"></span> Remove
                                    </Button>
                                </td>
                            </tr>
                            <tr>
                                <td>2</td>
                                <td>Input 2</td>
                                <td>Input 2</td>
                                <td>
                                    <Button type="submit" bsStyle="primary" bsSize="xsmall">
                                        <span className="glyphicon glyphicon-remove" aria-hidden="true"></span> Remove
                                    </Button>
                                </td>
                            </tr>
                            <tr>
                                <td>3</td>
                                <td>Input 3</td>
                                <td>Input 3</td>
                                <td>
                                    <Button type="submit" bsStyle="primary" bsSize="xsmall">
                                        <span className="glyphicon glyphicon-remove" aria-hidden="true"></span> Remove
                                    </Button>
                                </td>
                            </tr>

                            </tbody>
                        </Table>


                    </FormGroup>





                    <FormGroup>
                        <ControlLabel>Outputs</ControlLabel>
                        <div className="input-group">
                            <span className="input-group-addon">Name</span>
                            <FormControl type="text" placeholder="Name"/>
                            <span className="input-group-addon">Path</span>
                            <FormControl type="text" placeholder="Path"/>

                            <span className="input-group-btn">
                                <button type="submit" className="btn btn-default">
                                    <span className="glyphicon glyphicon-plus" aria-hidden="true"></span> Add
                                </button>
                            </span>
                        </div>



                        <Table responsive>
                            <tbody>
                            <tr>
                                <td>1</td>
                                <td>Input 1</td>
                                <td>Input 1</td>
                                <td>
                                    <Button type="submit" bsStyle="primary" bsSize="xsmall">
                                        <span className="glyphicon glyphicon-remove" aria-hidden="true"></span> Remove
                                    </Button>
                                </td>
                            </tr>
                            <tr>
                                <td>2</td>
                                <td>Input 2</td>
                                <td>Input 2</td>
                                <td>
                                    <Button type="submit" bsStyle="primary" bsSize="xsmall">
                                        <span className="glyphicon glyphicon-remove" aria-hidden="true"></span> Remove
                                    </Button>
                                </td>
                            </tr>
                            <tr>
                                <td>3</td>
                                <td>Input 3</td>
                                <td>Input 3</td>
                                <td>
                                    <Button type="submit" bsStyle="primary" bsSize="xsmall">
                                        <span className="glyphicon glyphicon-remove" aria-hidden="true"></span> Remove
                                    </Button>
                                </td>
                            </tr>

                            </tbody>
                        </Table>


                    </FormGroup>


                    <Button type="submit" className="btn btn-primary" onClick={handleSubmit}>
                        <span className="glyphicon glyphicon-floppy-disk" aria-hidden="true"></span> Save
                    </Button>


                </Col>
            </Row>

        </form>
    )
};

const ServiceEdit = reduxForm({form: 'ServiceEdit'} )(ServiceEditForm);

class Service extends React.Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.props.handleSubmit.bind(this);

        this.state = {
            showModal: false,
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

    renderModalWindow(curService) {
        return (
            <div>
                <Modal show={this.state.showModal} onHide={this.handleModalClose.bind(this)}>
                    <Modal.Header closeButton>
                        <Modal.Title>{curService.Name}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <ServiceEdit handleSubmit={this.handleSubmit} service={this.props.selectedService.Service}/>


                    </Modal.Body>
                    <Modal.Footer>
                        <Button onClick={this.handleModalClose.bind(this)}>Save</Button>
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
    volumes: PropTypes.array.isRequired,
};

export default Service;