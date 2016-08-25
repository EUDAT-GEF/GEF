'use strict';

import React, {PropTypes} from 'react';
import actions from '../actions/actions';
import {connect} from 'react-redux';
import {Grid, Row, Col} from 'react-bootstrap';
import ToolList from '../components/ToolList';
import BrowseJobs from '../components/BrowseJobs';
import ExecuteService from '../components/ExecuteService';
import BuildService from '../components/BuildService';


const pageNames = {
    browseJobs : 'Browse Jobs',
    buildService: 'Build a Service',
    executeService: 'Execute a Service'
};

const Main = ({currentPage, pageChange}) => {
    let page;
    switch (currentPage) {
        case pageNames.browseJobs:
            page =  <BrowseJobs />;
            break;
        case pageNames.buildService:
            page = <BuildService />;
            break;
        case pageNames.executeService:
            page = <ExecuteService />;
            break;
        default:
            page = <BrowseJobs />;
    };
    return (
        <Grid>
            <Row>
                <Col xs={12} sm={2} md={2}>
                    <ToolList currentPage={currentPage} onClick={pageChange}/>
                </Col>
                <Col xs={12} sm={2} md={2}>
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
