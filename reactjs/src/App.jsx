import React, { Component } from 'react';
import './style/app.scss';
import ChannelList from './components/ChannelList/ChannelList';
// import CreateChannel from './components/CreateChannel/CreateChannel';

class App extends Component {
    componentDidMount() {
        // services.setupInterceptors()
    }
    render() {
        return (
            <div>
                {/* <CreateChannel /> */}
                <ChannelList />
            </div>
        );
    }
}


export default App;
