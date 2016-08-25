'use strict';

import axios from 'axios';



const apiNames = {
    datasets: "/gef/api/datasets",
    builds:   "/gef/api/builds",
    services: "/gef/api/images",
    jobs: "/gef/api/jobs",
};

const getJobs = () => {
    return axios.get(apiNames.jobs);
};

const getBuilds = () => {
    return axios.get(apiNames.builds);
};

const getServices = () => {
    return axios.get(apiNames.services);
};

