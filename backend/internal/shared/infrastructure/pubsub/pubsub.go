package pubsub

import (
	"context"
	"errors"
	"fmt"
	"messenger/messenger/internal/platform/logger"
	"sync"

	"github.com/redis/go-redis/v9"
)

const (
	MessageTypeMessageSent = "messageSent"
)

var (
	ErrSubscriptionExists = errors.New("subscription already exists")
)

type PubSub struct {
	logger      *logger.Logger
	redisClient *redis.Client
	redisPubSub *redis.PubSub

	subscriptionsMu sync.Mutex
	subscriptions   map[string]*Subscription

	topicsMu sync.Mutex
	topics   map[string]map[string]*Subscription
}

func NewPubSub(
	ctx context.Context,
	redisClient *redis.Client,
) *PubSub {
	// ps := redisClient.Subscribe(ctx)

	return &PubSub{
		redisClient:   redisClient,
		// redisPubSub:   ps,
		subscriptions: make(map[string]*Subscription),
		topics: make(map[string]map[string]*Subscription),
	}
}

func (ps *PubSub) Run(ctx context.Context) error {
	if ps.redisPubSub != nil {
		return errors.New("pubsub is already running")
	}

	ps.redisPubSub = ps.redisClient.Subscribe(ctx)
	defer func() {
		ps.redisPubSub.Close()
		ps.redisPubSub = nil
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m := <-ps.redisPubSub.Channel():
			ps.dispatch(m)
		}
	}
}

func (ps *PubSub) dispatch(m *redis.Message) {
	// var message Message
	// if err := json.Unmarshal([]byte(m.Payload), &message); err != nil {
	// 	ps.logger.ErrorContext(context.Background(), "pubsub: failed to parse message", "error", err.Error(), "channel", m.Channel, "payload", m.Payload)
	// 	return
	// }

	// switch message.Type {
	// case MessageTypeEvent:
	// 	var eventEnvelope event.Envelope
	// 	if err := json.Unmarshal(message.Payload, &eventEnvelope); err != nil {
	// 		ps.logger.ErrorContext(context.Background(), "pubsub: failed to parse event", "error", err.Error(), "channel", m.Channel, "payload", m.Payload)
	// 		return
	// 	}
	// }

	// if topics, exists := ps.topics[m.Channel]; exists {
	// 	for _, subscription := range topics {
	// 		subscription.ch <- message
	// 	}
	// }

	if subscriptions, exists := ps.topics[m.Channel]; exists {
		for _, subscription := range subscriptions {
			subscription.ch <- m
		}
	}
}

func (ps *PubSub) Publish(ctx context.Context, channel string, message any) {
	ps.redisClient.Publish(ctx, channel, message)
}

func (ps *PubSub) NewSubscription(subscriberName string) (*Subscription, error) {
	if _, exists := ps.subscriptions[subscriberName]; exists {
		return nil, fmt.Errorf("failed to instantiate pubsub subscription with name %s: %w", subscriberName, ErrSubscriptionExists)
	}

	subscription := &Subscription{
		pubsub: ps,
		subscriberName: subscriberName,
		ch:             make(chan *redis.Message),
	}

	ps.subscriptionsMu.Lock()
	defer ps.subscriptionsMu.Unlock()

	ps.subscriptions[subscriberName] = subscription

	return subscription, nil
}

func (ps *PubSub) removeSubscription(ctx context.Context, s *Subscription) {
	ps.subscriptionsMu.Lock()
	delete(ps.subscriptions, s.subscriberName)
	ps.subscriptionsMu.Unlock()

	var topicsToUnsubscribe []string
	ps.topicsMu.Lock()
	for t := range s.topics {
		if _, exists := ps.topics[t]; exists {
			delete(ps.topics[t], s.subscriberName)
		}
		if len(ps.topics[t]) == 0 {
			delete(ps.topics, t)
			topicsToUnsubscribe = append(topicsToUnsubscribe, t)
		}
	}
	ps.topicsMu.Unlock()

	if len(topicsToUnsubscribe) > 0 {
		ps.redisPubSub.Unsubscribe(ctx, topicsToUnsubscribe...)
	}
}

func (ps *PubSub) Subscribe(ctx context.Context, subscriberName string, topic ...string) (*Subscription, error) {
	if _, exists := ps.subscriptions[subscriberName]; !exists {
		_, err := ps.NewSubscription(subscriberName)
		if err != nil {
			return nil, err
		}
	}
	subscription := ps.subscriptions[subscriberName]

	subscription.AddTopics(ctx, topic...)

	return subscription, nil
}