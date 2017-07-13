const urlRoot = ""; // window.location.origin;

export const apiNames = {
    apiinfo:        `${urlRoot}/api/`,
    user:           `${urlRoot}/api/user`,
    userTokens:     `${urlRoot}/api/user/tokens`,
    builds:         `${urlRoot}/api/builds`,
    services:       `${urlRoot}/api/services`,
    jobs:           `${urlRoot}/api/jobs`,
    volumes:        `${urlRoot}/api/volumes`,
};

export const wuiNames = {
    login:          `${urlRoot}/wui/login`,
    logout:         `${urlRoot}/wui/logout`,
};
