'use strict';

import BrowseJobs from '../components/BrowseJobs';
import actions from '../actions/actions';

import {connect} from 'react-redux';


const mapStateToProps = (state) => {
    return {
        jobs: state.jobs
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        fetchJobs: () => {
                const action = actions.fetchJobs();
                dispatch(action);
        }
    };
};


const BrowseJobsContainer = connect(mapStateToProps, mapDispatchToProps)(
   BrowseJobs
);

export default BrowseJobsContainer;
