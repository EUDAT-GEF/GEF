/**
 * Created by wqiu on 18/08/16.
 */
'use strict';
import SI from 'seamless-immutable';
import { combineReducers } from 'redux';
import actionTypes from '../actions/actionTypes';

const sampleState = {
    currentPage : 'BrowseFiles',
    isFetching : true,
    filesToUpload : [],
    services: [],
    jobs: [],
    tasks: []
};

function currentPage(state = 'buildService', action){
    switch (action.type) {
        case actionTypes.PAGE_CHANGE:
            return SI(action.data.page);
        default:
            return state;
    }
}


const rootReducer = combineReducers({
    currentPage
});

export default rootReducer;