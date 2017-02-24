import React, {PropTypes} from 'react';
import { Row, Col, Grid, Table } from 'react-bootstrap';
import { toPairs } from '../utils/utils';
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import actions from '../actions/actions';

class ConsoleOutput extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        console.log(this.props);
        if (this.props.task.ServiceExecution) {
            return (
                <div style={{border: "1px solid black"}}>
                    {this.props.task.ServiceExecution.ConsoleOutput}
                </div>

            )
        } else {
            return (
                <div>Press Show button to see the console output</div>
            )
        }
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

export default connect(mapStateToProps, mapDispatchToProps)(ConsoleOutput);