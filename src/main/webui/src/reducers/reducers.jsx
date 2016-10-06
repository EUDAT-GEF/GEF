/**
 * Created by wqiu on 18/08/16.
 */
'use strict';
import SI from 'seamless-immutable';
import { combineReducers } from 'redux';
import actionTypes from '../actions/actionTypes';
import {pageNames} from '../containers/Main';

const sampleState = {
    isFetching : true,
    services: [],
    jobs: [],
    workflows: []
};

function jobs(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.JOB_FETCH_SUCCESS:
            return SI(action.jobs);
        case actionTypes.JOB_FETCH_ERROR:
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

function selectedService(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.SERVICE_FETCH_SUCCESS:
            return SI(action.service);
        case actionTypes.SERVICE_FETCH_ERROR:
            return SI([]);
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

const rootReducer = combineReducers({
    jobs,
    services,
    volumes,
    currentService,
});

export default rootReducer;