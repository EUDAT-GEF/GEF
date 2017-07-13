import React from 'react';
import PropTypes from 'prop-types';
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
import {fetchApiInfo} from './actions/actions';
import BrowseJobsContainer from './containers/JobsContainer';
import BuildServiceContainer from './containers/BuildServiceContainer';
import BrowseServicesContainer from './containers/ServicesContainer';
import {UserContainer, UserProfileContainer} from './containers/UserContainer';
import {FooterContainer} from './containers/FooterContainer';


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
        store.dispatch(fetchApiInfo());
    }

    render() {
        return (
            <Provider store={store}>
                <div style={{background: '#FFFFFF url("/images/img_header.png") no-repeat top right'}}>
                    <div style={{'paddingBottom': 70}}>
                        <Router history={history}>
                            <Route path='/' component={Frame}>
                                <Route path='user' component={UserProfileContainer} />
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
                    <FooterContainer />
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


const NotFound = () => (
    <div>This page is not found!</div>
);


render(<App />, document.getElementById('react') );
