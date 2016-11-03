'use strict';

import React, {PropTypes} from 'react';
import bows from 'bows';
import _ from 'lodash';
import {Row, Col} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'


const Volume = ({volume})  => {
    return (
        <div>
            Inspect the content of a volume, not implemented yet
        </div>

    )
};

Volume.propTypes = {
    volume: PropTypes.object.isRequired
};

export default Volume
