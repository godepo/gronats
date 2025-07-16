# GroNATs

[![codecov](https://codecov.io/gh/godepo/gronats/graph/badge.svg?token=PxxK11IBFs)](https://codecov.io/gh/godepo/gronats)
[![Go Report Card](https://goreportcard.com/badge/godepo/gronats)](https://goreportcard.com/report/godepo/gronats)
[![License](https://img.shields.io/badge/License-MIT%202.0-blue.svg)](https://github.com/godepo/gronats/blob/main/LICENSE)

Gronats is a Go library that provides seamless NATS integration for testing using the [groat](https://github.com/godepo/groat) test suite. It simplifies the setup and management of NATS servers in Docker containers for integration tests.

## Features

- 🐳 **Docker Integration**: Automatically manages NATS server containers for testing
- 🔧 **Easy Configuration**: Simple API for configuring NATS connections
- 🧪 **Test-Friendly**: Designed specifically for integration testing scenarios
- 📦 **Dependency Injection**: Seamless integration with groat's dependency injection system
- 🚀 **Quick Setup**: Minimal configuration required to get started
- 🔄 **Lifecycle Management**: Automatic container startup and cleanup

## Installation

## Quick Start

Here's a basic example of how to use Gronats in your tests:



Create file with name main_test.go or cases_test.go:

```go
package themodule
import (
	"os"
	"testing"
	"github.com/godepo/groat"
	"github.com/godepo/groat/integration"
	"github.com/godepo/gronats"
	natsClient "github.com/nats-io/nats.go"

)

type Deps struct { 
	Client *natsClient.Conn  `groat:"nats"`
	ConnString string `groat:"nats.config"`
	ConnPrefix string `groat:"nats.prefix"`
}

type State struct { 
	// Your test state here 
}

type SystemUnderTest struct { 
	// Your system under test 
}

var suite *integration.Container[Deps, State, *SystemUnderTest]

func TestMain(m *testing.M) {
	suite = integration.New[Deps, State, *SystemUnderTest](
		m,
		func(t *testing.T) *groat.Case[Deps, State, *SystemUnderTest] {
			return groat.New[Deps, State, *SystemUnderTest](t, func(t *testing.T, deps Deps) *SystemUnderTest {
				return &SystemUnderTest{}
			})
		},
		gronats.New[Deps](),
	)
	os.Exit(suite.Go())

}

```

and write u firs test, similar like this:

```go
func TestNatsIntegration(t *testing.T) { 
	t.Run("should connect to NATS", func(t *testing.T) { tc := suite.Case(t)
        // Your NATS client is ready to use
        err := tc.Deps.Client.Publish("test.subject", []byte("hello"))
        require.NoError(t, err)
        
        // Connection string and prefix are also available
        t.Logf("Connected to: %s", tc.Deps.ConnString)
        t.Logf("Using prefix: %s", tc.Deps.ConnPrefix)
    })
}

```

## Configuration Options

Gronats provides several configuration options to customize the NATS container:

### `gronats.WithInjectLabel(label string)`
Sets the label for NATS client dependency injection.


### `WithContainerImage(image string)`
Specifies the Docker image to use for the NATS container.

### `WithImageEnvValue(envVar string)`
Sets the environment variable name that contains the Docker image specification.

### `WithInjectLabelPrefix(label string)`
Sets the label for test case prefix injection.

### `WithInjectLabelDSN(label string)`
Sets the label for connection string injection.


## Use Cases

### 1. Message Queue Testing

Test your message queue implementations with real NATS server:

```go

func TestMessageQueue(t *testing.T) { t.Run("should process messages", func(t *testing.T) { 
	tc := suite.Case(t)
    // Subscribe to a subject
    ch := make(chan *nats.Msg, 64)
    sub, err := tc.Deps.Client.ChanSubscribe("orders.created", ch)
    require.NoError(t, err)
    defer sub.Unsubscribe()
    
    // Publish a message
    err = tc.Deps.Client.Publish("orders.created", []byte(`{"orderId": "123"}`))
    require.NoError(t, err)
    
    // Verify message received
    select {
    case msg := <-ch:
        assert.Equal(t, `{"orderId": "123"}`, string(msg.Data))
    case <-time.After(5 * time.Second):
        t.Fatal("Message not received")
    }
})
}

```


### 2. Microservice Communication Testing

Test communication between microservices:

```go
func TestMicroserviceCommunication(t *testing.T) { 
	t.Run("should handle request-response", func(t *testing.T) { 
		tc := suite.Case(t)
    // Set up a service responder
    _, err := tc.Deps.Client.Subscribe("user.get", func(msg *nats.Msg) {
        response := `{"userId": "123", "name": "John Doe"}`
        msg.Respond([]byte(response))
    })
    require.NoError(t, err)
    
    // Test the request-response pattern
    response, err := tc.Deps.Client.Request("user.get", []byte(`{"userId": "123"}`), 5*time.Second)
    require.NoError(t, err)
    
    var user struct {
        UserID string `json:"userId"`
        Name   string `json:"name"`
    }
    err = json.Unmarshal(response.Data, &user)
    require.NoError(t, err)
    assert.Equal(t, "123", user.UserID)
    assert.Equal(t, "John Doe", user.Name)
})

```


## Environment Variables

- `GROAT_I9N_NATS_IMAGE`: Override the default NATS Docker image
- Custom environment variables can be configured using `WithImageEnvValue()`

## Requirements

- Go 1.19 or higher
- Docker (for running NATS containers)
- [groat](https://github.com/godepo/groat) framework

## Dependencies

- `github.com/godepo/groat`: Testing framework
- `github.com/nats-io/nats.go`: NATS client library
- `github.com/stretchr/testify`: Testing assertions (for examples)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 [Issue Tracker](https://github.com/godepo/gronats/issues)

## Related Projects

- [groat](https://github.com/godepo/groat) - Go testing framework
- [NATS](https://nats.io/) - Cloud native messaging system
- [nats.go](https://github.com/nats-io/nats.go) - NATS client for Go
