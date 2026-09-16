package pubsub

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

type SubscriptionIface[M any] interface {
	Read(ctx context.Context) (M, error)
	Cancel(ctx context.Context)
	AddTopics(ctx context.Context, topics ...string)
}

type Subscription struct {
	pubsub *PubSub
	subscriberName string
	ch             chan  *redis.Message

	topicsMu sync.Mutex
	topics map[string]struct{}
}

func (s *Subscription) Read(ctx context.Context) (Message, error) {
	for {
		select {
		case <-ctx.Done():
			return Message{}, ctx.Err()
		case m := <-s.ch:
			return m, nil
		}
	}
}

func (s *Subscription) AddTopics(ctx context.Context, topics ...string) {
	s.topicsMu.Lock()
	s.pubsub.topicsMu.Lock()

	var newTopics []string
	for _, t := range topics {
		if _, exists := s.pubsub.topics[t]; !exists {
			s.pubsub.topics[t] = make(map[string]*Subscription)
		}
		if len(s.pubsub.topics[t]) == 0 {
			newTopics = append(newTopics, t)
		}
		s.pubsub.topics[t][s.subscriberName] = s
		s.topics[t] = struct{}{}
	}

	s.topicsMu.Unlock()
	s.pubsub.topicsMu.Unlock()

	if len(newTopics) > 0 {
		s.pubsub.redisPubSub.Subscribe(ctx, newTopics...)
	}
}

func (s *Subscription) Cancel(ctx context.Context) {
	close(s.ch)
	s.pubsub.removeSubscription(ctx,s)
}