import React, {PropTypes} from 'react';
import { Row, Col, Grid, Table } from 'react-bootstrap';
import { toPairs } from '../utils/utils';
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import actions from '../actions/actions';


const VolumeFile = ({handleFileClick, file, iconClass, indentStyle}) => (

    <li className="row file" style={{lineHeight:2}} onClick={handleFileClick.bind(file)}>
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
        this.state = {
            itemChecked: {},
        };
    }

    handleFileClick() {
        console.log("File was clicked");
        console.log(this);
        //this.props.actions.volumeItemClick(this)
        //this.props.actions.inspectVolume(this.props.job.InputVolume)
    }


    renderFile(file, index, iconClass, indentStyle) {
            //return <VolumeFile handleFileClick={this.handleFileClick} key={fileListItem.path+"/"+fileListItem.name} file={fileListItem} iconClass={iconClass} indentStyle={indentStyle}/>
            return (
                 <li key={file.path+"/"+file.name}  className="row file" style={{lineHeight:2}} onClick={this.handleFileClick.bind(file)}>
                    <div className="col-sm-6">
                        <span style={indentStyle}/>
                        <button style={{width:20, background:'none', border:'none', fontSize:20, padding:0}}>+</button>

                        <span className={"glyphicon "+iconClass} aria-hidden={true} /> {file.name}
                    </div>
                    <div className="col-sm-3">{file.size}</div>
                    <div className="col-sm-3">{file.date}</div>
                </li>
            )
    }


    readVolumeContent(volumeContent, depthLevel, volumeItems) {
        volumeContent.map((fileListItem, index) => {
            let indentStyle = {paddingLeft: (3*(1+depthLevel))+'em'};
            let iconClass = "glyphicon-file";


            if (fileListItem.isFolder == true) {
                iconClass = "glyphicon-folder-close";
            }

            volumeItems.push(this.renderFile(fileListItem, index+volumeItems.length, iconClass, indentStyle))

            if (fileListItem.isFolder == true) {
                depthLevel += 1;
                volumeItems = this.readVolumeContent(fileListItem.folderTree, depthLevel, volumeItems);
                depthLevel -= 1;
            }
        })
        return volumeItems
    }
    render() {
        console.log(this.props);
        let sLines = [];
        if (this.props.selectedVolume.length > 0) {
            return (
                <div style={{margin:'1em'}}>
                    <ol className="list-unstyled fileList" style={{textAlign:'left', minHeight:'30em'}}>
                        <li className="heading row" style={{padding:'0.5em 0'}}>
                            <div className="col-sm-6" style={{fontWeight:'bold'}}>File Name</div>
                            <div className="col-sm-3" style={{fontWeight:'bold'}}>Size</div>
                            <div className="col-sm-3" style={{fontWeight:'bold'}}>Date</div>
                        </li>
                        {sLines = this.readVolumeContent(this.props.selectedVolume, 0, [])}


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