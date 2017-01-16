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
import { syncHistoryWithStore, routerReducer, routerMiddleware } from 'react-router-redux'
import thunkMiddleware from 'redux-thunk';
import gefReducers from '../reducers/reducers';
import { compose, applyMiddleware, combineReducers, createStore} from 'redux';
import createLogger from 'redux-logger';
import { Provider } from 'react-redux';

require('react-s-alert/dist/s-alert-default.css');
require('react-s-alert/dist/s-alert-css-effects/slide.css');


const logger = createLogger();

const finalCreateStore = compose(
    applyMiddleware(
        thunkMiddleware,
        routerMiddleware(browserHistory),
        logger
    ),
    window.devToolsExtension? window.devToolsExtension() : f=>f
)(createStore);

const store = finalCreateStore(gefReducers);

const history = syncHistoryWithStore(browserHistory, store);

class App extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {

    }


    render() {
        return (
            <Provider store={store}>
                <div>
                    <div style={{'paddingBottom': 70}}>
                        <Header />
                        <Router history={history}>
                            <Route path='/' component={Main}>
                                <Route path='jobs' component={BrowseJobsContainer} >
                                    <Route path=':id' />
                                </Route>
                                <Route path='buildImage' component={BuildServiceContainer} />
                                <Route path='services' component={BrowseServicesContainer} >
                                    <Route path=':id' />
                                </Route>
                                <Route path='volumes' component={BrowseVolumesContainer} >
                                    <Route path=':id' />
                                </Route>
                                <Route path='buildVolume' component={BuildVolumeContainer} />
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

export default App;
