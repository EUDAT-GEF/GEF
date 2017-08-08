import React from 'react';
import PropTypes from 'prop-types';
import Files from './Files';
import axios from 'axios';
import {apiNames} from '../GefAPI';
import Alert from 'react-s-alert';
import {errHandler} from '../actions/actions'

const log = console.log;

// there is possibility that the child component is rendered before the fetch of new build is finished?
class BuildService extends React.Component {
    constructor(props) {
        super(props);
        this.fileUploadStart = this.props.fileUploadStart.bind(this);
        this.fileUploadSuccess = this.props.fileUploadSuccess.bind(this);
        this.fileUploadError = this.props.fileUploadError.bind(this);
        this.state = {buildID : null};
        this.getApiURL = this.getApiURL.bind(this);
    }

    getApiURL(){
        // log("ApiURL get called");
        // log("buildID is", this.state.buildID);
        return apiNames.builds + '/' + this.state.buildID;
    }


    componentWillMount() {
        const resultPromise = axios.post(apiNames.builds);
        resultPromise.then(response => {
            this.setState({buildID : response.data.buildID});
        }).catch(errHandler());
    }
    render () {
        return <div>
            <h3>Build a Service</h3>
            <h4>Please select and upload the Dockerfile, together with other files which are part of the container</h4>
            <Files getApiURL={this.getApiURL} fileUploadStart={this.fileUploadStart} fileUploadSuccess={this.fileUploadSuccess} fileUploadError={this.fileUploadError} buttonText="Build Service"/>
        </div>
    }
}

BuildService.propTypes = {
    fileUploadStart: PropTypes.func.isRequired,
    fileUploadSuccess: PropTypes.func.isRequired,
    fileUploadError: PropTypes.func.isRequired,
};

export default BuildService;
