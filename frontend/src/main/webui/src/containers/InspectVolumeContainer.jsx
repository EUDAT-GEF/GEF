import Volume from '../components/Volume';
import actions from '../actions/actions';

import {connect} from 'react-redux';


const mapStateToProps = (state) => {
    return {
        volume: state.volume
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        fetchVolumeContent: () => {
                const action = actions.fetchVolumeContent();
                dispatch(action);
        }
    };
};


const InspectVolumeContainer = connect(mapStateToProps, mapDispatchToProps)(
   Volume
);

export default InspectVolumesContainer;
