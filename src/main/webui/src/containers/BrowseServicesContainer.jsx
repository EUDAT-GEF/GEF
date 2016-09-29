'use strict';

import BrowseServices from '../components/Services';
import actions from '../actions/actions';

import {connect} from 'react-redux';


const mapStateToProps = (state) => {
    return {
        services: state.services
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        fetchServices: () => {
            const action = actions.fetchServices();
            dispatch(action);
        }
    };
};


const BrowseServicesContainer = connect(mapStateToProps, mapDispatchToProps)(
    BrowseServices
);

export default BrowseServicesContainer;
