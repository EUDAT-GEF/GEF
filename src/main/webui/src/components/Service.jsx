/**
 * Created by wqiu on 17/08/16.
 */
'use strict';

import React, {PropTypes} from 'react';
import axios from 'axios';
import apiNames from '../utils/GefAPI';
import bows from 'bows';
import {Row, Col, Grid} from 'react-bootstrap';
// this is a detailed view of a service, user will be able to execute service in this view


const tagValueRow  = ({tag, value}) => (
    <Row>
           <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>{tag}</Col>
           <Col xs={12} sm={9} md={9} >{value}</Col>
    </Row>
);

class Service extends React.Component {

    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.props.fetchService(this.props.service.ID);
    }

    render() {
        return (
            <div>
                Execute a service, not implemented yet
            </div>

        )
    }
}


Service.propTypes = {
    service: PropTypes.object.isRequired,
    fetchService: PropTypes.func.isRequired
};

export default Service

