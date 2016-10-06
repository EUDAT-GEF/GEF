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

function servicesFetchStart() {
    return {
        type: actionTypes.SERVICES_FETCH_START
    }
}

function servicesFetchSuccess(services) {
    return {
        type: actionTypes.SERVICES_FETCH_SUCCESS,
        services: services
    }
}

function servicesFetchError(errorMessage) {
    return {
        type: actionTypes.SERVICES_FETCH_ERROR,
        errorMessage: errorMessage
    }
}

function serviceFetchStart() {
    return {
        type: actionTypes.SERVICE_FETCH_START
    }
}

function serviceFetchSuccess(service) {
    return {
        type: actionTypes.SERVICE_FETCH_SUCCESS,
        service: service
    }
}

function serviceFetchError(errorMessage) {
    return {
        type: actionTypes.SERVICE_FETCH_ERROR,
        errorMessage: errorMessage
    }
}

function jobsFetchStart() {
    return {
        type: actionTypes.JOB_FETCH_START
    }
}

function jobsFetchSuccess(jobs) {
    return {
        type: actionTypes.JOB_FETCH_SUCCESS,
        jobs: jobs
    }
}

function jobsFetchError(errorMessage) {
    return {
        type: actionTypes.JOB_FETCH_ERROR,
        errorMessage: errorMessage
    }
}

function volumesFetchStart() {
    return {
        type: actionTypes.VOLUME_FETCH_START
    }
}

function volumesFetchSuccess(volumes) {
    return {
        type: actionTypes.VOLUME_FETCH_SUCCESS,
        volumes: volumes
    }
}

function volumesFetchError(errorMessage) {
    return {
        type: actionTypes.VOLUME_FETCH_ERROR,
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
        dispatch(jobsFetchStart());
        const resultPromise = axios.get( apiNames.jobs);
        resultPromise.then(response => {
            log('fetched jobs:', response.data.Jobs);
            dispatch(jobsFetchSuccess(response.data.Jobs));
        }).catch(err => {
            log("An fetch error occurred");
            dispatch(jobsFetchError(err));
        })
    }
}


function fetchServices() {
    return function (dispatch, getState)  {
        dispatch(servicesFetchStart());
        const resultPromise = axios.get( apiNames.services);
        resultPromise.then(response => {
            log('fetched services:', response.data.Services);
            dispatch(servicesFetchSuccess(response.data.Services));
        }).catch(err => {
            log("A fetch error occurred");
            dispatch(servicesFetchError(err));
        })
    }
}

function fetchService(serviceID) {
    return function (dispatch, getState) {
        dispatch(serviceFetchStart());
        const resultPromise = axios.get( apiNames.services + '/' + serviceID);
        resultPromise.then(response => {
            log('fetched service:', response.data);
            dispatch(serviceFetchSuccess(response.data));
        }).catch(err => {
            log("A fetch error occurred");
            dispatch(serviceFetchError(err));
        })
    }
}

function fetchVolumes() {
    return function (dispatch, getState) {
        dispatch(volumesFetchStart());
        const resultPromise = axios.get(apiNames.volumes);
        resultPromise.then(response => {
            log('fetched volumes:', response.data.Volumes);
            dispatch(volumesFetchSuccess(response.data.Volumes))
        }).catch(err => {
            log("A fetch error occurred");
            dispatch(volumesFetchError(err));
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
        dispatch(jobsFetchStart());
        const resultPromise = axios.get( api.jobs + '/' + jobId);
        resultPromise.then(response => {
            log('fetched job:', response.data);
            //don't know what to do with it yet
        }).catch(err => {
            log("An fetch error occurred")
            dispatch(jobsFetchError(err));
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
    servicesFetchStart,
    servicesFetchSuccess,
    servicesFetchError,
    serviceFetchStart,
    serviceFetchSuccess,
    serviceFetchError,
    jobsFetchStart,
    jobsFetchSuccess,
    jobsFetchError,
    volumesFetchStart,
    volumesFetchSuccess,
    volumesFetchError,
    fetchJobs,
    fetchServices,
    fetchService,
    fetchVolumes,
    showErrorMessageWithTimeout,
    hideErrorMessage,
    fileUploadStart,
    fileUploadSuccess,
    fileUploadError,
    prepareNewBuild
};

