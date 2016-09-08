'use strict';

import React, {PropTypes} from 'react';
import Files from './Files';
import {Row, Col, Button, Glyphicon} from 'react-bootstrap'

const BuildService = () => (
    <div>
        <h3>Build Service</h3>
        <h4>Please select and upload the Dockerfile, together with other files which are part of the container</h4>
        <Files/>
        <Row>
            <Col md={4} mdOffset={8}> <Button type='submit' bsStyle='primary' style={{width: '100%'}}> <Glyphicon glyph='upload'/> Build Image</Button> </Col>
        </Row>
    </div>
);

export default BuildService;
