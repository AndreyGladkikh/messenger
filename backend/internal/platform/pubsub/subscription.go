package pubsub

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Subscription struct {
	pubsub         *PubSub
	subscriberName string
	ch             chan *redis.Message

	channelsMu sync.Mutex
	channels   map[string]struct{}
}

func (s *Subscription) Read(ctx context.Context) (*redis.Message, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case m := <-s.ch:
			return m, nil
		}
	}
}

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
	close(s.ch)
	s.pubsub.removeSubscription(ctx, s)
}
