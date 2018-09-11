import React, { Component } from 'react';
import { graphql, Mutation, compose } from 'react-apollo';
import gql from 'graphql-tag';
import _ from 'underscore';
import update from 'react-addons-update';
require('./ChannelList.css');

class ChannelsList extends Component {
    constructor(props) {
        super(props);
        this.state = {
            id: "",
            nameValue: "",
        }
    }

    handleKeyUp = (evt) => {
        if (evt.keyCode === 13) {
            evt.persist();
            if (this.state.id !== "") {
                this.props.updateChannelMutation({
                    variables: {
                        name: evt.target.value,
                        id: this.state.id
                    },
                }).then((data) => {
                    evt.target.value = '';
                });
                this.setState({
                    id: ""
                })
            } else {
                this.props.createChannelMutation({
                    variables: {
                        name: evt.target.value
                    },
                }).then((data) => {
                    evt.target.value = '';
                });
            }
        };
    }


    componentWillMount() {
        this.props.data.subscribeToMore({
            document: addChannelSubscription,  // Use the subscription
            updateQuery: (prev, { subscriptionData }) => {
                if (!subscriptionData.data) {
                    return prev;
                }
                const newChannel = subscriptionData.data.subscriptionChannelAdded;
                // Add check to prevent double adding of channels.
                if (!prev.channels.find((channel) => channel.id === newChannel.id)) {
                    let updatedChannels = Object.assign({}, prev, { channels: [...prev.channels, newChannel] });
                    return updatedChannels;
                } else {
                    return prev;
                }
            }
        });
        this.props.data.subscribeToMore({
            document: deleteChannelSubscription,  // Use the subscription
            updateQuery: (prev, { subscriptionData }) => {
                if (!subscriptionData.data) {
                    return prev;
                }
                const newChannel = subscriptionData.data.subscriptionChannelDeleted;
                const deleteIndex = _.findIndex(prev.channels, (item) => item.id === newChannel.id);
                if (deleteIndex < 0) {
                    return prev;
                }
                return update(prev, {
                    channels: {
                        $splice: [[deleteIndex, 1]]
                    }
                });
            }
        });
        this.props.data.subscribeToMore({
            document: updateChannelSubscription,  // Use the subscription
            updateQuery: (prev, { subscriptionData }) => {
                if (!subscriptionData.data) {
                    return prev;
                }
            }
        });
    }

    onEditClick = (ch) => {
        this.setState({
            nameValue: ch.name,
            id: ch.id
        });
    }
    onChangeFunc = (event) => {
        this.setState({
            nameValue: event.target.value
        })
    }
    render() {
        const { data: { loading, error, channels } } = this.props;
        if (loading) {
            return <p>Loading ...</p>;
        }
        if (error) {
            return <p>{error.message}</p>;
        }
        return (
            <div className="container">
                <div className="row justify-content-md-center">
                    <div className="col-md-8">
                        <center><h3 className="main-title">Channel List</h3></center>
                        <hr />
                        <div className="input-group mb-3">
                            <input
                                type="text"
                                className="form-control"
                                placeholder="Enter Channel"
                                onKeyUp={this.handleKeyUp}
                                value={this.state.nameValue}
                                onChange={this.onChangeFunc} />
                        </div>
                        <ul className="list-group">
                            {channels.map(ch =>
                                <li key={"div_" + ch.id} className="list-group-item">
                                    <label htmlFor="checkbox5">
                                        {ch.id} : {ch.name}
                                    </label>
                                    <div className="pull-right action-buttons">
                                        <a onClick={() => this.onEditClick(ch)} href="javascript:void(0)"><span className="fa fa-pencil"></span></a>

                                        <Mutation mutation={deleteChannel} >
                                            {(deleteChannelMutation, { data }) => (
                                                <a className="trash" href="javascript:void(0)" onClick={() => { deleteChannelMutation({ variables: { id: ch.id } }); }}><span className="fa fa-trash"></span></a>

                                            )}
                                        </Mutation>
                                    </div>
                                </li>
                            )}
                        </ul>
                    </div>
                </div>
            </div>
        )

    }

}

export const channelsListQuery = gql`
   query ChannelsListQuery {
     channels {
       id
       name
     }
   }
 `;

const addChannelSubscription = gql`
    subscription Channels {
     subscriptionChannelAdded {
       id
       name
     }
    }
`

const deleteChannelSubscription = gql`
    subscription Channels {
     subscriptionChannelDeleted {
       id
       name
     }
    }
`
const updateChannelSubscription = gql`
    subscription Channels {
     subscriptionChannelUpdated {
       id
       name
     }
    }
`
const deleteChannel = gql`
  mutation deleteChannelMutation($id: Int!) {
    deleteChannel(ID: $id) {
      id
      name
    }
  }
`;

const updateChannel = gql`
  mutation updateChannelMutation($id: Int!,$name:String!) {
    updateChannel(id: $id,name:$name) {
      id
      name
    }
  }
`;

const CreateChannelMutation = gql`
  mutation addChannel($name: String!) {
    addChannel(name: $name) {
      id
      name
    }
  }
`;
const multipleMutation = compose(
    graphql(CreateChannelMutation, { name: 'createChannelMutation' }),
    graphql(updateChannel, { name: 'updateChannelMutation' })
)

export default compose(multipleMutation, graphql(channelsListQuery))(ChannelsList);
