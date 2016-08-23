'use strict';

import React from 'react';
import Radium from 'radium';
import {PropTypes} from 'react';
import {Grid, Row, Col} from 'react-bootstrap';

const styles = {
    position: 'absolute',
    bottom: 0,
    width: '100%',
    height: 70   /* Height of the footer */
};

const Footer = ({version}) => (
    <Grid style={styles}>
        <Row>
            <Col xs={12} md={6} sm={6}>
                <p>	<img width="45" height="31" src="images/flag-ce.jpg" style={{float:'left', marginRight:10}}/>
                    EUDAT receives funding from the European Union’s Horizon 2020 research
                    and innovation programme under grant agreement No. 654065.&nbsp;
                    <a href="#">Legal Notice</a>.
                </p>
            </Col>
            <Col xs={12} sm={6} md={6}>
                <ul className="list-inline pull-right" style={{marginLeft:20}}>
                    <li><span style={{color:'#173b93', fontWeight:'500'}}> GEF v.{version}</span></li>
                </ul>
                <ul className="list-inline pull-right">
                    <li><a target="_blank" href="http://eudat.eu/what-eudat">About EUDAT</a></li>
                    <li><a href="https://github.com/GEFx">Go to GitHub</a></li>
                    <li><a href="mailto:emanuel.dima@uni-tuebingen.de">Contact</a></li>
                </ul>
            </Col>
        </Row>
    </Grid>
)

Footer.propTypes = {
    version: PropTypes.string.isRequired
};

export default Radium(Footer);