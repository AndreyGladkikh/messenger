package rabbitmq

import (
	"encoding/json/jsontext"
	"fmt"
	"messenger/messenger/internal/platform/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QNameRealtimePersonalDelivery  = "realtime.personal-delivery"
	QNameRealtimeBroadcastDelivery = "realtime.broadcast-delivery"
)

func NewConnection(cfg config.Config) (*amqp.Connection, func(), error) {
	conn, err := amqp.Dial(cfg.Mq.Dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	cleanup := func() {
		conn.Close()
	}

	return conn, cleanup, nil
}

type Envelope struct {
	Type    string
	Payload jsontext.Value
}
