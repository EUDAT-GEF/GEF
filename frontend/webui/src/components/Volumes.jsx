import React, {PropTypes} from 'react';
import bows from 'bows';
import {Row, Col} from 'react-bootstrap';
import {LinkContainer} from 'react-router-bootstrap'
import Volume from './Volume'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import actions from '../actions/actions'


const log = bows('Volumes');

const VolumeRow = ({volume}) => (
    <LinkContainer to={`/volumes/${volume.ID}`}>
        <Row>
            <Col xs={12} sm={4} md={4}><i className="glyphicon glyphicon-transfer"/>{volume.ID}</Col>
            <Col xs={12} sm={4} md={4}>{volume.Mountpoint}</Col>
        </Row>
    </LinkContainer>
);

const Header = () => (
    <div className="row table-head">
        <div className="col-xs-12 col-sm-4">ID</div>
        <div className="col-xs-12 col-sm-4">Internal Location</div>
    </div>
);

class Volumes extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.props.fetchVolumes();
    }

    render() {
        return (
            <div>
                <h3>Browse Volumes</h3>
                <h4>All Volumes</h4>
                <Header/>
                {this.props.volumes.map((volume) => {
                    //console.log(this.props.selectedVolume);
                    if(volume.ID === this.props.params.id)
                        return <Volume actions={this.props.actions} key={volume.ID} volume={volume} selectedVolume={this.props.selectedVolume}/>;
                    else
                        return <VolumeRow actions={this.props.actions} key={volume.ID} volume={volume} selectedVolume={this.props.selectedVolume}/>;

                })}
            </div>
        );
    }

}

function mapStateToProps(state) {
    return state
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators(actions, dispatch)
    }
}

Volumes.propTypes = {
    fetchVolumes: PropTypes.func.isRequired,
    volumes: PropTypes.array.isRequired,
    volumeID: PropTypes.string

};

export default connect(mapStateToProps, mapDispatchToProps)(Volumes);








