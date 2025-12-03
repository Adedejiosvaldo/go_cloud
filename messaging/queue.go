package messaging

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

type PubSubQueue struct {
	client       *pubsub.Client
	log          *zap.Logger
	topicName    string
	subName      string
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
}

type NewPubSubQueueOptions struct {
	ProjectID string
	Topic     string
	Sub       string
	Log       *zap.Logger
}

func NewPubSubQueue(ctx context.Context, opts NewPubSubQueueOptions) (*PubSubQueue, error) {
	client, err := pubsub.NewClient(ctx, opts.ProjectID)
	if err != nil {
		return nil, err
	}

	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	q := &PubSubQueue{
		client:    client,
		log:       opts.Log,
		topicName: opts.Topic,
		subName:   opts.Sub,
	}

	// Make sure topic exists
	topic := client.Topic(opts.Topic)
	exists, err := topic.Exists(ctx)
	if err != nil {
		topic, err = client.CreateTopic(ctx, opts.Topic)
		if err != nil {
			return nil, err
		}
	}
	q.topic = topic

	//
	sub := client.Subscription(opts.Sub)
	exists, err = sub.Exists(ctx)

	if err != nil {
		return nil, err
	}

	if !exists {
		sub, err = client.CreateSubscription(ctx, opts.Sub, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}
	q.subscription = sub
	return q, nil
}
