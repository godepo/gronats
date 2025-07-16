package gronats

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/godepo/groat/pkg/generics"
)

func newContainer[T any](
	ctx context.Context,
	natsContainer NATsContainer,
	cfg config,
) (*Container[T], error) {
	connectionString, err := natsContainer.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get connection string: %w", err)
	}

	natsClient, err := cfg.natsConnector(connectionString)
	if err != nil {
		return nil, fmt.Errorf("can't connect to nats: %w", err)
	}

	container := &Container[T]{
		forks:                       &atomic.Int32{},
		ctx:                         ctx,
		connString:                  connectionString,
		injectLabel:                 cfg.injectLabel,
		injectLabelConnectionString: cfg.injectLabelsConnectionString,
		injectLabelPrefix:           cfg.injectLabelPrefix,
		natsClient:                  natsClient,
		connector:                   cfg.natsConnector,
	}

	return container, nil
}

func (c *Container[T]) Injector(t *testing.T, to T) T {
	t.Helper()
	prefix := strconv.Itoa(int(c.forks.Add(1))) + "_"
	res := generics.Injector(t, c.natsClient, to, c.injectLabel)
	res = generics.Injector(t, c.connString, res, c.injectLabelConnectionString)
	res = generics.Injector(t, prefix, res, c.injectLabelPrefix)

	return res
}
