package rabbitmq

import (
	"encoding/json/jsontext"
	"fmt"
	"messenger/messenger/internal/platform/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QNameRealtimePersonalDelivery  = "realtime.personal-delivery"
	QNameRealtimeBroadcastDelivery = "realtime.broadcast-delivery"
)

func NewConnection(cfg config.Config) (*amqp.Connection, func(), error) {
	// todo rm
	config := amqp.Config{
		Heartbeat: time.Hour,
	}
	conn, err := amqp.DialConfig(cfg.Mq.Dsn, config)
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
