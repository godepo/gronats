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
	dsn string,
	cfg config,
) (*Container[T], error) {
	natsClient, err := cfg.natsConnector(dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to nats: %w", err)
	}

	container := &Container[T]{
		forks:                       &atomic.Int32{},
		ctx:                         ctx,
		connString:                  dsn,
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

	if c.injectLabelPrefix != "" {
		prefix = c.injectLabelPrefix + "_" + prefix
	}

	res := generics.Injector(t, c.natsClient, to, c.injectLabel)
	res = generics.Injector(t, c.connString, res, c.injectLabelConnectionString)
	res = generics.Injector(t, prefix, res, c.injectLabelPrefix)

	return res
}
