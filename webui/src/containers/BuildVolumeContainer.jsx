import BuildVolume from '../components/BuildVolume';
import actions from '../actions/actions';

import {connect} from 'react-redux';
import { push } from 'react-router-redux';

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
            dispatch(push('/volumes/' + response.Volume.ID));
            dispatch(action)
        },

        fileUploadError: (files, errorMessage) => {
            const action = actions.fileUploadError(errorMessage);
            dispatch(action);
        }
    }
};

const BuildVolumeContainer = connect(mapStateToProps, mapDispatchToProps)(BuildVolume);

export default BuildVolumeContainer;
