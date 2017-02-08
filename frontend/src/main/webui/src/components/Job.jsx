import React, {PropTypes} from 'react';
import { Row, Col, Grid, Table } from 'react-bootstrap';
import { toPairs } from '../utils/utils';
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import actions from '../actions/actions';

const Value = ({value}) => {
    if (typeof value === 'object') {
        toPairs(value).map(({k, v}) =>
            (
                 <div><dt>{k}</dt><dd>{v}</dd></div>
            ))
    } else {
        return <div>{value}</div>;
    }
};

const JobRow = ({tag, value, style}) => (
    <Row style={style}>
        <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>{tag}</Col>
        <Col xs={12} sm={3} md={3} ><Value value={value}/></Col>
    </Row>
);



///var iconClass = !file.isdir ? "glyphicon-file"
//                        : file.children == undefined ? "glyphicon-folder-close"
//                        : "glyphicon-folder-open";
//


// const indentStyle = {paddingLeft: (3*file.indent)+'em'};
//      const handlerStyle = {width:20, background:'none', border:'none', fontSize:20, padding:0};

const volumeFile = ({file, index, iconClass, indentStyle}) => (
    <li className="row file" key={file.path} style={{lineHeight:2}} onClick={this.handleFileClick.bind(this, file, index)}>
        <div className="col-sm-6">
            <span style={indentStyle}/>
            { file.isdir ?
                file.children == undefined ?
                    <button style={{width:20, background:'none', border:'none', fontSize:20, padding:0}}>+</button> :
                    <button style={{width:20, background:'none', border:'none', fontSize:20, padding:0}}>-</button> :
                <input type="checkbox" style={{width:20}} checked={file.selected}/> }
            <span className={"glyphicon "+iconClass} aria-hidden={true} /> {file.name}
        </div>
        <div className="col-sm-3">{file.size}</div>
        <div className="col-sm-3">{file.date}</div>
    </li>
);

const VolumeFilesTable = ({fileList}) => (
    <div style={{margin:'1em'}}>
        <ol className="list-unstyled fileList" style={{textAlign:'left', minHeight:'30em'}}>
            <li className="heading row" style={{padding:'0.5em 0'}}>
                <div className="col-sm-6" style={{fontWeight:'bold'}}>File Name</div>
                <div className="col-sm-3" style={{fontWeight:'bold'}}>Size</div>
                <div className="col-sm-3" style={{fontWeight:'bold'}}>Date</div>
            </li>

            {fileList.map((fileListItem, index) => {
                console.log(index);
                let indentStyle = {paddingLeft: (3*1)+'em'};
                let iconClass = "glyphicon-file";
                if (fileListItem.isFolder == true) {
                    iconClass = "glyphicon-folder-close";
                }
                return <li className="row file" key={index} style={{lineHeight:2}}>
                   <div className="col-sm-6">
                       <span style={indentStyle}/>
                       <span className={"glyphicon "+iconClass} aria-hidden={true} /> {fileListItem.name}
                   </div>
                   <div className="col-sm-3">{fileListItem.size}</div>
                   <div className="col-sm-3">{fileListItem.modified}</div>
               </li>
            })}
        </ol>
    </div>
);

class Job extends React.Component {
    constructor(props) {
        super(props);
    }

    handleInspectInputVolume() {
        this.props.actions.inspectVolume(this.props.job.InputVolume)
    }

    handleInspectOutputVolume() {
        this.props.actions.inspectVolume(this.props.job.OutputVolume)
    }

    componentDidMount() {
        this.props.actions.inspectVolume(); // send an empty list of files when a new box is drown
    }

    render() {
        console.log(this.props);
        let job = this.props.job;
        let service = this.props.service;
        let title = this.props.title;

        console.log(this.props.selectedVolume);
        let filesTable = null;
        if (this.props.selectedVolume.length > 0) {
            //filesTable = this.props.selectedVolume.map((fileList))
            filesTable = <VolumeFilesTable fileList={this.props.selectedVolume}/>
        }

        return (
            <div style={{border: "1px solid black"}}>
                <h4> Selected job</h4>
                <JobRow tag="ID" value={job.ID}/>
                <JobRow tag="Name" value={title}/>
                <JobRow tag="Input" value={job.Input}/>
                <JobRow tag="Service ID" value={job.ServiceID}/>
                <JobRow tag="Service Description" value={service ? service.Description : false}/>
                <JobRow tag="Service Version" value={service ? service.Version : false}/>
                <JobRow style={{marginTop:'1em'}} tag="Status" value={job.State.Status}/>
                <JobRow style={{marginTop:'1em'}} tag="Error" value={job.State.Error ? job.State.Error : false}/>
                <Row style={{marginTop:'1em'}}>
                    <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>Input Volume</Col>
                    <Col xs={12} sm={9} md={9} >
                        <button type="submit" className="btn btn-default" onClick={this.handleInspectInputVolume.bind(this)}>Inspect</button>
                    </Col>
                </Row>
                <Row style={{marginTop:'1em'}}>
                    <Col xs={12} sm={3} md={3} style={{fontWeight:700}}>Output Volume</Col>
                    <Col xs={12} sm={9} md={9} >
                        <button type="submit" className="btn btn-default" onClick={this.handleInspectOutputVolume.bind(this)}>Inspect</button>
                    </Col>
                </Row>
                {filesTable}

            </div>

        )
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

Job.propTypes = {
    job: PropTypes.object.isRequired
};

//export default Job
export default connect(mapStateToProps, mapDispatchToProps)(Job);