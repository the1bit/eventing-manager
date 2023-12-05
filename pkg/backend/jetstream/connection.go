package jetstream

import (
	"github.com/nats-io/nats.go"

	"github.com/kyma-project/eventing-manager/pkg/env"
	"github.com/kyma-project/eventing-manager/pkg/errors"
)

type Builder interface {
	Build() (ConnectionInterface, error)
}

// ConnectionBuilder helps in establishing a connection to NATS.
type ConnectionBuilder struct {
	config env.NATSConfig
}

func NewConnectionBuilder(config env.NATSConfig) Builder {
	return ConnectionBuilder{config: config}
}

// Build connects to NATS and returns the connection. If an error occurs, ErrConnect is returned.
func (b ConnectionBuilder) Build() (ConnectionInterface, error) {
	config := b.config
	jsOptions := []nats.Option{
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(config.MaxReconnects),
		nats.ReconnectWait(config.ReconnectWait),
		nats.Name("Kyma Controller"),
	}
	conn, err := nats.Connect(config.URL, jsOptions...)
	if err != nil || !conn.IsConnected() {
		return nil, errors.MakeError(ErrConnect, err)
	}

	return conn, nil
}

// ConnectionInterface is a contract for a NATS connection object.
type ConnectionInterface interface {
	IsConnected() bool
	SetClosedHandler(cb nats.ConnHandler)
	SetReconnectHandler(rcb nats.ConnHandler)
	JetStream(opts ...nats.JSOpt) (nats.JetStreamContext, error)
}
