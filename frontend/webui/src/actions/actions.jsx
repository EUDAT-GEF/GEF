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

function serviceUpdateStart() {
    return {
        type: actionTypes.SERVICE_UPDATE_START
    }
}

function serviceUpdateSuccess() {
    return {
        type: actionTypes.SERVICE_UPDATE_SUCCESS,
    }
}

function serviceUpdateError(errorMessage) {
    return {
        type: actionTypes.SERVICE_UPDATE_ERROR,
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

function jobRemovalStart() {
    return {
        type: actionTypes.JOB_REMOVAL_START
    }
}

function jobRemovalSuccess(data) {
    return {
        type: actionTypes.JOB_REMOVAL_SUCCESS,
        data: data
    }
}

function jobRemovalError(errorMessage) {
    return {
        type: actionTypes.JOB_REMOVAL_ERROR,
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

function consoleOutputFetchStart() {
    return {
        type: actionTypes.CONSOLE_OUTPUT_FETCH_START
    }
}

function consoleOutputFetchSuccess(data) {
    return {
        type: actionTypes.CONSOLE_OUTPUT_FETCH_SUCCESS,
        data: data
    }
}

function consoleOutputFetchEmpty() {
    return {
        type: actionTypes.CONSOLE_OUTPUT_FETCH_EMPTY,
        data: []
    }
}

function consoleOutputFetchError(errorMessage) {
    return {
        type: actionTypes.CONSOLE_OUTPUT_FETCH_ERROR,
        errorMessage: errorMessage
    }
}

function ioAddStart() {
    return {
        type: actionTypes.IO_ADD_START
    }
}

function ioAddSuccess(data) {
    return {
        type: actionTypes.IO_ADD_SUCCESS,
        data: data
    }
}

function ioAddError(errorMessage) {
    return {
        type: actionTypes.IO_ADD_ERROR,
        errorMessage: errorMessage
    }
}

function ioRemoveStart() {
    return {
        type: actionTypes.IO_REMOVE_START
    }
}

function ioRemoveSuccess(data) {
    return {
        type: actionTypes.IO_REMOVE_SUCCESS,
        data: data
    }
}

function ioRemoveError(errorMessage) {
    return {
        type: actionTypes.IO_REMOVE_ERROR,
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

function removeJob(jobID) {
    return function (dispatch, getState)  {
        dispatch(jobRemovalStart());
        const resultPromise = axios.delete( apiNames.jobs + "/" + jobID);
        resultPromise.then(response => {
            log('removed job:', response.data);
            dispatch(jobRemovalSuccess(response.data));
            dispatch(fetchJobs());
        }).catch(err => {
            Alert.error("Cannot remove the job.");
            log("An error occurred during the job removal", err);
            dispatch(jobRemovalError(err));
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

        if (!volumeId) {
            dispatch(inspectVolumeEmpty());
        } else {
            const resultPromise = axios.get( apiNames.volumes + '/' + volumeId + "/");
            resultPromise.then(response => {
                dispatch(inspectVolumeSuccess(response.data))
            }).catch(err => {
                Alert.error("Cannot fetch volume content information from the server.");
                log("A fetch error occurred");
                dispatch(inspectVolumeError(err));
            })
        }
    }
}

export function consoleOutputFetch(jobId) {
    return function (dispatch, getState) {
        dispatch(consoleOutputFetchStart());

        if (!jobId) {
            dispatch(consoleOutputFetchEmpty());
        } else {
            const resultPromise = axios.get( apiNames.jobs + '/' + jobId + "/output");
            resultPromise.then(response => {
                dispatch(consoleOutputFetchSuccess(response.data))
            }).catch(err => {
                Alert.error("Cannot fetch the console content.");
                log("A fetch error occurred");
                dispatch(consoleOutputFetchError(err));
            })
        }
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
function handleUpdateService() {
    return function (dispatch, getState) {
        const selectedService = getState().selectedService;
        const allServices = getState().services;
        const serviceEdit = getState().form.ServiceEdit;

        let outputObject =  {
            'Created': selectedService.Service.Created,
            'Description': serviceEdit.values.serviceDescription, // modified
            'ID': selectedService.Service.ID,
            'ImageID': selectedService.Service.ImageID,
            'Input':  selectedService.Service.Input,
            'Name': serviceEdit.values.serviceName, // modified
            'Output': selectedService.Service.Output,
            'RepoTag': selectedService.Service.RepoTag,
            'Size': selectedService.Service.Size,
            'Version': serviceEdit.values.serviceVersion // modified
        };

        dispatch(serviceUpdateStart());
        const resultPromise = axios.put(apiNames.services + '/' + selectedService.Service.ID, outputObject);

        resultPromise.then(response => {
            log('updated service:', response.data);
            Alert.info("Service metadata has been successfully updated");
            // Updating the list of services on the client side without requesting data from the server
            let updatedServices = [];
            let responseService = response.data.Service;
            allServices.map((curService) => {
                if (curService.ID == responseService.ID) {
                    updatedServices.push(responseService);
                } else {
                    updatedServices.push(curService);
                }
            });

            dispatch(serviceUpdateSuccess());
            dispatch(serviceFetchSuccess(response.data));
            dispatch(servicesFetchSuccess(updatedServices)); // forcing to update the list of services
        }).catch(err => {
            Alert.error("Cannot update the service.");
            log("An update error occurred");
            dispatch(serviceUpdateError(err));
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


function addIOPort(isInput) {
    return function (dispatch, getState)  {
        const selectedService = getState().selectedService;
        const serviceEdit = getState().form.ServiceEdit;
        dispatch(ioAddStart());

        let inputs = [];
        let newInput = {};
        let outputs = [];
        let newOutput = {};
        let selectedInput = [];
        let selectedOutput = [];
        if (selectedService.Service.Input) {
            selectedInput = selectedService.Service.Input
        }
        if (selectedService.Service.Output) {
            selectedOutput = selectedService.Service.Output;
        }
        if (isInput) {
            newInput.ID = "input" + selectedInput.length;
            newInput.Name = serviceEdit.values.inputSourceName;
            newInput.Path = serviceEdit.values.inputSourcePath;
            if ((!newInput.Name) && (!newInput.Path)){
                Alert.error("Input name and path cannot be empty");
                dispatch(ioAddError());
            } else {
                selectedInput.map((input) => {
                    inputs.push(input);
                });
                inputs.push(newInput);
            }

            outputs = selectedOutput;
        } else {
            newOutput.ID = "output" + selectedOutput.length;
            newOutput.Name = serviceEdit.values.outputSourceName;
            newOutput.Path = serviceEdit.values.outputSourcePath;
            if ((!newOutput.Name) && (!newOutput.Path)){
                Alert.error("Output name and path cannot be empty");
                dispatch(ioAddError());
            } else {
                selectedOutput.map((out) => {
                    outputs.push(out);
                });
                outputs.push(newOutput);
            }

            inputs = selectedInput;
        }

        let outputObject = {
            'Created': selectedService.Service.Created,
            'Description': selectedService.Service.Description,
            'ID': selectedService.Service.ID,
            'ImageID': selectedService.Service.ImageID,
            'Input': inputs,
            'Name': selectedService.Service.Name,
            'Output': outputs,
            'RepoTag': selectedService.Service.RepoTag,
            'Size': selectedService.Service.Size,
            'Version': selectedService.Service.Version
        };

        if ((inputs.length>0) || (outputs.length>0)) {
            dispatch(ioAddSuccess(outputObject));
            dispatch(serviceUpdateStart());
            const resultPromise = axios.put(apiNames.services + '/' + selectedService.Service.ID, outputObject);

            resultPromise.then(response => {
                log('updated service:', response.data);
                dispatch(serviceUpdateSuccess());
                dispatch(serviceFetchSuccess(response.data));
            }).catch(err => {
                Alert.error("Cannot update the service.");
                log("An update error occurred");
                dispatch(serviceUpdateError(err));
            })
        }
    }
}

function removeIOPort(isInput, removeIndex) {
    return function (dispatch, getState)  {
        const selectedService = getState().selectedService;
        dispatch(ioRemoveStart());
        let inputs = [];
        let outputs = [];

        if (isInput) {
            selectedService.Service.Input.map((input, currentIndex) => {
                if (currentIndex != removeIndex) {
                    inputs.push(input);
                }
            });
            outputs = selectedService.Service.Output;
        } else {

            selectedService.Service.Output.map((out, currentIndex) => {
                if (currentIndex != removeIndex) {
                    outputs.push(out);
                }
            });
            inputs = selectedService.Service.Input;
        }

        let outputObject = {
            'Created': selectedService.Service.Created,
            'Description': selectedService.Service.Description,
            'ID': selectedService.Service.ID,
            'ImageID': selectedService.Service.ImageID,
            'Input': inputs,
            'Name': selectedService.Service.Name,
            'Output': outputs,
            'RepoTag': selectedService.Service.RepoTag,
            'Size': selectedService.Service.Size,
            'Version': selectedService.Service.Version
        };

        dispatch(ioRemoveSuccess(outputObject));
        dispatch(serviceUpdateStart());
        const resultPromise = axios.put(apiNames.services + '/' + selectedService.Service.ID, outputObject);

        resultPromise.then(response => {
            log('updated service:', response.data);
            dispatch(serviceUpdateSuccess());
            dispatch(serviceFetchSuccess(response.data));
        }).catch(err => {
            Alert.error("Cannot update the service.");
            log("An update error occurred");
            dispatch(serviceUpdateError(err));
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

    serviceUpdateStart,
    serviceUpdateSuccess,
    serviceUpdateError,

    jobListFetchStart,
    jobListFetchSuccess,
    jobListFetchError,
    jobRemovalStart,
    jobRemovalSuccess,
    jobRemovalError,
    volumesFetchStart,
    volumesFetchSuccess,
    volumesFetchError,
    inspectVolumeStart,
    inspectVolumeSuccess,
    inspectVolumeError,
    consoleOutputFetchStart,
    consoleOutputFetchSuccess,
    consoleOutputFetchError,

    ioAddStart,
    ioAddSuccess,
    ioAddError,

    ioRemoveStart,
    ioRemoveSuccess,
    ioRemoveError,

    fetchJobs,
    removeJob,
    fetchServices,
    fetchService,
    handleUpdateService,
    fetchVolumes,
    inspectVolume,
    consoleOutputFetch,
    showErrorMessageWithTimeout,
    hideErrorMessage,
    fileUploadStart,
    fileUploadSuccess,
    fileUploadError,
    getNewUploadEndpoint,
    handleSubmitJob,
    addIOPort,
    removeIOPort,

};
