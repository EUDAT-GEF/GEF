import React, {PropTypes} from 'react';
import { Row, Col, Grid, Table, Button, Modal, OverlayTrigger } from 'react-bootstrap';
import { toPairs } from '../utils/utils';
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import actions from '../actions/actions';

class ModalWindow extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            showModal: false,
        };
    }

    close() {
        this.setState({ showModal: false });
    }

    open() {
        this.setState({ showModal: true });
    }

    render() {

        let modalBody = "";
        let modalTitle = "";
        if (this.props.modalBody) {
            modalBody = this.props.modalBody;
        }

        return (
            <div>
                <Button
                    bsStyle="primary"
                    bsSize="large"
                    onClick={this.open.bind(this)}
                >
                    Launch demo modal
                </Button>

                <Modal show={this.state.showModal} onHide={this.close.bind(this)}>
                    <Modal.Header closeButton>
                        <Modal.Title>{modalTitle}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <h4>Text in a modal</h4>
                        <p>Duis mollis, est non commodo luctus, nisi erat porttitor ligula.</p>

                        <hr />
                        {modalBody}
                    </Modal.Body>
                    <Modal.Footer>
                        <Button onClick={this.close.bind(this)}>Close</Button>
                    </Modal.Footer>
                </Modal>
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

export default connect(mapStateToProps, mapDispatchToProps)(ModalWindow);