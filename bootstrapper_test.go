package gronats

import (
	"context"
	"sync"
	"testing"

	"github.com/godepo/groat/pkg/ctxgroup"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

func TestBootstrapper(t *testing.T) {
	t.Run("should be able failure", func(t *testing.T) {
		t.Run("when  can't run container", func(t *testing.T) {
			exp := UnexpectError()

			_, err := bootstrapper[Deps](config{
				runner: func(
					ctx context.Context,
					img string, opts ...testcontainers.ContainerCustomizer,
				) (NATsContainer, error) {
					return nil, exp
				},
			})(t.Context())
			require.ErrorIs(t, err, exp)
		})

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
	})
}
