package gronats

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	natsClient "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func UnexpectError() error {
	return errors.New(uuid.NewString())
}

func TestNewContainer(t *testing.T) {
	t.Run("should be able to fail", func(t *testing.T) {
		t.Run("when can't get connection string", func(t *testing.T) {
			exp := UnexpectError()

			cont := NewMockNATsContainer(t)

			cont.EXPECT().ConnectionString(mock.Anything).Return("", exp)

			_, err := newContainer[Deps](t.Context(), cont, config{})
			require.ErrorIs(t, err, exp)
		})

		t.Run("when construction nats client failed", func(t *testing.T) {
			exp := UnexpectError()

			cont := NewMockNATsContainer(t)

			cont.EXPECT().ConnectionString(mock.Anything).Return("", nil)

			_, err := newContainer[Deps](t.Context(), cont, config{
				natsConnector: func(url string, options ...natsClient.Option) (*natsClient.Conn, error) {
					return nil, exp
				},
			})
			require.ErrorIs(t, err, exp)
		})
	})

}
