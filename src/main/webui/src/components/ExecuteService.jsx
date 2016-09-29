'use strict';

import React, {PropTypes} from 'react';
import axios from 'axios';
import apiNames from '../utils/GefAPI';
import bows from 'bows';
// this is a detailed view of a service, user will be able to execute service in this view

const ExecuteService = ({service})  => {
    return (
        <div>
            run a service
        </div>

    )
};

ExecuteService.propTypes = {
    service: PropTypes.object.isRequired
};

export default ExecuteService