package mq

type ClientInterface interface {
	Subscribe(channel string, topic string, onmessage SubscriptionCallback) error
	Publish(channel string, topic string, payload interface{}) error
}
