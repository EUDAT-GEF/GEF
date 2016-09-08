'use strict';

import React, {PropTypes} from 'react';
import {Row, InputGroup} from 'react-bootstrap';
import DropzoneComponent from 'react-dropzone-component';


const FileItemRow = ({file, handleRemove}) => (
    <Row key={f.name}>
        <Col md={12}>
            <div className="input-group">
                <input type="text" className="input-large" readOnly="readonly"
                       style={{width:'100%', lineHeight:'26px', paddingLeft:10}} value={f.name}/>
                <span className="input-group-btn">
							<button type="button" className="btn btn-warning btn-sm" onClick={handleRemove(file)}>
								<i className="glyphicon glyphicon-remove"/>
							</button>
						</span>
            </div>
        </Col>
    </Row>
);

FileItemRow.propTypes = {
    file: PropTypes.object.isRequired,
    handleRemove: PropTypes.func.isRequired
};

const BottomRow = ({files, handleAdd}) => (
    <Row style={{margin: '5px 0px'}}>
        <Col md={3}>
            <FileAddButton> </FileAddButton>
        </Col>
    </Row>
);

BottomRow.propTypes = {
    files: PropTypes.array.isRequired,
    handleAdd: PropTypes.func.isRequired
};

const fileAddButtonStyle = {
    noshow: {width:0, height:0, margin:0, border:'none'}
};

class FileAddButton extends React.Component {
    constructor(props) {
        super(props);
    }

    doBrowse(event) {
        this.element.click();
    }

    doAddFiles(event) {
        this.props.addFile(event.target.files);
    }
}

FileAddButton.propTypes = {
    addFile: PropTypes.func.isRequired
};







