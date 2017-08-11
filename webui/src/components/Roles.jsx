import React from 'react';
import PropTypes from 'prop-types';
import {NavItem, DropdownButton, MenuItem, Row, Col, Grid, Table, Button, Modal, OverlayTrigger } from 'react-bootstrap';
import {Link} from 'react-router';
import axios from 'axios';
import {apiNames, wuiNames} from '../GefAPI';
import {Field, reduxForm, initialize} from 'redux-form';
import moment from 'moment';
import { toPairs } from '../utils/utils';


export const Roles = React.createClass({
    componentDidMount() {
        this.props.fetchRoles();
    },

    renderRole([rID, r]) {
        return (
            <Role
                key={r.ID}
                role={r}
                fetchRoleUsers={this.props.fetchRoleUsers}
                newRoleUser={this.props.newRoleUser}
                deleteRoleUser={this.props.deleteRoleUser}
                />
        );
    },

    render() {
        const roles = this.props.roles || {};
        return (
            <div>
                <h1>Role Management</h1>
                <div className="container-fluid">
                    <ul className="list-unstyled">
                        { toPairs(roles).map(this.renderRole) }
                    </ul>
                </div>
            </div>
        );
    }
});

const Role = React.createClass({
    getInitialState() {
        return {
            open: true,
        }
    },

    componentDidMount() {
        this.props.fetchRoleUsers(this.props.role.ID);
    },

    toggle(e) {
        const open = !this.state.open;
        if (open) {
            this.props.fetchRoleUsers(this.props.role.ID);
        }
        this.setState({open});
    },

    newRole({userEmail}) {
        this.props.newRoleUser(this.props.role.ID, userEmail);
    },

    renderUser(user, idx) {
        return (
            <li className="list-group-item" key={idx}>
                <span className="glyphicon glyphicon-user" aria-hidden="true"/>&nbsp;
                {user.Name}&nbsp;
                {user.Email}
                <a className="btn btn-xs btn-warning" style={{float:'right'}}
                    onClick={()=>this.props.deleteRoleUser(this.props.role.ID, user.ID)}>
                    <span className="glyphicon glyphicon-remove" aria-hidden="true"/>&nbsp;
                </a>
            </li>
        );
    },

    render() {
        const role = this.props.role;
        const users = role.users || [];
        const NewRoleForm = newRoleForm('NewRoleUser'+role.ID);
        return (
            <li className="row" key={role.ID} style={{marginTop:'2em'}}>
                <h4 onClick={this.toggle}>
                    { this.state.open ?
                        <span className="glyphicon glyphicon-chevron-down" aria-hidden="true"/> :
                        <span className="glyphicon glyphicon-chevron-right" aria-hidden="true"/>
                    }&nbsp;
                    { role.Description }
                    { role.CommunityName ?
                        <span> {" (" + role.CommunityName + ")" } </span>
                        : false }
                </h4>
                { this.state.open ?
                    <div style={{marginLeft:'2em'}}>
                        { users.length ?
                            <ul className="list-unstyled">
                                { users.map(this.renderUser) }
                            </ul> :
                            <span>No users</span> }
                        <div style={{marginTop:'1em'}}>
                            <NewRoleForm onSubmit={this.newRole}/>
                        </div>
                    </div>
                    : false }
            </li>
        );
    }
});

function newRoleForm(formDataName) {
    const NewRole = ({handleSubmit}) => {
        const inputStyle = {
            height: '34px',
            padding: '6px 12px',
            fontSize: '14px',
            lineHeight: '1.42857143',
            color: '#555',
            backgroundColor: '#fff',
            backgroundImage: 'none',
            border: '1px solid #ccc',
            borderRadius: '4px',
        }

        return (
            <form onSubmit={handleSubmit}>
                <div className="input-group">
                    <span className="input-group-addon">Add user (by email) </span>
                    <Field name="userEmail" component="input" type="text" placeholder="user@example.com"
                           style={inputStyle} className="form-control"/>
                    <span className="input-group-btn">
                        <button type="submit" className="btn btn-default">
                            <span className="glyphicon glyphicon-play" aria-hidden="true"></span> Submit
                        </button>
                    </span>
                </div>
            </form>
        )
    };
    return reduxForm({form: formDataName}) (NewRole);
}
