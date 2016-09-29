'use strict';

import Workflows from '../components/Workflows';
import actions from '../actions/actions';

import {connect} from 'react-redux';


const mapStateToProps = (state) => {
    return {
        workflows: state.workflows
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        fetchWorkflows: (dispatch)  => {return {}}
    }

};


const WorkflowsContainer = connect(mapStateToProps, mapDispatchToProps)(
    Workflows
);

export default WorkflowsContainer;
