package resolver

import (
	context "context"
	"fmt"
	"golang-gqlgen-reactjs-subscription-demo/golang/app/model"
)

// type channelResolver struct{ *resolver }
var channelList = []model.Channel{
	{ID: 1, Name: "Channel A"},
	{ID: 2, Name: "Channel B"},
}
var lastId = 2

var addChannelObserver map[string]chan model.Channel
var deleteChannelObserver map[string]chan model.Channel
var updateChannelObserver map[string]chan model.Channel

func init() {
	addChannelObserver = map[string]chan model.Channel{}
	deleteChannelObserver = map[string]chan model.Channel{}
	updateChannelObserver = map[string]chan model.Channel{}
}

/*
	- Function for run query
	-example :
	 	query{
			channels{
				name
				id
			}
		}
*/
func (r *queryResolver) Channels(ctx context.Context) ([]model.Channel, error) {
	return channelList, nil
}

func (r *mutationResolver) AddChannel(ctx context.Context, name string) (model.Channel, error) {
	fmt.Println("---------AddChannel-----------")

	lastId++
	newID := lastId
	newChannel := model.Channel{
		ID:   newID,
		Name: name,
	}
	channelList = append(channelList, newChannel)
	for _, observer := range addChannelObserver {
		observer <- newChannel
	}
	return newChannel, nil
}
func remove(s []model.Channel, i int) []model.Channel {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
func (r *mutationResolver) DeleteChannel(ctx context.Context, ID int) (model.Channel, error) {

	fmt.Println("---------DeleteChannel-----------")
	var newChannel model.Channel
	for i, v := range channelList {
		if v.ID == ID {
			channelList = remove(channelList, i)
			newChannel = model.Channel{ID: v.ID, Name: ""}
		}
	}
	for _, observer := range deleteChannelObserver {
		observer <- newChannel
	}
	return newChannel, nil
}
func (r *mutationResolver) UpdateChannel(ctx context.Context, id int, name string) (model.Channel, error) {

	fmt.Println("---------UpdateChannel-----------")
	var newChannel model.Channel
	for i, v := range channelList {
		if v.ID == id {
			channelList[i].Name = name
			newChannel = model.Channel{ID: v.ID, Name: name}
		}
	}

	for _, observer := range updateChannelObserver {
		observer <- newChannel
	}
	return newChannel, nil
}

func (r *subscriptionResolver) SubscriptionChannelAdded(ctx context.Context) (<-chan model.Channel, error) {
	id := randString(8)
	events := make(chan model.Channel, 1)

	addChannelObserver[id] = events

	return events, nil
}

func (r *subscriptionResolver) SubscriptionChannelDeleted(ctx context.Context) (<-chan model.Channel, error) {
	id := randString(8)
	events := make(chan model.Channel, 1)

	deleteChannelObserver[id] = events
	return events, nil
}

func (r *subscriptionResolver) SubscriptionChannelUpdated(ctx context.Context) (<-chan model.Channel, error) {
	id := randString(8)
	events := make(chan model.Channel, 1)

	updateChannelObserver[id] = events
	return events, nil
}
