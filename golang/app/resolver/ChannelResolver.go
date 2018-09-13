package resolver

import (
	context "context"
	"golang-gqlgen-reactjs-subscription-demo/golang/app/config/connection"
	"golang-gqlgen-reactjs-subscription-demo/golang/app/model"
)

type channelResolver struct{ *Resolver }

var addChannelObserver map[string]chan model.Channel
var deleteChannelObserver map[string]chan model.Channel
var updateChannelObserver map[string]chan model.Channel

func init() {
	addChannelObserver = map[string]chan model.Channel{}
	deleteChannelObserver = map[string]chan model.Channel{}
	updateChannelObserver = map[string]chan model.Channel{}
}
func (r *queryResolver) Channels(ctx context.Context) ([]model.Channel, error) {
	db := connection.DbConn()
	var query = "SELECT * FROM channel"
	selDB, err := db.Query(query)
	var arrChannel []model.Channel
	for selDB.Next() {
		var name string
		var id int64
		err = selDB.Scan(&id, &name)
		if err != nil {
			panic(err.Error())
		}
		todo1 := model.Channel{ID: int(id), Name: name}
		arrChannel = append(arrChannel, todo1)
	}

	defer db.Close()
	return arrChannel, nil
}

func (r *mutationResolver) AddChannel(ctx context.Context, name string) (model.Channel, error) {
	db := connection.DbConn()

	insForm, err := db.Prepare("INSERT INTO channel(name) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}

	var newChannel model.Channel
	res, err := insForm.Exec(name)
	if err != nil {
		println("Exec err:", err.Error())
	} else {
		var id int64
		id, err := res.LastInsertId()
		if err != nil {
			println("Error:", err.Error())
		} else {
			newChannel = model.Channel{ID: int(id), Name: name}
		}
	}
	defer db.Close()
	//add new chanel in observer
	for _, observer := range addChannelObserver {
		observer <- newChannel
	}
	return newChannel, nil
}
func (r *mutationResolver) DeleteChannel(ctx context.Context, ID int) (model.Channel, error) {
	db := connection.DbConn()

	delForm, err := db.Prepare("DELETE FROM channel WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(ID)
	var newChannel model.Channel
	defer db.Close()
	newChannel = model.Channel{ID: ID, Name: ""}
	for _, observer := range deleteChannelObserver {
		observer <- newChannel
	}
	return newChannel, nil
}
func (r *mutationResolver) UpdateChannel(ctx context.Context, id int, name string) (model.Channel, error) {
	db := connection.DbConn()

	insForm, err := db.Prepare("UPDATE channel SET name=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	var newChannel model.Channel
	insForm.Exec(name, id)
	newChannel = model.Channel{ID: id, Name: name}

	defer db.Close()
	//add new chanel in observer
	for _, observer := range updateChannelObserver {
		observer <- newChannel
	}
	return newChannel, nil
}

func (r *subscriptionResolver) SubscriptionChannelAdded(ctx context.Context) (<-chan model.Channel, error) {
	id := randString(8)
	events := make(chan model.Channel, 1)

	go func() {
		<-ctx.Done()
		delete(addChannelObserver, id)
	}()

	addChannelObserver[id] = events

	return events, nil
}

func (r *subscriptionResolver) SubscriptionChannelDeleted(ctx context.Context) (<-chan model.Channel, error) {
	id := randString(8)
	events := make(chan model.Channel, 1)

	go func() {
		<-ctx.Done()
		delete(deleteChannelObserver, id)
	}()

	deleteChannelObserver[id] = events

	return events, nil
}

func (r *subscriptionResolver) SubscriptionChannelUpdated(ctx context.Context) (<-chan model.Channel, error) {
	id := randString(8)
	events := make(chan model.Channel, 1)

	go func() {
		<-ctx.Done()
		delete(updateChannelObserver, id)
	}()

	updateChannelObserver[id] = events

	return events, nil
}
