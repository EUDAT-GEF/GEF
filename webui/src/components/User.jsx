import React, {PropTypes} from 'react';
import {NavItem, DropdownButton, MenuItem} from 'react-bootstrap'
import {Link} from 'react-router';
import axios from 'axios';
import apiNames from '../GefAPI';


const NoUser = (props) => (
    <NavItem onClick={()=>{window.location.href = "/login/"}} className="login">
        <i className="glyphicon glyphicon-log-in"/> Login
    </NavItem>
);


const ActiveUser = ({user}) => {
    const title = <span>
            <i className="glyphicon glyphicon-user"></i>
            {" "} {user.Name || user.Email}
        </span>;
    return <NavItem className="user">
        <DropdownButton title={title} style={{border:'none'}}>
            <MenuItem>
                <Link to="/user"> <i className="fa fa-info"></i> Profile </Link>
            </MenuItem>
            <MenuItem divider />
            <MenuItem>
                <a onClick={()=>window.location.href="/api/user/logout/"}> <i className="glyphicon glyphicon-log-out"></i> Logout </a>
            </MenuItem>
        </DropdownButton>
    </NavItem>
}


class User extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.props.fetchUser();
    }

    render() {
        if (!this.props.user || !this.props.user.Email) {
            return <NoUser/>;
        }
        return <ActiveUser user={this.props.user}/> ;
    }
};

export default User;
