import React, {PropTypes} from 'react';
import { Row, Col, Grid, Table } from 'react-bootstrap';
import { toPairs } from '../utils/utils';
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import actions from '../actions/actions';


const VolumeFile = ({handleFileClick, file, iconClass, indentStyle}) => (

    <li className="row file" style={{lineHeight:2}} onClick={handleFileClick()}>
        <div className="col-sm-6">
            <span style={indentStyle}/>
            <button style={{width:20, background:'none', border:'none', fontSize:20, padding:0}}>+</button>

            <span className={"glyphicon "+iconClass} aria-hidden={true} /> {file.name}
        </div>
        <div className="col-sm-3">{file.size}</div>
        <div className="col-sm-3">{file.date}</div>
    </li>
);


// const iconClass = !file.isdir ? "glyphicon-file"
//                               : file.children == undefined ? "glyphicon-folder-close"
//                               : "glyphicon-folder-open";
//           const size = file.size ? humanSize(file.size) : "";
//           const date = moment(file.date).format('ll');
//           const indentStyle = {paddingLeft: (3*file.indent)+'em'};
//           const handlerStyle = {width:20, background:'none', border:'none', fontSize:20, padding:0};
//
const VolumeFilesTable1 = ({fileList}) => (
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
                       fileListItem.isFolder == true ?
                              (<button style={{width:20, background:'none', border:'none', fontSize:20, padding:0}}>+</button>)

                       <span className={"glyphicon "+iconClass} aria-hidden={true} /> {fileListItem.name}
                   </div>
                   <div className="col-sm-3">{fileListItem.size}</div>
                   <div className="col-sm-3">{fileListItem.modified}</div>
               </li>
            })}
        </ol>
    </div>
);

class FileTree extends React.Component {
    constructor(props) {
        super(props);
    }

    handleFileClick() {
        console.log("File was clicked");
        //this.props.actions.inspectVolume(this.props.job.InputVolume)
    }
    render() {
        console.log(this.props);
        if (this.props.selectedVolume.length > 0) {
            return (
                <div style={{margin:'1em'}}>
                    <ol className="list-unstyled fileList" style={{textAlign:'left', minHeight:'30em'}}>
                        <li className="heading row" style={{padding:'0.5em 0'}}>
                            <div className="col-sm-6" style={{fontWeight:'bold'}}>File Name</div>
                            <div className="col-sm-3" style={{fontWeight:'bold'}}>Size</div>
                            <div className="col-sm-3" style={{fontWeight:'bold'}}>Date</div>
                        </li>

                        {this.props.selectedVolume.map((fileListItem, index) => {

                            let indentStyle1 = {paddingLeft: (3*1)+'em'};
                            let iconClass1 = "glyphicon-file";
                            if (fileListItem.isFolder == true) {
                                iconClass1 = "glyphicon-folder-close";
                            }

                            //this.volumeFile(fileListItem, index, iconClass, indentStyle)
                            //console.log(index)
                            //console.log(fileListItem)
                            //console.log(fileListItem)
                            return <VolumeFile handleFileClick={this.handleFileClick} key={"vf"+index} file={fileListItem} iconClass={iconClass1} indentStyle={indentStyle1}/>





                        })}
                    </ol>
                </div>
            )
        }
        else {
            return (
                <div>Press Inspect to see files</div>
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

FileTree.propTypes = {
    handleFileClick: PropTypes.func.isRequired,
};

export default connect(mapStateToProps, mapDispatchToProps)(FileTree);