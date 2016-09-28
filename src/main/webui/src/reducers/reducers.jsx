/**
 * Created by wqiu on 18/08/16.
 */
'use strict';
import SI from 'seamless-immutable';
import { combineReducers } from 'redux';
import actionTypes from '../actions/actionTypes';
import {pageNames} from '../containers/Main';

const sampleState = {
    currentPage : pageNames.browseJobs,
    isFetching : true,
    filesToUpload : [],
    services: [],
    jobs: [],
    tasks: []
};

function currentPage(state = pageNames.browseJobs, action){
    switch (action.type) {
        case actionTypes.PAGE_CHANGE:
            return SI(action.page);
        default:
            return state;
    }
}

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

const rootReducer = combineReducers({
    currentPage,
    jobs,
});

export default rootReducer;