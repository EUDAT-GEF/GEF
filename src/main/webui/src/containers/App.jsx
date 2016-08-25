'use strict';

import React from 'react';
import Header from '../components/Header'
import Main from './Main';
import Footer from '../components/Footer';

class App extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {

    }


    render() {
        return (
            <div>
                <Header />
                <Main />
                <Footer version="0.4.0" />
            </div>
        );
    }
}

export default App;