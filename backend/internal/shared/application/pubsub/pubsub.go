package pubsub

type PubSub interface {
	Publish()
	Subscribe()
}