'use strict';

import React from 'react';
import Alert from 'react-s-alert';
import Header from '../components/Header'
import Main from './Main';
import Footer from '../components/Footer';
import NotFound from '../components/NotFound';
import {Router, Route, browserHistory, hashHistory } from 'react-router';
import BrowseWorkflowsContainer from './WorkflowsContainer';
import BrowseJobsContainer from './JobsContainer';
import BuildServiceContainer from '../containers/BuildServiceContainer';
import BrowseServicesContainer from './ServicesContainer';
import BrowseVolumesContainer from './VolumesContainer';
import BuildVolumeContainer from '../containers/BuildVolumeContainer';

require('react-s-alert/dist/s-alert-default.css');
require('react-s-alert/dist/s-alert-css-effects/slide.css');

class App extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {

    }


    render() {
        return (
            <div>
                <div style={{'paddingBottom': 70}}>
                    <Header />
                    <Router history={browserHistory}>
                        <Route path='/' component={Main}>
                            <Route path='workflows' component={BrowseWorkflowsContainer} />
                            <Route path='jobs(/:id)' component={BrowseJobsContainer} />
                            <Route path='buildImage' component={BuildServiceContainer} />
                            <Route path='services(/:id)' component={BrowseServicesContainer} />
                            <Route path='volumes(/:id)' component={BrowseVolumesContainer} />
                            <Route path='buildVolume' component={BuildVolumeContainer} />
                            <Route path='*' component={NotFound} />
                        </Route>
                    </Router>
                </div>
                <Alert stack={{limit:5}}/>
                <Footer version="0.4.0" />
            </div>
        );
    }
}

export default App;