/**
 * Created by wqiu on 18/08/16.
 */
import SI from 'seamless-immutable';
import { combineReducers } from 'redux';
import actionTypes from '../actions/actionTypes';
import { reducer as formReducer } from 'redux-form';
import { syncHistoryWithStore, routerReducer } from 'react-router-redux';

function jobs(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.JOB_LIST_FETCH_SUCCESS:
            return SI(action.jobs);
        case actionTypes.JOB_LIST_FETCH_ERROR:
            return SI([]);
        default:
            return state;
    }
}

function services(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.SERVICES_FETCH_SUCCESS:
            return SI(action.services);
        case actionTypes.SERVICES_FETCH_ERROR:
            return SI([]);
        default:
            return state;
    }
}

function selectedService(state = SI({}), action) {
    switch (action.type) {
        case actionTypes.SERVICE_FETCH_SUCCESS:
            return SI(action.service);
        case actionTypes.SERVICE_FETCH_ERROR:
            return SI({});
        default:
            return state;
    }
}

function volumes(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.VOLUME_FETCH_SUCCESS:
            return SI(action.volumes);
        case actionTypes.VOLUME_FETCH_ERROR:
            return SI([]);
        default:
            return state;
    }
}

function selectedVolume(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.INSPECT_VOLUME_SUCCESS:
            return SI(action.data);
        case actionTypes.INSPECT_VOLUME_EMPTY:
            return SI([]);
        case actionTypes.INSPECT_VOLUME_ERROR:
            return SI([]);
        default:
            return state;
    }
}


const rootReducer = combineReducers({
    jobs,
    services,
    volumes,
    selectedService,
    selectedVolume,
    form: formReducer,
    routing: routerReducer
});

export default rootReducer;
