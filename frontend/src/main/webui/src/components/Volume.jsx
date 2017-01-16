/**
 * Created by Alexandr Chernov on 16/12/16.
 */
import React, {PropTypes} from 'react';
import bows from 'bows';
import _ from 'lodash';
import {Row, Col, Table} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'
import actions from '../actions/actions'

const log = bows("Volume");

const styles = {
    volumeRowStyle: {
        fontWeight:700
    },
    volumeStyle: {
        height: "1em"
    }
};

const volumeRowStyle = {
    fontWeight:700
};

const Value = ({value}) => {
    if (typeof value === 'object') {
        _.toPairs(value).map(({k, v}) =>
            (
                 <div><dt>{k}</dt><dd>{v}</dd></div>
            ))
    } else {
        return <div>{value}</div>;
    }
};

const VolumeRow = ({tag, value, fileList}) => (
    <Row>
        <Col xs={12} sm={3} md={3} style={styles.volumeRowStyle}>{tag}</Col>
        <Col xs={12} sm={3} md={3} ><Value value={value}/></Col>
    </Row>
);

const VolumeFilesTable = ({fileList}) => (
    <Row>
        <Col xs={12} sm={3} md={3}>

        <Table striped bordered condensed hover>
            <thead>
              <tr>
                <th>#</th>
                <th>Name</th>
                <th>Size</th>
                <th>Modified</th>
              </tr>
            </thead>
            <tbody>
                {_.map(fileList, (fileListItem, index) => {
                    return <tr key={index}>
                           <td>{index+1}</td>
                           <td>{fileListItem.Name}</td>
                           <td>{fileListItem.Size}</td>
                           <td>@{fileListItem.Modified}</td>
                         </tr>;
                })}
            </tbody>
        </Table>
        </Col>
    </Row>
);



class Volume extends React.Component {
    constructor(props) {
        super(props);
    }

    handleInspect() {
        this.props.actions.inspectVolume(this.props.volume.ID)
    }

    componentDidMount() {
        this.props.actions.inspectVolume(); // send an empty list of files when a new box is drown
    }

    render() {

        let filesTable = null;
        if (this.props.selectedVolume.length > 0) {
            filesTable = <VolumeFilesTable fileList={this.props.selectedVolume}/>
        } else {
            filesTable = null;
        }
        return(
            <div style={{border: "1px solid black"}}>
                <div style={styles.volumeStyle}></div>
                <h4>Selected volume</h4>
                <VolumeRow tag="ID" value={this.props.volume.ID}/>
                <button type="submit" onClick={this.handleInspect.bind(this)}>Show list of files</button>
                {filesTable}
            </div>
        )
    }

}

Volume.propTypes = {
    volume: PropTypes.object.isRequired
};

export default Volume
