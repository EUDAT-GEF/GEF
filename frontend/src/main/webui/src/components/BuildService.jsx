import React, {PropTypes} from 'react';
import Files from './Files';
import axios from 'axios';
import apiNames from '../GefAPI';
import bows from 'bows';

const log = bows('BuildService');

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
        return apiNames.buildImages + '/' + this.state.buildID;
    }


    componentWillMount() {
        const resultPromise = axios.post(apiNames.buildImages);
        resultPromise.then(response => {
            this.setState({buildID : response.data.buildID});
            log('New service URL:', this.state.buildID);
        }).catch(err => {
            log("An error occurred during creating new service URL");
        });
    }
    render () {
        return <div>
            <h3>Build a Service</h3>
            <h4>Please select and upload the Dockerfile, together with other files which are part of the container</h4>
            <Files getApiURL={this.getApiURL} fileUploadStart={this.fileUploadStart} fileUploadSuccess={this.fileUploadSuccess} fileUploadError={this.fileUploadError} buttonText="Build Image"/>
        </div>
    }
}

BuildService.propTypes = {
    fileUploadStart: PropTypes.func.isRequired,
    fileUploadSuccess: PropTypes.func.isRequired,
    fileUploadError: PropTypes.func.isRequired,
};

export default BuildService;
