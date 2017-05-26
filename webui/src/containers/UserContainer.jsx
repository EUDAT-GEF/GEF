import User from '../components/User';
import {fetchUser} from '../actions/actions';
import {connect} from 'react-redux';

const mapStateToProps = (state) => ({user: state.user});

const mapDispatchToProps = (dispatch) => ({
    fetchUser: () => dispatch(fetchUser()),
});

const UserContainer = connect(mapStateToProps, mapDispatchToProps)(User);

export default UserContainer;
