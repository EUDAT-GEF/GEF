const urlRoot = ""; // window.location.origin;

const apiNames = {
    user:           `${urlRoot}/api/user`,
    builds:         `${urlRoot}/api/builds`,
    services:       `${urlRoot}/api/services`,
    jobs:           `${urlRoot}/api/jobs`,
    volumes:        `${urlRoot}/api/volumes`,
};

export default apiNames;
