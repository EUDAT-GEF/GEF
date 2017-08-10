/**
 * Created by wqiu on 18/08/16.
 */
import SI from 'seamless-immutable';
import { combineReducers } from 'redux';
import actionTypes from '../actions/actionTypes';
import { reducer as formReducer } from 'redux-form';
import { syncHistoryWithStore, routerReducer } from 'react-router-redux';

function apiinfo(state = SI({}), action) {
    switch (action.type) {
        case actionTypes.APIINFO_FETCH_SUCCESS:
            return SI(action.apiinfo);
        default:
            return state;
    }
}

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

function task(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.CONSOLE_OUTPUT_FETCH_SUCCESS:
            return SI(action.data);
        case actionTypes.CONSOLE_OUTPUT_FETCH_EMPTY:
            return SI([]);
        case actionTypes.CONSOLE_OUTPUT_FETCH_ERROR:
            return SI([]);
        default:
            return state;
    }
}

function user(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.USER_FETCH_SUCCESS:
            return SI(action.user);
        case actionTypes.USER_FETCH_ERROR:
            return SI([]);
        default:
            return state;
    }
}

function tokens(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.USER_TOKENS_FETCH_SUCCESS:
            return SI(action.tokens);
        case actionTypes.USER_TOKENS_FETCH_ERROR:
            return SI([]);
        default:
            return state;
    }
}

function roles(state = SI({}), action) {
    switch (action.type) {
        case actionTypes.ROLES_FETCH_SUCCESS:
            action.roles.map(r => state = state.set(r.ID,  r));
            return state;
        case actionTypes.ROLE_USERS_FETCH_SUCCESS:
            state = SI.setIn(
                state,
                [action.roleID, 'users'],
                SI(action.roleUsers))
            return state;
        case actionTypes.ROLES_FETCH_ERROR:
            return SI([]);
        case actionTypes.ROLE_USERS_FETCH_ERROR:
            return SI([]);
        default:
            return state;
    }
}


const rootReducer = combineReducers({
    apiinfo,
    jobs,
    services,
    volumes,
    selectedService,
    selectedVolume,
    task,
    user,
    tokens,
    roles,
    form: formReducer,
    routing: routerReducer
});

export default rootReducer;
