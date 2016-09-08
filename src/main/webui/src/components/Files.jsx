'use strict';

import React, {PropTypes} from 'react';
import ReactDOMServer from 'react-dom/server';
import {Row, InputGroup} from 'react-bootstrap';
import DropzoneComponent from 'react-dropzone-component';

require('react-dropzone-component/styles/filepicker.css');
require('dropzone/dist/min/dropzone.min.css');



class Files extends React.Component {
    constructor(props) {
        super(props);

        this.djsConfig = {
            addRemoveLinks: true,
            autoProcessQueue: false,
            previewTemplate: ReactDOMServer.renderToStaticMarkup(
                <div className="dz-preview dz-file-preview">
                    <div className="dz-filename"><span data-dz-name="true"></span></div>
                    <img data-dz-thumbnail="true" />
                    {/*<div className="dz-details">*/}
                    {/*</div>*/}
                    <div className="dz-progress"><span className="dz-upload" data-dz-uploadprogress="true"></span></div>
                    <div className="dz-success-mark"><span>✔</span></div>
                    <div className="dz-error-mark"><span>✘</span></div>
                    <div className="dz-error-message"><span data-dz-errormessage="true"></span></div>
                </div>
            )
        };

        this.componentConfig = {
            postUrl: 'no-url'
        };
    }

    handleFileAdded(file) {
        console.log(file);
    }

    render() {
        const config = this.componentConfig;
        const djsConfig = this.djsConfig;

        const eventHandlers = {
            addedfile: this.handleFileAdded.bind(this)
        };

        return <DropzoneComponent config={config} eventHandlers={eventHandlers} djsConfig={djsConfig} />
    }

}

export default Files;






