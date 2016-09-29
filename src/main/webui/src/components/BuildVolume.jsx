'use strict';

import React, {PropTypes} from 'react';
import Files from './Files';
import axios from 'axios';
import apiNames from '../utils/GefAPI';
import bows from 'bows';

const log = bows('BuildService');

// there is possibility that the child component is rendered before the fetch of new build is finished?
class BuildVolume extends React.Component {
    constructor(props) {
        super(props);
        this.fileUploadStart = this.props.fileUploadStart.bind(this);
        this.fileUploadSuccess = this.props.fileUploadSuccess.bind(this);
        this.fileUploadError = this.props.fileUploadError.bind(this);
        this.state = {buildID : null};
        this.getApiURL = this.getApiURL.bind(this);
    }

    getApiURL(){
        log("ApiURL get callded");
        log("buildID is", this.state.buildID);
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
            <h3>Build a Volume</h3>
            <h4>Please select and upload files, the files will be put in a named docker volume, you can use the volume as input for services later. </h4>
            <Files getApiURL={this.getApiURL} fileUploadStart={this.fileUploadStart} fileUploadSuccess={this.fileUploadSuccess} fileUploadError={this.fileUploadError} buttontText='Build Volume'/>
        </div>
    }
}

BuildVolume.propTypes = {
    fileUploadStart: PropTypes.func.isRequired,
    fileUploadSuccess: PropTypes.func.isRequired,
    fileUploadError: PropTypes.func.isRequired,
};

export default BuildVolume;
