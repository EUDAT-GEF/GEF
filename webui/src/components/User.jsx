import React from 'react';
import PropTypes from 'prop-types';
import {NavItem, DropdownButton, MenuItem, Row, Col, Grid, Table, Button, Modal, OverlayTrigger } from 'react-bootstrap';
import {Link} from 'react-router';
import axios from 'axios';
import {apiNames, wuiNames} from '../GefAPI';
import {Field, reduxForm, initialize} from 'redux-form';
import moment from 'moment';


const NoUser = () => (
    <NavItem onClick={()=>{window.location.href = wuiNames.login}} className="login">
        <i className="glyphicon glyphicon-log-in"/> Login
    </NavItem>
);


const ActiveUser = ({user, isSuperAdmin}) => {
    const title = <span>
            <i className="glyphicon glyphicon-user"></i>
            {" "} {user.Name || user.Email}
        </span>;
    const style = {border:'none'};
    if (isSuperAdmin) {
        style.backgroundColor ='orange';
    }
    return <NavItem className="user">
        <DropdownButton id="usermenu" title={title} style={style}>
            <MenuItem>
                <Link to="/user"> <i className="glyphicon glyphicon-info-sign"></i> Profile </Link>
            </MenuItem>
            { isSuperAdmin ?
                <MenuItem>
                    <Link to="/roles"> <i className="glyphicon glyphicon-tasks"></i> Manage Roles </Link>
                </MenuItem>:false
            }
            <MenuItem divider />
            <MenuItem>
                <a onClick={()=>window.location.href = wuiNames.logout}>
                    <i className="glyphicon glyphicon-log-out"></i> Logout </a>
            </MenuItem>
        </DropdownButton>
    </NavItem>
}


export class User extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.props.fetchUser();
    }

    render() {
        const data = this.props.user || {};
        const user = data.User || {};
        if (!user.Email) {
            return <NoUser/>;
        }
        return <ActiveUser user={user} isSuperAdmin={data.IsSuperAdmin}/> ;
    }
};


export const UserProfile = React.createClass({
    componentDidMount() {
        this.props.fetchTokens();
    },

    renderNoUser() {
        return (
            <div>
                <h1>User Profile</h1>
                <div className="container-fluid">
                    <div className="row">
                        No authenticated user found, please login.
                    </div>
                </div>
            </div>
        );
    },

    renderRole(r) {
        return (
            <li key={r.ID}>
                {r.Description}
                {r.CommunityName ? <span> {" (" + r.CommunityName + ")" } </span> : false }
            </li>
        );
    },

    render() {
        const data = this.props.user || {};
        const user = data.User || {};
        const roles = data.Roles || [];
        const communityMap = data.CommunityMap || {};
        if (!user.Email) {
            return this.renderNoUser();
        }

        return (
            <div>
                <h1>User Profile</h1>
                <div className="container-fluid">
                    <div className="row">
                        <p><span style={{fontWeight:'bold'}}>Name:</span> {user.Name}</p>
                        <p><span style={{fontWeight:'bold'}}>Email:</span> {user.Email || "-"}</p>
                    </div>
                    <div className="row">
                        <h3>Roles</h3>
                        {data.IsSuperAdmin ?
                            <li style={{fontWeight:'bold', color:"red"}}>
                                SuperAdministrator
                            </li> : false }
                        { roles.map(this.renderRole) }
                    </div>
                    <div className="row">
                        <h3>API Tokens</h3>
                        <TokenList tokens={this.props.tokens}
                            submitNewAccessToken={this.props.submitNewAccessToken}
                            deleteAccessToken={this.props.deleteAccessToken} />
                    </div>
                </div>
            </div>
        );
    }
});


const TokenList = React.createClass({
    getInitialState() {
        return {
            newtoken: null,
        };
    },

    submitNewToken(event) {
        const {tokenName} = event;
        this.props.submitNewAccessToken(tokenName, data => {
            this.setState({newtoken: data.Token});
        });
    },

    deleteToken(tokenID) {
        this.props.deleteAccessToken(tokenID);
    },

    renderTokenLine(token) {
        return (
            <li className="list-group-item" key={token.ID}>
                {token.Name}
                <a className="btn btn-xs btn-warning" style={{float:'right'}}
                    onClick={()=>this.deleteToken(token.ID)}>Delete</a>
                <code style={{margin:'0 2em', float:'right'}}>{token.Secret}</code>
            </li>
        );
    },

    render() {
        const tokens = this.props.tokens;
        const newtoken = this.state.newtoken;
        return (
            <div>
                {tokens && tokens.length ?
                    <div>
                        <p>Active tokens:</p>
                            <ul className="list-group">
                                { tokens.map(this.renderTokenLine) }
                            </ul>
                    </div> : <div><p> You have no active tokens </p></div>
                }

                {newtoken ?
                    <div className="alert alert-success">
                        <h4>A new access token has just been created:</h4>
                        <dl className="dl-horizontal">
                            <dt>Name</dt>
                            <dd>{newtoken.Name}</dd>
                            <dt>Access token</dt>
                            <dd className="access-token">{newtoken.Secret}</dd>
                            <dt>Expiration date</dt>
                            <dd>{moment(newtoken.Expire).format('ll')}</dd>
                        </dl>
                        <p className="alert alert-warning">
                            Please copy the personal access token now. You won't see it again!
                        </p>
                    </div> : false
                }

                <div style={{marginTop:'2em'}}>
                    <h4>Create new token</h4>
                    <NewTokenForm onSubmit={this.submitNewToken}/>
                </div>
            </div>
        );
    }
});


const NewToken = ({handleSubmit}) => {
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
            <Row>
                <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>
                    Token name
                </Col>
                <Col xs={12} sm={9} md={9} >
                    <div className="input-group">
                        <Field name="tokenName" component="input" type="text" placeholder="Token Name"
                               style={inputStyle} className="form-control"/>
                        <span className="input-group-btn">
                            <button type="submit" className="btn btn-default">
                                <span className="glyphicon glyphicon-play" aria-hidden="true"></span> Submit
                            </button>
                        </span>
                    </div>
                </Col>
            </Row>
        </form>
    )
};
const NewTokenForm = reduxForm({form: 'NewToken'})(NewToken);
