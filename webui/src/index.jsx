import React, {PropTypes} from 'react';
import {render} from 'react-dom';
import {Provider, connect} from 'react-redux';
import {Router, Route, browserHistory, hashHistory, Link } from 'react-router';
import {syncHistoryWithStore, routerMiddleware} from 'react-router-redux'
import {compose, applyMiddleware, createStore} from 'redux';
import createLogger from 'redux-logger';
import thunkMiddleware from 'redux-thunk';

require('react-s-alert/dist/s-alert-default.css');
require('react-s-alert/dist/s-alert-css-effects/slide.css');

import bows from 'bows';
import {Grid, Row, Col, ListGroup, ListGroupItem, Navbar, Nav, NavItem } from 'react-bootstrap'
import {LinkContainer} from 'react-router-bootstrap'
import Alert from 'react-s-alert';

import rootReducers from './actions/reducers';
import BrowseJobsContainer from './containers/JobsContainer';
import BuildServiceContainer from './containers/BuildServiceContainer';
import BrowseServicesContainer from './containers/ServicesContainer';
import UserContainer from './containers/UserContainer';


const log = bows('app');

const middleware = compose(
    applyMiddleware(thunkMiddleware,
                    routerMiddleware(browserHistory),
                    createLogger()),
    window.devToolsExtension || (f=>f)
);
const store = createStore(rootReducers, middleware);
const history = syncHistoryWithStore(browserHistory, store);


class App extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <Provider store={store}>
                <div style={{background: '#FFFFFF url("/images/img_header.png") no-repeat top right'}}>
                    <div style={{'paddingBottom': 70}}>
                        <Router history={history}>
                            <Route path='/' component={Frame}>
                                <Route path='builds' component={BuildServiceContainer} />
                                <Route path='services' component={BrowseServicesContainer} >
                                    <Route path=':id' />
                                </Route>
                                <Route path='jobs' component={BrowseJobsContainer} >
                                    <Route path=':id' />
                                </Route>
                                <Route path='*' component={NotFound} />
                            </Route>
                        </Router>
                    </div>
                    <Alert stack={{limit:5}}/>
                    <Footer version="0.4.0" />
                </div>
            </Provider>
        );
    }
}

const Frame = (props) => (
    <Grid fluid={true}>
        <Row>
            <Col xs={12} sm={3} md={3} >
                <div>
                    <a href="/"><img width="232" height="128" src="/images/logo.png" alt=""/></a>
                </div>
            </Col>
            <Col xs={12} sm={8} md={8}>
                <Navbar fluid={true} collapseOnSelect style={{marginTop:'80px'}}>
                    <Navbar.Header>
                        <Navbar.Toggle />
                    </Navbar.Header>
                    <Navbar.Collapse>
                        <Nav>
                            <LinkContainer to='/builds' >
                                <NavItem> Build </NavItem>
                            </LinkContainer>
                            <LinkContainer to='/services' >
                                <NavItem> Services </NavItem>
                            </LinkContainer>
                            <LinkContainer to='/jobs' >
                                <NavItem> Jobs </NavItem>
                            </LinkContainer>
                        </Nav>
                        <Nav pullRight>
                            <UserContainer/>
                        </Nav>
                    </Navbar.Collapse>
                </Navbar>
            </Col>
        </Row>
        <Row>
            <Col xs={12}>
                <div style={{borderBottom:'1px solid #ddd'}}/>
            </Col>
        </Row>
        <Row>
            <Col xs={12} sm={12} md={1}/>
            <Col xs={12} sm={12} md={10}>
                {props.children}
            </Col>
        </Row>
    </Grid>
);


const footerStyle = {
    position: 'fixed',
    bottom: 0,
    width: '100%',
    height: 70,   /* Height of the footer */
    background: '#F7F3E9 url("/images/color-line.jpg") repeat-x top left',
    padding: '20px 10px 0px 10px',
    fontSize: 12
};
const Footer = ({version}) => (
    <Grid style={footerStyle}>
        <Row>
            <Col xs={12} md={6} sm={6}>
                <p> <img width="45" height="31" src="/images/flag-ce.jpg" style={{float:'left', marginRight:10}}/>
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
                    <li><a href="https://github.com/EUDAT-GEF">Go to GitHub</a></li>
                    <li><a href="mailto:emanuel.dima@uni-tuebingen.de">Contact</a></li>
                </ul>
            </Col>
        </Row>
    </Grid>
);
Footer.propTypes = {
    version: PropTypes.string.isRequired
};


const NotFound = () => (
    <div>This page is not found!</div>
);


render(<App />, document.getElementById('react') );
