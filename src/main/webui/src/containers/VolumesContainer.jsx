'use strict';

import Volumes from '../components/Volumes';
import actions from '../actions/actions';

import {connect} from 'react-redux';


const mapStateToProps = (state) => {
    return {
        volumes: state.volumes
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        fetchVolumes: () => {
            const action = actions.fetchVolumes();
            dispatch(action);
        }
    };
};


const VolumesContainer = connect(mapStateToProps, mapDispatchToProps)(
    Volumes
);

export default VolumesContainer;
