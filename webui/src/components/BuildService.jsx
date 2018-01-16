import React from 'react';
import PropTypes from 'prop-types';
import Files from './Files';
import axios from 'axios';
import {apiNames} from '../GefAPI';
import Alert from 'react-s-alert';
import {errHandler, fetchBuildID, fetchApiInfo} from '../actions/actions';

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

        // this.buildFetchStart = this.props.buildFetchStart.bind(this);
        // this.buildFetchSuccess = this.props.buildFetchSuccess.bind(this);
        // this.buildFetchError = this.props.buildFetchError.bind(this);
    }

    getApiURL(){
        return apiNames.builds + '/' + this.state.buildID;
    }


    componentWillMount() {
        console.log("BUILD----");
        console.log(sessionStorage.getItem('buildID'));
        const resultPromise = axios.post(apiNames.builds);
        resultPromise.then(response => {
            this.setState({buildID : response.data.buildID});
        }).catch(errHandler());
    }

    render () {
        return <div>
            <h3>Build a Service</h3>
            <Files getApiURL={this.getApiURL}
              fileUploadStart={this.fileUploadStart}
              fileUploadSuccess={this.fileUploadSuccess}
              fileUploadError={this.fileUploadError}
              buildID={this.state.buildID}
              buttonText="Build Service"/>
        </div>
    }
}

BuildService.propTypes = {
    fileUploadStart: PropTypes.func.isRequired,
    fileUploadSuccess: PropTypes.func.isRequired,
    fileUploadError: PropTypes.func.isRequired,
    // buildFetchStart: PropTypes.func.isRequired,
    // buildFetchSuccess: PropTypes.func.isRequired,
    // buildFetchError: PropTypes.func.isRequired,
};

export default BuildService;
