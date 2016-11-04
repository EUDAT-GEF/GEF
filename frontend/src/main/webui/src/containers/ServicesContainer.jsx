import Services from '../components/Services';
import actions from '../actions/actions';
import bows from 'bows';

import {connect} from 'react-redux';

const log = bows("ServiceContainer");

const mapStateToProps = (state) => {
    return {
        services: state.services,
        selectedService: state.selectedService,
        volumes: state.volumes
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        fetchServices: () => {
            const action = actions.fetchServices();
            const action2 = actions.fetchVolumes();
            dispatch(action2);
            dispatch(action);
        },
        fetchService: (serviceID) => {
            const action = actions.fetchService(serviceID);
            dispatch(action);
        },
        handleSubmit: (e) => {
            e.preventDefault();
            log("handleSubmit called");
            const action = actions.handleSubmitJob();
            dispatch(action);
        }
    };
};


const ServicesContainer = connect(mapStateToProps, mapDispatchToProps)(
    Services
);

export default ServicesContainer;
