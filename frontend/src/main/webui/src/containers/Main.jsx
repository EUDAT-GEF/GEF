import React, {PropTypes} from 'react';
import actions from '../actions/actions';
import {connect} from 'react-redux';
import {Grid, Row, Col} from 'react-bootstrap';
import BrowseJobsContainer from './JobsContainer';
import BuildServiceContainer from '../containers/BuildServiceContainer';
import {Router, Route, hashHistory, Link } from 'react-router';
import {ListGroup, ListGroupItem} from 'react-bootstrap'
import {LinkContainer} from 'react-router-bootstrap'


const ToolNav = () => (
    <div>
        <ListGroup>
            <LinkContainer to='/jobs' >
                <ListGroupItem> Browse Jobs </ListGroupItem>
            </LinkContainer>
            <LinkContainer to='/services' >
                <ListGroupItem> Browse Services </ListGroupItem>
            </LinkContainer>
        </ListGroup>
        <ListGroup>
            <LinkContainer to='/buildImage' >
                <ListGroupItem> Build a Service </ListGroupItem>
            </LinkContainer>
        </ListGroup>
    </div>
);

const Main = (props) => (
    <Grid fluid={true}>
        <Row>
            <Col xs={12} sm={2} md={2}>
                <ToolNav></ToolNav>
            </Col>
            <Col xs={12} sm={10} md={10}>
                {props.children}
            </Col>
        </Row>
    </Grid>
);

const mapStateToProps = (state) => {
    return {
        currentPage: state.currentPage
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        pageChange: (pageName) => {
            return () => {
                const action = actions.pageChange(pageName);
                dispatch(action);
            };
        }
    };
};


const MainContainer = connect(mapStateToProps, mapDispatchToProps)(
    Main
);

export default MainContainer
