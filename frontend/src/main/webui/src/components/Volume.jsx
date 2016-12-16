/**
 * Created by Alexandr Chernov on 16/12/16.
 */
import React, {PropTypes} from 'react';
import bows from 'bows';
import _ from 'lodash';
import {Row, Col} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'

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

const VolumeRow = ({tag, value}) => (
    <Row>
        <Col xs={12} sm={3} md={3} style={styles.volumeRowStyle}>{tag}</Col>
        <Col xs={12} sm={3} md={3} ><Value value={value}/></Col>
    </Row>
);




class Volume extends React.Component {
    constructor(props) {
        super(props);
    }




    render() {
        return (
            <div style={{border: "1px solid black"}}>
                <div style={styles.volumeStyle}></div>
                <h4>Selected volume</h4>

                <VolumeRow tag="ID" value={this.props.volume.ID}/>
            </div>
        )

    }


}




Volume.propTypes = {
    volume: PropTypes.object.isRequired,
    fetchVolume: PropTypes.func.isRequired,
    handleSubmit: PropTypes.func.isRequired,
};

export default Volume
