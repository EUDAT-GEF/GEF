import Jobs from '../components/Jobs';
import {fetchJobs, fetchServices} from '../actions/actions';
import {connect} from 'react-redux';

const mapStateToProps = (state) => {
    return {
        jobs: state.jobs,
        services: state.services
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        fetchJobs: () => {
            const action = fetchJobs();
            dispatch(action);
        },
        fetchServices: () => {
            const action = fetchServices();
            dispatch(action);
        }
    };
};

const JobsContainer = connect(mapStateToProps, mapDispatchToProps)(
   Jobs
);

export default JobsContainer;