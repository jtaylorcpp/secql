package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/jtaylorcpp/secql/graph/generated"
	"github.com/jtaylorcpp/secql/graph/model"
)

func (r *queryResolver) Ec2Instances(ctx context.Context) ([]*model.EC2Instance, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
