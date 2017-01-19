const urlRoot = ""; // window.location.origin;

const apiNames = {
    datasets:       `${urlRoot}/api/datasets`,
    buildImages:    `${urlRoot}/api/buildImages`,
    buildVolumes:   `${urlRoot}/api/buildVolumes`,
    volumes:        `${urlRoot}/api/volumes`,
    services:       `${urlRoot}/api/images`,
    jobs:           `${urlRoot}/api/jobs`,
};

export default apiNames;
