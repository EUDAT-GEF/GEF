import React from 'react';
import {render} from 'react-dom';
import {Router, Route, browserHistory, hashHistory} from 'react-router';
import {syncHistoryWithStore, routerMiddleware} from 'react-router-redux'
import thunkMiddleware from 'redux-thunk';
import {compose, applyMiddleware, createStore} from 'redux';
import createLogger from 'redux-logger';
import {Provider} from 'react-redux';
import bows from 'bows';
import Alert from 'react-s-alert';
require('react-s-alert/dist/s-alert-default.css');
require('react-s-alert/dist/s-alert-css-effects/slide.css');

import gefReducers from './reducers/reducers';

import Header from './components/Header'
import Footer from './components/Footer';
import NotFound from './components/NotFound';

import BrowseJobsContainer from './containers/JobsContainer';
import BuildServiceContainer from './containers/BuildServiceContainer';
import BrowseServicesContainer from './containers/ServicesContainer';
import BrowseVolumesContainer from './containers/VolumesContainer';
import BuildVolumeContainer from './containers/BuildVolumeContainer';
import Main from './containers/Main';


const log = bows('app');

const middleware = compose(
    applyMiddleware(thunkMiddleware,
                    routerMiddleware(browserHistory),
                    createLogger()),
    window.devToolsExtension || (f=>f)
);
const store = createStore(gefReducers, middleware);
const history = syncHistoryWithStore(browserHistory, store);

class App extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <Provider store={store}>
                <div>
                    <div style={{'paddingBottom': 70}}>
                        <Header />
                        <Router history={history}>
                            <Route path='/' component={Main}>
                                <Route path='builds' component={BuildServiceContainer} />
                                <Route path='services' component={BrowseServicesContainer} >
                                    <Route path=':id' />
                                </Route>
                                <Route path='jobs' component={BrowseJobsContainer} >
                                    <Route path=':id' />
                                </Route>
                                <Route path='*' component={NotFound} />
                            </Route>
                        </Router>
                    </div>
                    <Alert stack={{limit:5}}/>
                    <Footer version="0.4.0" />
                </div>
            </Provider>
        );
    }
}

render(<App />, document.getElementById('react') );
