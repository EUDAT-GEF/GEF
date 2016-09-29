/**
 * Created by wqiu on 18/08/16.
 */
'use strict';

import _ from 'lodash';
import actionTypes from './actionTypes';
import {pageNames} from '../containers/Main';
import bows from 'bows';
import axios from 'axios';
import apiNames from '../utils/GefAPI';

const log = bows('actions');
//sync actions
//these are just plain action creators


function pageChange(pageName) {
    return {
        type: actionTypes.PAGE_CHANGE,
        page: pageName
    }
}

function errorOccur(errorMessage) {
    return {
        type: actionTypes.ERROR_OCCUR,
        page: errorMessage
    }
}

function serviceFetchStart() {
    return {
        type: actionTypes.SERVICE_FETCH_START
    }
}

function serviceFetchSuccess(services) {
    return {
        type: actionTypes.SERVICE_FETCH_SUCCESS,
        services: services
    }
}

function serviceFetchError(errorMessage) {
    return {
        type: actionTypes.SERVICE_FETCH_ERROR,
        errorMessage: errorMessage
    }
}

function jobFetchStart() {
    return {
        type: actionTypes.JOB_FETCH_START
    }
}

function jobFetchSuccess(jobs) {
    return {
        type: actionTypes.JOB_FETCH_SUCCESS,
        jobs: jobs
    }
}

function jobFetchError(errorMessage) {
    return {
        type: actionTypes.JOB_FETCH_ERROR,
        errorMessage: errorMessage
    }
}

function fileUploadStart() {
    return {
        type: actionTypes.FILE_UPLOAD_START
    }
}

function newBuild(buildID) {
    return {
        type: actionTypes.NEW_BUILD,
        buildID: buildID
    }
}

function newBuildError(errorMessage) {
    return {
        type: actionTypes.NEW_BUILD_ERROR,
        errorMessage: errorMessage
    }
}

function fileUploadSuccess(data) {
    return {
        type: actionTypes.FILE_UPLOAD_SUCCESS,
        data: data
    }
}

function fileUploadError(errorMessage) {
    return {
        type: actionTypes.FILE_UPLOAD_ERROR,
        errorMessage: errorMessage
    }
}


//TODO: catch seems to swallow all of the exceptions, not only the exceptions occurred in fetch
//async actions
//these do some extra async stuff before the real actions are dispatched
function fetchJobs() {
    return function (dispatch, getState)  {
        dispatch(jobFetchStart());
        const resultPromise = axios.get( apiNames.jobs);
        resultPromise.then(response => {
            log('fetched jobs:', response.data.Jobs);
            dispatch(jobFetchSuccess(response.data.Jobs));
        }).catch(err => {
            log("An fetch error occurred");
            dispatch(jobFetchError(err));
        })
    }
}


function fetchServices() {
    return function (dispatch, getState)  {
        dispatch(serviceFetchStart());
        const resultPromise = axios.get( apiNames.services);
        resultPromise.then(response => {
            log('fetched services:', response.data.Services);
            dispatch(serviceFetchSuccess(response.data.Services));
        }).catch(err => {
            log("An fetch error occurred");
            dispatch(serviceFetchError(err));
        })
    }
}

//this creates a new build endpoint on the server,
//files are posted to this endpoint to construct a docker image
function prepareNewBuild() {
    return function (dispatch, getState) {
        const resultPromise = axios.get( apiNames.buildImages);
        resultPromise.then(response => {
            const buildID = response.data.buildID;
            log('Preapred a new buildID', buildID);
            dispatch(newBuild(buildID))
        }).catch(err => {
            log("failed to get a new buildID");
            dispatch(newBuildError(err));
        })
    }

}

function fetchJobById(jobId) {
    return function (dispatch, getState)  {
        dispatch(jobFetchStart());
        const resultPromise = axios.get( api.jobs + '/' + jobId);
        resultPromise.then(response => {
            log('fetched job:', response.data);
            //don't know what to do with it yet
        }).catch(err => {
            log("An fetch error occurred")
            dispatch(jobFetchError(err));
        })
    }
}

function showErrorMessageWithTimeout(id, timeout) {

}

function hideErrorMessage(id) {

}

export default {
    pageChange,
    errorOccur,
    serviceFetchStart,
    serviceFetchSuccess,
    serviceFetchError,
    jobFetchStart,
    jobFetchSuccess,
    jobFetchError,
    fetchJobs,
    fetchServices,
    showErrorMessageWithTimeout,
    hideErrorMessage,
    fileUploadStart,
    fileUploadSuccess,
    fileUploadError,
    prepareNewBuild
};

