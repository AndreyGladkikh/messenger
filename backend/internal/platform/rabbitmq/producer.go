package rabbitmq

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"log"
	"messenger/messenger/internal/platform/logger"
	"messenger/messenger/internal/platform/rabbitmq/exchanges"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	logger  *logger.Logger
	rmqChan *amqp.Channel
	xName   string
}

func NewProducer(
	logger *logger.Logger,
	amqpConn *amqp.Connection,
	xName string,
) (*Producer, func(), error) {
	ch, err := amqpConn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("rabbitmq producer: failed to open a channel: %w", err)
	}
	cleanup := func() {
		ch.Close()
	}

	err = exchanges.Declare(xName, ch)
	if err != nil {
		return nil, cleanup, err
	}

	return &Producer{
		logger:  logger,
		rmqChan: ch,
		xName:   xName,
	}, cleanup, nil
}

func (p *Producer) Produce(ctx context.Context, typ string, payload any) error {
	var (
		jsonPayload jsontext.Value
		err         error
	)

	logErr := func(logMsg string, err error) {
		p.logger.ErrorContext(ctx, fmt.Sprintf("rabbitmq producer: %s", logMsg),
			"err", err,
			"xName", p.xName,
			"typ", typ,
			"payload", payload,
		)
	}

	b, ok := payload.([]byte)
	if ok {
		jsonPayload = jsontext.Value(b)
		if !jsonPayload.IsValid() {
			logErr("payload is not a valid json", nil)
			return fmt.Errorf("rabbitmq producer: payload is not a valid json")
		}
	} else {
		jsonPayload, err = json.Marshal(payload)
		if err != nil {
			logErr("failed to serialize payload", err)
			return fmt.Errorf("rabbitmq producer: failed to serialize payload: %w", err)
		}
	}

	body, err := json.Marshal(Envelope{
		Type:    typ,
		Payload: jsonPayload,
	})
	if err != nil {
		logErr("failed to serialize payload", err)
		return fmt.Errorf("failed to serialize payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.rmqChan.PublishWithContext(ctx,
		p.xName,
		typ,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf(" [x] Sent %s", body)

	return err
}
