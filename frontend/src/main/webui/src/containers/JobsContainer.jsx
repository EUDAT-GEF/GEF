import Jobs from '../components/Jobs';
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


const JobsContainer = connect(mapStateToProps, mapDispatchToProps)(
   Jobs
);

export default JobsContainer;
