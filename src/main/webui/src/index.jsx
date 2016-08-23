/**
 * Created by wqiu on 18/08/16.
 */
'use strict';
import React from 'react';
import createLogger from 'redux-logger';
import { Provider } from 'react-redux';
import { render } from 'react-dom';
import thunkMiddleware from 'redux-thunk';
import gefReducers from './reducers/reducers';
import { compose, applyMiddleware, combineReducers, createStore} from 'redux';
import App from './containers/App';
import bows from 'bows';


const log = bows('app');
const logger = createLogger();

const finalCreateStore = compose(
    applyMiddleware(
        thunkMiddleware,
        logger
    ),
    window.devToolsExtension? window.devToolsExtension() : f=>f
)(createStore);

const store = finalCreateStore(gefReducers);

render(
    <Provider store={store}>
        <App />
    </Provider>,
    document.getElementById('react')
);


