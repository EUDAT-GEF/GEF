import React, {PropTypes} from 'react';
import { Row, Col, Grid, Table } from 'react-bootstrap';
import { toPairs } from '../utils/utils';
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import actions from '../actions/actions';
import moment from 'moment';

var path = require('path');

class FileTree extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            folderOpen: {},
        };
    }

    handleFolderClick(child, isOpen, e) {
        let folderOpen = this.state.folderOpen;
        let volumeInternalPath = path.join(child.path, child.name);

        folderOpen[volumeInternalPath] = isOpen;
        this.setState({folderOpen});



        console.log(this.state);

    }

    renderFile(file, indentStyle) {
        let volumeInternalPath = path.join(file.path, file.name);
        let iconClass = "glyphicon-file";
        let browseButton;
        let isContentVisible = true;

        if (file.isFolder) {
            iconClass = "glyphicon-folder-close";
            browseButton = <button style={{width:20, background:'none', border:'none', fontSize:20, padding:0}} onClick={ (e) => this.handleFolderClick(file, true, e)}>+</button>
            if (this.state.folderOpen[volumeInternalPath] != null) {
                if (this.state.folderOpen[volumeInternalPath]) {
                    //isContentVisible = true;
                    iconClass = "glyphicon-folder-open";
                    browseButton = <button style={{width:20, background:'none', border:'none', fontSize:20, padding:0}} onClick={ (e) => this.handleFolderClick(file, false, e)}>-</button>
                } //else {
                 //   isContentVisible = false;
                //}
            }
        }


        if (this.state.folderOpen[file.path] != null) {
            if (this.state.folderOpen[file.path]) {
                isContentVisible = true
            } else {
                isContentVisible = false
            }
        } else {
            isContentVisible = false
        }


        for (var folderName in this.state.folderOpen) {
            console.log(folderName);
            console.log("********")

            if ((file.path.indexOf(folderName) == 0) && (this.state.folderOpen[folderName]==false)) {
                isContentVisible = false;
            }

        }




        if (file.path == "") {
            isContentVisible = true;
        }

        if (isContentVisible) {
            return (
                <li key={volumeInternalPath} className="row file" style={{lineHeight: 2}}>
                    <div className="col-sm-6">
                        <span style={indentStyle}/>
                        {browseButton}
                        <span className={"glyphicon " + iconClass} aria-hidden={true}/> {file.name}
                    </div>
                    <div className="col-sm-3">{file.size}</div>
                    <div className="col-sm-3">{moment(file.date).format('ll')}</div>
                </li>
            )
        }
    }


    readVolumeContent(volumeContent, depthLevel, volumeItems) {
        volumeContent.map((fileListItem) => {
            let indentStyle = {paddingLeft: (3*(1+depthLevel))+'em'};


            volumeItems.push(this.renderFile(fileListItem, indentStyle))

            if (fileListItem.isFolder == true) {
                depthLevel += 1;
                volumeItems = this.readVolumeContent(fileListItem.folderTree, depthLevel, volumeItems);
                depthLevel -= 1;
            }
        })
        return volumeItems
    }

    isInsideFolder(volumeContent, folderPath, fileName, currentFolder) {
        volumeContent.map((fileListItem) => {



            //volumeItems.push(this.renderFile(fileListItem, indentStyle))

            if (fileListItem.isFolder == true) {

                isFound = this.isInsideFolder(volumeContent, folderPath, fileName);

            }

        })
        return isFound
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

export default connect(mapStateToProps, mapDispatchToProps)(FileTree);