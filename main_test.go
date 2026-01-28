package gronats

import (
	"os"
	"testing"

	"github.com/godepo/groat"
	"github.com/godepo/groat/integration"
	natsClient "github.com/nats-io/nats.go"
)

type (
	SystemUnderTest struct {
	}

	State struct {
	}
	Deps struct {
		Client     *natsClient.Conn `groat:"gronats"`
		ConnString string           `groat:"gronats.config"`
		ConnPrefix string           `groat:"gronats.prefix"`
	}
)

var suite *integration.Container[Deps, State, *SystemUnderTest]

func TestMain(m *testing.M) {
	if os.Getenv("GROAT_I9N_NATS_IMAGE") == "" {
		_ = os.Setenv("GROAT_I9N_NATS_IMAGE", "nats:2.9")
	}

	suite = integration.New[Deps, State, *SystemUnderTest](
		m,
		func(t *testing.T) *groat.Case[Deps, State, *SystemUnderTest] {
			tcs := groat.New[Deps, State, *SystemUnderTest](t, func(t *testing.T, deps Deps) *SystemUnderTest {
				return &SystemUnderTest{}
			})
			return tcs
		},
		New[Deps](
			WithInjectLabel("gronats"),
			WithContainerImage("nats:2.6"),
			WithImageEnvValue("GROAT_I9N_NATS_IMAGE"),
			WithInjectLabelCasePrefix("gronats.prefix"),
			WithInjectLabelDSN("gronats.config"),
			WithNameSpaceLabel("gronats"),
		),
	)
	os.Exit(suite.Go())
}
