/**
 * Created by wqiu on 18/08/16.
 */
import actionTypes from './actionTypes';
import bows from 'bows';
import axios from 'axios';
import apiNames from '../GefAPI';
import Alert from 'react-s-alert';
import { push } from 'react-router-redux';
import { toPairs } from '../utils/utils';

const log = console.log; // bows('actions');
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

function jobListFetchStart() {
    return {
        type: actionTypes.JOB_LIST_FETCH_START
    }
}

function jobListFetchSuccess(jobs) {
    return {
        type: actionTypes.JOB_LIST_FETCH_SUCCESS,
        jobs: jobs
    }
}

function jobListFetchError(errorMessage) {
    return {
        type: actionTypes.JOB_LIST_FETCH_ERROR,
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

function inspectVolumeStart() {
    return {
        type: actionTypes.INSPECT_VOLUME_START
    }
}

function inspectVolumeSuccess(data) {
    return {
        type: actionTypes.INSPECT_VOLUME_SUCCESS,
        data: data
    }
}

function inspectVolumeEmpty() {
    return {
        type: actionTypes.INSPECT_VOLUME_EMPTY,
        data: []
    }
}

function inspectVolumeError(errorMessage) {
    return {
        type: actionTypes.INSPECT_VOLUME_ERROR,
        errorMessage: errorMessage
    }
}




//TODO: catch seems to swallow all of the exceptions, not only the exceptions occurred in fetch
//async actions
//these do some extra async stuff before the real actions are dispatched
function fetchJobs() {
    return function (dispatch, getState)  {
        dispatch(jobListFetchStart());
        const resultPromise = axios.get( apiNames.jobs);
        resultPromise.then(response => {
            log('fetched jobs:', response.data.Jobs);
            dispatch(jobListFetchSuccess(response.data.Jobs));
        }).catch(err => {
            Alert.error("Cannot fetch job information from the server.");
            log("An fetch error occurred", err);
            dispatch(jobListFetchError(err));
        })
    }
}


function fetchServices() {
    return function (dispatch, getState)  {
        dispatch(servicesFetchStart());
        const resultPromise = axios.get( apiNames.services);
        resultPromise.then(response => {
            // Alert.info('Services loaded from server');
            log('fetched services:', response.data.Services);
            dispatch(servicesFetchSuccess(response.data.Services));
        }).catch(err => {
            Alert.error("Cannot fetch service information from the server.");
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
            Alert.error("Cannot fetch service information from the server.");
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
            Alert.error("Cannot fetch volume information from the server.");
            log("A fetch error occurred");
            dispatch(volumesFetchError(err));
        })
    }
}

export function inspectVolume(volumeId) {
    return function (dispatch, getState) {
        dispatch(inspectVolumeStart());
        const resultPromise = axios.get( apiNames.volumes + '/' + volumeId + "/");
        if (volumeId == null) {
            dispatch(inspectVolumeEmpty());
        }
        resultPromise.then(response => {
            dispatch(inspectVolumeSuccess(response.data))
        }).catch(err => {
            Alert.error("Cannot fetch volume content information from the server.");
            log("A fetch error occurred");
            dispatch(inspectVolumeError(err));
        })
    }
}

//this creates a new upload endpoint on the server,
//the upload endpoint can be used for building services and volumes
function getNewUploadEndpoint() {
    return function (dispatch, getState) {
        const resultPromise = axios.get( apiNames.builds);
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
        dispatch(jobListFetchStart());
        const resultPromise = axios.get( apiNames.jobs + '/' + jobId);
        resultPromise.then(response => {
            log('fetched job:', response.data);
            //don't know what to do with it yet
        }).catch(err => {
            log("An fetch error occurred");
            dispatch(jobListFetchError(err));
        })
    }
}

function handleSubmitJob() {
    return function (dispatch, getState) {
        const selectedService = getState().selectedService;
        const jobCreater = getState().form.JobCreator;
        log("selectedService", selectedService);
        var fd = new FormData();
        fd.append("serviceID", selectedService.Service.ID);
        toPairs(jobCreater.values).forEach(([k, v]) => fd.append(k, v));
        const resultPromise = axios.post( apiNames.jobs, fd);
        resultPromise.then(response => {
            Alert.info("Your job has been successfully submitted");
            dispatch(push('/jobs' + '/' + response.data.jobID));
            log("created job:", response.data)
        }).catch(err => {
            Alert.error("An error occurred during submitting your job");
            log("An error occurred during creating a job", err);
        });
        console.log("submitting current job:", fd)
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
    jobListFetchStart,
    jobListFetchSuccess,
    jobListFetchError,
    volumesFetchStart,
    volumesFetchSuccess,
    volumesFetchError,
    inspectVolumeStart,
    inspectVolumeSuccess,
    inspectVolumeError,
    fetchJobs,
    fetchServices,
    fetchService,
    fetchVolumes,
    inspectVolume,
    showErrorMessageWithTimeout,
    hideErrorMessage,
    fileUploadStart,
    fileUploadSuccess,
    fileUploadError,
    getNewUploadEndpoint,
    handleSubmitJob,

};
