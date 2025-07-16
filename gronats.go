// Package gronats provides NATS integration for testing using the groat framework.
// It allows easy setup and usage of NATS server in Docker containers for integration tests.
//
//go:generate go tool mockery
package gronats

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/godepo/groat/integration"
	"github.com/godepo/groat/pkg/ctxgroup"
	"github.com/godepo/gronats/internal/pkg/containersync"
	natsClient "github.com/nats-io/nats.go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/nats"
)

type (
	NATsContainer interface {
		ConnectionString(ctx context.Context) (string, error)
		Terminate(ctx context.Context, opts ...testcontainers.TerminateOption) error
	}

	containerRunner func(
		ctx context.Context,
		img string,
		opts ...testcontainers.ContainerCustomizer,
	) (NATsContainer, error)

	// Container represents a NATS integration container for testing.
	// Provides NATS server instance and manages its lifecycle during tests.
	Container[T any] struct {
		forks                       *atomic.Int32
		ctx                         context.Context
		connString                  string
		injectLabel                 string
		injectLabelConnectionString string
		injectLabelPrefix           string
		natsClient                  *natsClient.Conn
		connector                   func(url string, options ...natsClient.Option) (*natsClient.Conn, error)
	}

	config struct {
		imageEnvValue                string
		containerImage               string
		runner                       containerRunner
		injectLabel                  string
		injectLabelsConnectionString string
		injectLabelPrefix            string
		natsConnector                func(url string, options ...natsClient.Option) (*natsClient.Conn, error)
	}

	// Option represents a configuration option for configure the NATS container construction.
	// Used to customize container behavior and specify dependency injection.
	Option func(*config)
)

// WithInjectLabel sets the label for NATS client dependency injection.
// label is the label tag name in Deps structure, used in groat for dependency identification.
// By default, the value is "nats"
//
//		example:
//	 type Deps structure {
//			NATSClient *nats.Conn `groat:"nats"`
//	 }
func WithInjectLabel(label string) Option {
	return func(cfg *config) {
		cfg.injectLabel = label
	}
}

// WithInjectLabelDSN sets the label for connection string injection.
// Label is the label name used for connection string dependency injection.
// By default, the value is "nats.config"
//
//		example:
//	 type Deps structure {
//			DSN string `groat:"nats.config"`
//	 }
func WithInjectLabelDSN(label string) Option {
	return func(cfg *config) {
		cfg.injectLabelsConnectionString = label
	}
}

// WithInjectLabelCasePrefix sets the label for prefix injection.
// Label is the label name used for prefix dependency injection.
// Prefix is a string like "$num_", where $num is a numerical value for a test case.
// By default, the value is "nats.prefix"
//
//		example:
//	 type Deps structure {
//			CasePrefix string `groat:"nats.prefix"`
//	 }
func WithInjectLabelCasePrefix(label string) Option {
	return func(cfg *config) {
		cfg.injectLabelPrefix = label
	}
}

// WithContainerImage sets the Docker image to use for the NATS container.
// image is the Docker image name and tag (e.g., "nats:2.6").
func WithContainerImage(img string) Option {
	return func(cfg *config) {
		cfg.containerImage = img
	}
}

// WithImageEnvValue sets the environment variable name for the Docker image.
// envVar is the environment variable name that contains the image specification.
// By default, GROAT_I9N_NATS_IMAGE.
func WithImageEnvValue(env string) Option {
	return func(cfg *config) {
		cfg.imageEnvValue = env
	}
}

// New creates a new factory method for build NATS integration container with the given options.
func New[T any](options ...Option) integration.Bootstrap[T] {
	cfg := config{
		imageEnvValue:                "GROAT_I9N_NATS_IMAGE",
		containerImage:               "nats:2.9",
		injectLabel:                  "nats",
		injectLabelsConnectionString: "nats.config",
		injectLabelPrefix:            "nats.prefix",
		natsConnector:                natsClient.Connect,
		runner: func(ctx context.Context, img string, opts ...testcontainers.ContainerCustomizer) (NATsContainer, error) {
			return nats.Run(ctx, img, opts...)
		},
	}

	for _, op := range options {
		op(&cfg)
	}

	if env := os.Getenv(cfg.imageEnvValue); env != "" {
		cfg.containerImage = env
	}

	return bootstrapper[T](cfg)
}

func bootstrapper[T any](cfg config) integration.Bootstrap[T] {
	return func(ctx context.Context) (integration.Injector[T], error) {
		natsContainer, err := cfg.runner(ctx, cfg.containerImage)
		if err != nil {
			return nil, fmt.Errorf("nats container failed to run: %w", err)
		}

		ctxgroup.IncAt(ctx)

		go containersync.Terminator(ctx, natsContainer.Terminate)()

		container, err := newContainer[T](ctx, natsContainer, cfg)
		if err != nil {
			return nil, err
		}

		return container.Injector, nil
	}
}
