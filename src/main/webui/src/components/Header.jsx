'use strict';

import React from 'react';
import {Grid, Row, Col} from 'react-bootstrap';
import Radium from 'radium';

const style = {
    padding: '1% 0',
    background: '#FFFFFF url("images/img_header.png") no-repeat top right'
};

const Header = () => (
    <Grid style={style}>
        <Row>
            <Col xs={12} md={3} sm={3}>
                <div>
                    <a href="/gef"><img width="232" height="128" src="images/logo.png" alt=""/></a>
                </div>
            </Col>
            <Col xs={12} md={9} sm={9}>
            </Col>
        </Row>

    </Grid>
);

export default Radium(Header);
