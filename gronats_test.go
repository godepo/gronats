package gronats

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGronats(t *testing.T) {
	t.Run("should be able to run", func(t *testing.T) {
		tc := suite.Case(t)
		require.NotNil(t, tc.Deps.Client)
		require.NotEmpty(t, tc.Deps.ConnString)
		require.NotEmpty(t, tc.Deps.ConnPrefix)
	})
}
