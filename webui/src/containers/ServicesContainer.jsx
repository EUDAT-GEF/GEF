import Services from '../components/Services';
import {fetchServices, fetchService, handleUpdateService, addIOPort, removeIOPort, handleSubmitJob} from '../actions/actions';

import {connect} from 'react-redux';

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
            const action = fetchServices();
            dispatch(action);
        },
        fetchService: (serviceID) => {
            const action = fetchService(serviceID);
            dispatch(action);
        },
        handleUpdate: (e) => {
            e.preventDefault();
            const action = handleUpdateService();
            dispatch(action);
        },
        handleAddIO: (isInput, e) => {
            e.preventDefault();
            const action = addIOPort(isInput);
            dispatch(action);
        },
        handleRemoveIO: (isInput, index, e) => {
            e.preventDefault();
            const action = removeIOPort(isInput, index);
            dispatch(action);
        },
        handleSubmit: (e) => {
            e.preventDefault();
            const action = handleSubmitJob();
            dispatch(action);
        },
    };
};

const ServicesContainer = connect(mapStateToProps, mapDispatchToProps)(
    Services
);

export default ServicesContainer;
