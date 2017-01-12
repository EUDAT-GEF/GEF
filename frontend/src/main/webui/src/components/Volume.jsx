/**
 * Created by Alexandr Chernov on 16/12/16.
 */
import React, {PropTypes} from 'react';
import bows from 'bows';
import _ from 'lodash';
import {Row, Col} from 'react-bootstrap';
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


const FormattedList = ({formattedList}) => {
    console.log(formattedList);
    console.log(formattedList.length);
    console.log(formattedList[0]);


 {_.map(formattedList, (fileListItem) => {
                        console.log(fileListItem.Name);


                    })}


    var volumeFiles = [];
    for (var i = 0; i < formattedList.length; i++) {
        volumeFiles.push("<li>" + formattedList[i].Name + "</li>");
    }
    if (volumeFiles.length > 0) {
    var lst = volumeFiles.join("\n");
        return <div><ul>{lst}</ul></div>;
    } else {
        return <div>Files not found</div>;
    }




};

const VolumeRow = ({tag, value, fileList}) => (
    <Row>
        <Col xs={12} sm={3} md={3} style={styles.volumeRowStyle}>{tag}</Col>
        <Col xs={12} sm={3} md={3} ><Value value={value}/><br/>
        <ul>
        {_.map(fileList, (fileListItem) => {
            console.log(fileListItem.Name);
            return <li>{fileListItem.Name}</li>;

        })}
        </ul>

        </Col>
    </Row>
);




class Volume extends React.Component {
    constructor(props) {
        super(props);
    }

    handleInspect() {
        console.log("checking");
        this.props.actions.inspectVolume(this.props.volume.ID)
    }


    render() {

        console.log(this.props);
        console.log(this.props.selectedVolume);
        var someVar = "Some text";

        return (

            <div style={{border: "1px solid black"}} onClick={this.handleInspect.bind(this)}>
                <div style={styles.volumeStyle}></div>
                <h4>Selected volume</h4>

                <VolumeRow tag="ID" value={this.props.volume.ID} fileList={this.props.selectedVolume}/>



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
