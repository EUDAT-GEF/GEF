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

function selectedVolumeContent(state = SI([]), action) {
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

function selectedVolumeID(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.SELECT_VOLUME_SUCCESS:
            return SI(action.data);
        case actionTypes.SELECT_VOLUME_EMPTY:
            return SI([]);
        case actionTypes.SELECT_VOLUME_ERROR:
            return SI([]);
        default:
            return state;
    }
}

function downloadedFile(state = SI([]), action) {
    switch (action.type) {
        case actionTypes.DOWNLOAD_VOLUME_FILE_SUCCESS:


            var blob = new Blob([action.data]);
            if (window.navigator.msSaveOrOpenBlob)  // IE hack; see http://msdn.microsoft.com/en-us/library/ie/hh779016.aspx
                window.navigator.msSaveBlob(blob, "filename.txt");
            else
            {
                var a = window.document.createElement("a");
                a.href = window.URL.createObjectURL(blob, {type: "text/plain"});
                a.download = "filename.csv";
                document.body.appendChild(a);
                a.click();  // IE: "Access is denied"; see: https://connect.microsoft.com/IE/feedback/details/797361/ie-10-treats-blob-url-as-cross-origin-and-denies-access
                document.body.removeChild(a);
            }

            return SI(action.data);
        case actionTypes.DOWNLOAD_VOLUME_FILE_START:
            return SI([]);
        case actionTypes.DOWNLOAD_VOLUME_FILE_ERROR:
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
    selectedVolumeContent,
    selectedVolumeID,
    downloadedFile,
    form: formReducer,
    routing: routerReducer
});

export default rootReducer;
