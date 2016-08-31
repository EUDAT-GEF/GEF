/**
 * Created by wqiu on 18/08/16.
 */
'use strict';

import _ from 'lodash';
import actionTypes from './actionTypes';
import {pageNames} from '../containers/Main';
import bows from 'bows';
import axios from 'axios';

const log = bows('actions');
//sync actions
//these are just plain action creators

const apiNames = {
    datasets: "/gef/api/datasets",
    builds:   "/gef/api/builds",
    services: "/gef/api/images",
    jobs: "/gef/api/jobs",
};

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

function serviceFetchSuccess() {
    return {
        type: actionTypes.SERVICE_FETCH_SUCCESS
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

//async actions
//these do some extra async stuff before the real actions are dispatched
function fetchJobs() {
    return function (dispatch, getState)  {
        dispatch(jobFetchStart());
        const resultPromise = axios.get(apiNames.jobs);
        resultPromise.then(response => {
            log('fetched jobs:', response.data.Jobs);
            dispatch(jobFetchSuccess(response.data.Jobs));
        }).catch(err => {
            log("An fetch error occurred")
            dispatch(jobFetchError(err));
        })
    }
}

function fetchJobById(jobId) {
    return function (dispatch, getState)  {
        dispatch(jobFetchStart());
        const resultPromise = axios.get(apiNames.jobs + '/' + jobId);
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
    showErrorMessageWithTimeout,
    hideErrorMessage
};

