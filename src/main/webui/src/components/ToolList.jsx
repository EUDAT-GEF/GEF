'use strict';

import React, {PropTypes} from 'react';

import {ListGroup, ListGroupItem} from 'react-bootstrap'
import {pageNames} from '../containers/Main';
import _ from 'lodash';


const ToolList = ({currentPage, onClick}) => {
    return (
        <ListGroup>
            {_.values(pageNames).map((page) => (GefListItem(page, currentPage, onClick)))}
        </ListGroup>
    );
};

const GefListItem = (page, currentPage, onClick) => {
    if (page === currentPage) {
        return (<ListGroupItem href="#" active onClick={onClick(page)}> {page} </ListGroupItem>);
    }
    else {
        return (<ListGroupItem href="#" onClick={onClick(page)}> {page} </ListGroupItem>)

    }
};


ToolList.propTypes = {
    currentPage: PropTypes.string.isRequired,
    onClick: PropTypes.func.isRequired
};

export default ToolList;
