//go:generate gorunpkg github.com/99designs/gqlgen

package resolver

import (
	context "context"
	graph "golang-gqlgen-reactjs-subscription-demo/golang/app/graph"
	model "golang-gqlgen-reactjs-subscription-demo/golang/app/model"
)

type Resolver struct{}

func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() graph.SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) AddChannel(ctx context.Context, name string) (model.Channel, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateChannel(ctx context.Context, id int, name string) (model.Channel, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteChannel(ctx context.Context, ID int) (model.Channel, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

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
	return nil, nil
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) SubscriptionChannelAdded(ctx context.Context) (<-chan model.Channel, error) {
	panic("not implemented")
}
func (r *subscriptionResolver) SubscriptionChannelDeleted(ctx context.Context) (<-chan model.Channel, error) {
	panic("not implemented")
}
func (r *subscriptionResolver) SubscriptionChannelUpdated(ctx context.Context) (<-chan model.Channel, error) {
	panic("not implemented")
}
