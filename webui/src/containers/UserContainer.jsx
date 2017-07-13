import {User, UserProfile} from '../components/User';
import {fetchUser, fetchTokens, submitNewAccessToken, deleteAccessToken} from '../actions/actions';
import {connect} from 'react-redux';

const mapStateToProps = (state) => ({
    user: state.user,
    tokens: state.tokens,
});

const mapDispatchToProps = (dispatch) => ({
    fetchUser: () => dispatch(fetchUser()),
    fetchTokens: () => dispatch(fetchTokens()),
    submitNewAccessToken: (tokenName, successFn) => dispatch(submitNewAccessToken(tokenName, successFn)),
    deleteAccessToken: (tokenID) => dispatch(deleteAccessToken(tokenID)),
});

export const UserContainer = connect(mapStateToProps, mapDispatchToProps)(User);
export const UserProfileContainer = connect(mapStateToProps, mapDispatchToProps)(UserProfile);
