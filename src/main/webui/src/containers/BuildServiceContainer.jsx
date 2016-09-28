'use strict';


import BuildService from '../components/BuildService';
import actions from '../actions/actions';

import {connect} from 'react-redux';

const mapStateToProps = (state) => {
    return {
        files: state.files
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        fileUploadStart: () => {
            const action = actions.fileUploadStart();
            dispatch(action);
        },

        fileUploadSuccess: (response) => {
            const action = actions.fileUploadSuccess(response);
            dispatch(action);
        },

        fileUploadError: (files, errorMessage) => {
            const action = actions.fileUploadError(errorMessage);
            dispatch(action);
        }
    }
};

const BuildServiceContainer = connect(mapStateToProps, mapDispatchToProps)(BuildService);

export default BuildServiceContainer;
