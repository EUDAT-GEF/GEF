/**
 * Created by wqiu on 17/08/16.
 */
'use strict';

import React, {PropTypes} from 'react';
import axios from 'axios';
import apiNames from '../utils/GefAPI';
import bows from 'bows';
// this is a detailed view of a service, user will be able to execute service in this view

const Service = ({service})  => {
    return (
        <div>
            Execute a service, not implemented yet
        </div>

    )
};

Service.propTypes = {
    service: PropTypes.object.isRequired
};

export default Service

