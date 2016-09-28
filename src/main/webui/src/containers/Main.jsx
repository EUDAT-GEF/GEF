'use strict';

import React, {PropTypes} from 'react';
import actions from '../actions/actions';
import {connect} from 'react-redux';
import {Grid, Row, Col} from 'react-bootstrap';
import ToolList from '../components/ToolList';
import BrowseJobsContainer from '../containers/BrowseJobsContainer';
import ExecuteService from '../components/ExecuteService';
import BuildServiceContainer from '../containers/BuildServiceContainer';


const pageNames = {
    browseJobs : 'Browse Jobs',
    buildService: 'Build a Service',
    executeService: 'Execute a Service'
};

const Main = ({currentPage, pageChange}) => {
    let page;
    switch (currentPage) {
        case pageNames.browseJobs:
            page =  <BrowseJobsContainer />;
            break;
        case pageNames.buildService:
            page = <BuildServiceContainer />;
            break;
        case pageNames.executeService:
            page = <ExecuteService />;
            break;
        default:
            page = <BrowseJobs />;
    };
    return (
        <Grid fluid={true}>
            <Row>
                <Col xs={12} sm={2} md={2}>
                    <ToolList currentPage={currentPage} onClick={pageChange}/>
                </Col>
                <Col xs={12} sm={10} md={10}>
                    {page}
                </Col>
            </Row>
        </Grid>
    )
};

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

export { pageNames };
export default MainContainer
