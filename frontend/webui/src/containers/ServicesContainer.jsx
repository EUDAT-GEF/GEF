import Services from '../components/Services';
import actions from '../actions/actions';

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
            const action = actions.fetchServices();
            dispatch(action);
        },
        fetchService: (serviceID) => {
            const action = actions.fetchService(serviceID);
            dispatch(action);
        },
        handleUpdate: (e) => {
            e.preventDefault();
            console.log("UPDATE PRESSED");
            //console.log(changedService);
            //console.log(state.services);
           // var changedService = {};
            const action = actions.handleUpdateService();
            dispatch(action);
        },
        handleSubmit: (e) => {
            e.preventDefault();
            const action = actions.handleSubmitJob();
            dispatch(action);
        },

    };
};


const ServicesContainer = connect(mapStateToProps, mapDispatchToProps)(
    Services
);

export default ServicesContainer;
