package pubsub

import (
	"context"
	"errors"
	"messenger/messenger/internal/platform/logger"
	"sync"

	"github.com/redis/go-redis/v9"
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

	channelsMu sync.Mutex
	channels   map[string]map[string]*Subscription
}

func NewPubSub(
	redisClient *redis.Client,
) (*PubSub, func()) {
	redisPS := redisClient.Subscribe(context.Background())

	ps := &PubSub{
		redisClient:   redisClient,
		redisPubSub: redisPS,
		subscriptions: make(map[string]*Subscription),
		channels:      make(map[string]map[string]*Subscription),
	}

	go ps.run()

	return ps, func() {
		redisPS.Close()
	}
}

func (ps *PubSub) run() {
	defer ps.redisPubSub.Close()

	for m := range ps.redisPubSub.Channel() {
		ps.dispatch(m)
	}

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return ctx.Err()
	// 	case m := <-ps.redisPubSub.Channel():
	// 		ps.dispatch(m)
	// 	}
	// }
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

	if subscriptions, exists := ps.channels[m.Channel]; exists {
		for _, subscription := range subscriptions {
			subscription.ch <- m
		}
	}
}

func (ps *PubSub) Publish(ctx context.Context, channel string, message any) {
	ps.redisClient.Publish(ctx, channel, message)
}

func (ps *PubSub) Subscription(subscriberName string) *Subscription {
	if subscription, exists := ps.subscriptions[subscriberName]; exists {
		return subscription
	}

	subscription := &Subscription{
		pubsub:         ps,
		subscriberName: subscriberName,
		ch:             make(chan *redis.Message),
		channels: make(map[string]struct{}),
	}

	ps.subscriptionsMu.Lock()
	defer ps.subscriptionsMu.Unlock()

	ps.subscriptions[subscriberName] = subscription

	return subscription
}

func (ps *PubSub) removeSubscription(ctx context.Context, s *Subscription) {
	ps.subscriptionsMu.Lock()
	delete(ps.subscriptions, s.subscriberName)
	ps.subscriptionsMu.Unlock()

	var channelsToUnsubscribe []string
	ps.channelsMu.Lock()
	for c := range s.channels {
		if _, exists := ps.channels[c]; exists {
			delete(ps.channels[c], s.subscriberName)
		}
		if len(ps.channels[c]) == 0 {
			delete(ps.channels, c)
			channelsToUnsubscribe = append(channelsToUnsubscribe, c)
		}
	}
	ps.channelsMu.Unlock()

	if len(channelsToUnsubscribe) > 0 {
		ps.redisPubSub.Unsubscribe(ctx, channelsToUnsubscribe...)
	}
}

func (ps *PubSub) Subscribe(ctx context.Context, subscriberName string, channel ...string) *Subscription {
	subscription := ps.Subscription(subscriberName)
	subscription.AddChannels(ctx, channel...)

	return subscription
}
