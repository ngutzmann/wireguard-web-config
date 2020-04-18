package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/ngutzmann/wireguard-web-config/graph/generated"
	"github.com/ngutzmann/wireguard-web-config/graph/model"
	"github.com/ngutzmann/wireguard-web-config/graph/resolvers"
)

func (r *mutationResolver) CreatePeer(ctx context.Context, input model.NewPeer) (*model.Peer, error) {
	db := r.Resolver.DB
	return resolvers.CreatePeer(ctx, db, input)
}

func (r *queryResolver) Peers(ctx context.Context) ([]*model.Peer, error) {
	db := r.Resolver.DB
	return resolvers.GetPeers(ctx, db)
}

func (r *queryResolver) Peer(ctx context.Context, id string) (*model.Peer, error) {
	db := r.Resolver.DB
	return resolvers.GetPeer(ctx, db, id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
