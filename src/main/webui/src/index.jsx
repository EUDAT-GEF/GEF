/**
 * Created by wqiu on 18/08/16.
 */
'use strict';
import React from 'react';
import createLogger from 'redux-logger';
import { Provider } from 'react-redux';
import { render } from 'react-dom';
import App from './containers/App';
import bows from 'bows';


const log = bows('app');


render(
    <App />,
    document.getElementById('react')
);


