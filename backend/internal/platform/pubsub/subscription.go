package pubsub

import (
	"context"
	"encoding/json/v2"
	"sync"
	"uuid"

	"github.com/redis/go-redis/v9"
)

type Subscription struct {
	pubsub         *PubSub
	subscriberName string
	ch             chan *redis.Message

	channelsMu sync.Mutex
	channels   map[string]struct{}
}

// func (s *Subscription) Read(ctx context.Context) (*redis.Message, error) {
// 	select {
// 	case <-ctx.Done():
// 		return nil, ctx.Err()
// 	case m := <-s.ch:
// 		return m, nil
// 	}
// }

func (s *Subscription) Read(ctx context.Context) (Message, error) {
	var msg Message
	select {
	case <-ctx.Done():
		return msg, ctx.Err()
	case m := <-s.ch:
		err := json.Unmarshal([]byte(m.Payload), &msg)
		if err != nil {
			return msg, err
		}

		return msg, nil
	}
}

// func (s *Subscription) ReadAs[T any](ctx context.Context) (T, error) {
// 	var t T

// 	m, err := s.Read(ctx)
// 	if err != nil {
// 		return t, err
// 	}

// 	// todo unmarshal RawEnvelope
// 	err = json.Unmarshal(jsontext.Value(m.Payload), &t)
// 	if err != nil {
// 		return t, err
// 	}

// 	return t, nil
// }

func (s *Subscription) AddChannels(ctx context.Context, channels ...string) {
	s.channelsMu.Lock()
	s.pubsub.channelsMu.Lock()

	var newChannels []string
	for _, c := range channels {
		if _, exists := s.pubsub.channels[c]; !exists {
			s.pubsub.channels[c] = make(map[string]*Subscription)
		}
		if len(s.pubsub.channels[c]) == 0 {
			newChannels = append(newChannels, c)
		}
		s.pubsub.channels[c][s.subscriberName] = s
		s.channels[c] = struct{}{}
	}

	s.channelsMu.Unlock()
	s.pubsub.channelsMu.Unlock()

	if len(newChannels) > 0 {
		s.pubsub.redisPubSub.Subscribe(ctx, newChannels...)
	}
}

func (s *Subscription) Cancel(ctx context.Context) {
	s.pubsub.removeSubscription(ctx, s)
}

func SubscribeToUserEvents(ctx context.Context, subscription *Subscription, userID uuid.UUID) {
	subscription.AddChannels(ctx, ChannelNameUserEvents(userID))
}

func SubscribeToConversationEvents(ctx context.Context, subscription *Subscription, cnvIDs ...uuid.UUID) {
	channels := make([]string, 0, len(cnvIDs))
	for _, cid := range cnvIDs {
		channels = append(channels, ChannelNameConversationEvents(cid))
	}
	subscription.AddChannels(ctx, channels...)
}
