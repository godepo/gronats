package gronats

import (
	"context"
	"sync"
	"testing"

	"github.com/godepo/groat/pkg/ctxgroup"
	natsClient "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestBootstrapper(t *testing.T) {
	t.Run("should be able failure", func(t *testing.T) {
		t.Run("when  can construct inner container injector", func(t *testing.T) {
			exp := UnexpectError()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			cont := NewMockNATsContainer(t)

			wg := &sync.WaitGroup{}
			ctx = ctxgroup.WithWaitGroup(ctx, wg)

			term := cont.EXPECT().Terminate(mock.Anything)
			term.Return(nil)

			cont.EXPECT().ConnectionString(mock.Anything).Return("", exp)

			_, err := bootstrapper[Deps](config{
				natsConnector: func(url string, options ...natsClient.Option) (*natsClient.Conn, error) {
					return nil, exp
				},
				runner: func(
					ctx context.Context,
					img string, opts ...testcontainers.ContainerCustomizer,
				) (NATsContainer, error) {
					return cont, nil
				},
			})(ctx)
			require.ErrorIs(t, err, exp)

			cancel()

			wg.Wait()
		})

		t.Run("when cant construct connect to nats", func(t *testing.T) {
			exp := UnexpectError()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			cont := NewMockNATsContainer(t)

			wg := &sync.WaitGroup{}
			ctx = ctxgroup.WithWaitGroup(ctx, wg)

			term := cont.EXPECT().Terminate(mock.Anything)
			term.Return(nil)

			cont.EXPECT().ConnectionString(mock.Anything).Return("", nil)

			_, err := bootstrapper[Deps](config{
				natsConnector: func(url string, options ...natsClient.Option) (*natsClient.Conn, error) {
					return nil, exp
				},
				runner: func(
					ctx context.Context,
					img string, opts ...testcontainers.ContainerCustomizer,
				) (NATsContainer, error) {
					return cont, nil
				},
			})(ctx)
			require.ErrorIs(t, err, exp)

			cancel()

			wg.Wait()
		})

		t.Run("when cant run container", func(t *testing.T) {
			exp := UnexpectError()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			wg := &sync.WaitGroup{}
			ctx = ctxgroup.WithWaitGroup(ctx, wg)

			_, err := bootstrapper[Deps](config{
				natsConnector: func(url string, options ...natsClient.Option) (*natsClient.Conn, error) {
					return nil, exp
				},
				runner: func(
					ctx context.Context,
					img string, opts ...testcontainers.ContainerCustomizer,
				) (NATsContainer, error) {
					return nil, exp
				},
			})(ctx)
			require.ErrorIs(t, err, exp)

			cancel()

			wg.Wait()
		})
	})
}
