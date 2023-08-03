package spanner

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
)

func NewClient(
	projectID string,
	instanceID string,
	databaseID string,
) *spanner.Client {
	ctx := context.Background()
	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, databaseID)
	c, err := spanner.NewClient(ctx, dsn)
	if err != nil {
		panic(err)
	}
	return c
}
