package rabbitmq

import (
	"context"
	"encoding/json/v2"
	"fmt"
	cnvEvnt "messenger/messenger/internal/conversations/app/event"
	"messenger/messenger/internal/shared/app/event"
	"time"

	rmq "github.com/rabbitmq/rabbitmq-amqp-go-client/pkg/rabbitmqamqp"
)

type Producer struct {
	publisher *rmq.Publisher
}

func (p *Producer) Produce(ctx context.Context, payload any) error {
	var msgType string

	switch t := payload.(type) {
	case event.Envelope[cnvEvnt.DirectConversationWithFirstMessageCreated]:
		msgType = t.Event.Name()
	default:
		return fmt.Errorf("rabbitmq producer: failed to infer message type")
	}

	pl, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rabbitmq producer: failed to infer serialize message %q^ %w", msgType, err)
	}

	msg := Message{
		Type:    msgType,
		Payload: pl,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("rabbitmq producer: failed to infer serialize message %q^ %w", msgType, err)
	}

	res, err := p.publisher.Publish(ctx, rmq.NewMessage(body))
	if err != nil {
		return fmt.Errorf("rabbitmq producer: failed to publish a message: %w", err)
	}

	switch res.Outcome.(type) {
	case *rmq.StateAccepted:
	case *rmq.StateRejected:
		return fmt.Errorf("rabbitmq producer: message was rejected: %v", res.Outcome)
	case *rmq.StateReleased:
		return fmt.Errorf("rabbitmq producer: message was released: %v", res.Outcome)
	case *rmq.StateModified:
		return fmt.Errorf("rabbitmq producer: message was modified: %v", res.Outcome)
	default:
		return fmt.Errorf("rabbitmq producer: unexpected publish outcome: %v", res.Outcome)
	}

	return nil
}

type ConversationEventsProducer struct {
	*Producer
}

func NewConversationEventsProducer(rmqConn *rmq.AmqpConnection) (*ConversationEventsProducer, func(), error) {
	xName := "conversation.events"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rmqConn.Management().DeclareExchange(ctx, &rmq.TopicExchangeSpecification{Name: xName})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare an exchange %q: %w", xName, err)
	}

	publisher, err := rmqConn.NewPublisher(ctx, &rmq.ExchangeAddress{
		Exchange: xName,
	}, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create publisher %q: %w", xName, err)
	}
	cleanup := func() { publisher.Close(context.Background()) }

	return &ConversationEventsProducer{
		Producer: &Producer{publisher: publisher},
	}, cleanup, nil
}
