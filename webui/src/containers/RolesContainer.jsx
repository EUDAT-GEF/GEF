import {Roles} from '../components/Roles';
import {fetchRoles, fetchRoleUsers, newRoleUser, deleteRoleUser} from '../actions/actions';
import {connect} from 'react-redux';

const mapStateToProps = (state) => ({
    roles: state.roles,
    roleUsers: state.roleUsers,
});

const mapDispatchToProps = (dispatch) => ({
    fetchRoles: () => dispatch(fetchRoles()),
    fetchRoleUsers: (roleID) => dispatch(fetchRoleUsers(roleID)),
    newRoleUser: (roleID, user) => dispatch(newRoleUser(roleID, user)),
    deleteRoleUser: (roleID, user) => dispatch(deleteRoleUser(roleID, user)),
});

export const RolesContainer = connect(mapStateToProps, mapDispatchToProps)(Roles);
