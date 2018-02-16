import React from 'react';
import PropTypes from 'prop-types';
import ReactDOMServer from 'react-dom/server';
import DropzoneComponent from 'react-dropzone-component';
import {Row, Col, Button, Glyphicon} from 'react-bootstrap';
import {apiNames} from '../GefAPI';
import axios from 'axios';
import bows from 'bows';
import {errHandler} from '../actions/actions'

require('react-dropzone-component/styles/filepicker.css');
require('dropzone/dist/min/dropzone.min.css');


const log = bows('Files');
const BuildProgress = ({isInProgress, statusMessage, errorMessage}) => {
    if (!isInProgress) {
        return <div className="text-center">{statusMessage} {(errorMessage ? errorMessage : "")}</div>
    } else {
        return <div className="text-center"><img src="/images/progress-animation.gif" /> {statusMessage}</div>
    }
};

const GoBackLink = ({id, buildFinished}) => {
    if ((id) && (buildFinished)) {
        return <div className="text-center"><a href="/builds">Build another service</a></div>
    } else {
      return <div></div>
    }
};

class Files extends React.Component {
    constructor(props) {
        super(props);
        this.djsConfig = {
          disabled: true,
            addRemoveLinks: true,
            autoProcessQueue: false,
            uploadMultiple: true,
            parallelUploads: 999,
            maxFilesize: 4000,
            previewTemplate: ReactDOMServer.renderToStaticMarkup(
                <div className="dz-preview dz-file-preview">
                    <div className="dz-filename"><span data-dz-name="true"></span></div>
                    <img data-dz-thumbnail="true" />
                    <div className="dz-progress"><span className="dz-upload" data-dz-uploadprogress="true"></span></div>
                    <div className="dz-success-mark"><span>✔</span></div>
                    <div className="dz-error-mark"><span>✘</span></div>
                    <div className="dz-error-message"><span data-dz-errormessage="true"></span></div>
                </div>
            )
        };

        this.state = {
            myDropzone: undefined,
            serviceBuildInProgress: false,
            statusMessage: "Ready to build a service",
            errorMessage: "",
            build : null,
            buildID: null,
            buildFinished: false,
        };
        this.fileUploadSuccess = this.props.fileUploadSuccess.bind(this);
        this.fileUploadError = this.props.fileUploadError.bind(this);
    }

    componentDidMount() {
      if (sessionStorage.getItem("buildID")) {
          this.setState({serviceBuildInProgress: true})
      }
      let buildStatusUpdateTimer = setInterval(() => this.tick(), 1000);
      this.setState({buildStatusUpdateTimer: buildStatusUpdateTimer, buildFinished: false});
    }

    componentWillUnmount() {
        clearInterval(this.state.buildStatusUpdateTimer);
    }

    tick() {
        if (sessionStorage.getItem("buildID")) {
            const resultPromise = axios.get(apiNames.builds + '/' + sessionStorage.getItem("buildID"));
            resultPromise.then(response => {
                this.setState({build : response.data.Build});
                let build = response.data.Build;

                this.setState({statusMessage: build.State.Status});
                this.setState({errorMessage: build.State.Error});
                if (build.State.Code>-1) {
                    clearInterval(this.state.buildStatusUpdateTimer);
                    this.setState({serviceBuildInProgress: false, buildFinished: true});
                    sessionStorage.removeItem("buildID");
                }
            }).catch(errHandler());
        }
    }

    render() {
        const getApiURL = this.props.getApiURL;
        const buttonText = this.props.buttonText;
        const djsConfig = this.djsConfig;

        const eventHandlers = {
            init: (passedDropZone) => {
                this.setState({myDropzone: passedDropZone});
            },

            successmultiple: (files, response) => {
                log('successmultiple, response is: ', response);
                this.fileUploadSuccess(response);
                sessionStorage.setItem('buildID', response.buildID);
                this.setState({
                    serviceBuildInProgress: true,
                    statusMessage: "Files have been successfully uploaded. Starting to build a service...",
                    buildID: response.buildID,
                });
                this.tick();
            },

            error: (files, errorMessage) => {
                this.fileUploadError(errorMessage);
                this.setState({serviceBuildInProgress: false, statusMessage: errorMessage });
            }
        };

        const fileUploadStart = this.props.fileUploadStart.bind(this);
        const submitHandler = ()  => {
            fileUploadStart();
            this.state.myDropzone.processQueue();
        };

        if(getApiURL() != null) {
            const config = {
                postUrl: getApiURL()
            };
            let isDropZoneVisible = (sessionStorage.getItem("buildID")) ? false : true;
            return <div>
                {isDropZoneVisible == true &&
                    <span>
                        <h4>Please select and upload the Dockerfile, together with other files which are part of the container</h4>
                        <DropzoneComponent config={config} eventHandlers={eventHandlers} djsConfig={djsConfig} />
                        <Row>
                            <Col md={4} mdOffset={4}> <Button type='submit' bsStyle='primary' style={{width: '100%'} } onClick={submitHandler}> <Glyphicon glyph='upload'/> {buttonText} </Button> </Col>
                        </Row>
                    </span>
                }
                <span>
                    <BuildProgress isInProgress={this.state.serviceBuildInProgress} statusMessage={this.state.statusMessage} errorMessage={this.state.errorMessage}/>
                    <GoBackLink id={sessionStorage.getItem("buildID")} buildFinished={this.state.buildFinished}/>
                </span>
            </div>
        } else {
            return <div> loading </div>
        }
    }
}

Files.propTypes = {
    buttonText: PropTypes.string.isRequired,
    fileUploadStart: PropTypes.func.isRequired,
    getApiURL: PropTypes.func.isRequired,
    fileUploadSuccess: PropTypes.func.isRequired,
    fileUploadError: PropTypes.func.isRequired,
    buildID: PropTypes.string,
};


export default Files;
