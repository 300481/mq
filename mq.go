package mq

import (
	"context"
	"errors"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

const (
	errGCPTopicDontExist        = "GCP PubSub Topic not existing and not allowed to create one."
	errGCPSubscriptionDontExist = "GCP PubSub Subscription not existing and not allowed to create one."
)

type GCP struct {
	CredentialsFile    string
	TopicName          string
	CreateTopic        bool
	SubscriptionName   string
	CreateSubscription bool
	ProjectID          string
}

// NewGCP creates new GCP PubSub struct
func NewGCP() *GCP {
	return &GCP{
		CredentialsFile:    os.Getenv("GCP_CREDENTIALS_FILE"),
		TopicName:          os.Getenv("GCP_TOPIC_NAME"),
		CreateTopic:        os.Getenv("GCP_CREATE_TOPIC") == "TRUE",
		SubscriptionName:   os.Getenv("GCP_SUBSCRIPTION_NAME"),
		CreateSubscription: os.Getenv("GCP_CREATE_SUBSCRIPTION") == "TRUE",
		ProjectID:          os.Getenv("GCP_PROJECT_ID"),
	}
}

/*
Publish publishes a message on GCP PubSub
Config needed:
	"GCP_PROJECT_ID"
	"GCP_CREDENTIALS_FILE"
	"GCP_TOPIC_NAME"
	"GCP_CREATE_TOPIC"
*/
func (m *GCP) Publish(payload []byte) (id string, err error) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, m.ProjectID)
	if err != nil {
		return "", err
	}

	topic, err := m.createTopicIfNotExists(client, ctx)
	if err != nil {
		return "", err
	}

	result := topic.Publish(ctx, &pubsub.Message{
		Data: payload,
	})

	id, err = result.Get(ctx)
	if err != nil {
		return "", err
	}

	return id, nil
}

/*
Subscribe subscribes to an existing subscription
Config needed:
	"GCP_PROJECT_ID"
	"GCP_CREDENTIALS_FILE"
	"GCP_TOPIC_NAME"
	"GCP_CREATE_TOPIC"
	"GCP_SUBSCRIPTION_NAME"
	"GCP_CREATE_SUBSCRIPTION"
*/
func (m *GCP) Subscribe(handleFunc func(ctx context.Context, m *pubsub.Message)) (err error) {
	ctx := context.Background()

	opts := option.WithCredentialsFile(m.CredentialsFile)

	client, err := pubsub.NewClient(ctx, m.ProjectID, opts)
	if err != nil {
		return err
	}

	sub, err := m.createSubscriptionIfNotExists(client, ctx)
	if err != nil {
		return err
	}

	err = sub.Receive(ctx, handleFunc)
	if err != nil {
		return err
	}

	return nil
}

/*
createTopicIfNotExists creates a Topic if its not existing
and allowed to create one
*/
func (m *GCP) createTopicIfNotExists(client *pubsub.Client, ctx context.Context) (topic *pubsub.Topic, err error) {
	topic = client.Topic(m.TopicName)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		if m.CreateTopic {
			topic, err = client.CreateTopic(ctx, m.TopicName)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New(errGCPTopicDontExist)
		}
	}
	return topic, nil
}

/*
createSubscriptionIfNotExists creates a Subscription if its not existing
and allowed to create one
*/
func (m *GCP) createSubscriptionIfNotExists(client *pubsub.Client, ctx context.Context) (topic *pubsub.Subscription, err error) {
	sub := client.Subscription(m.SubscriptionName)
	ok, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		if m.CreateSubscription {
			topic, err := m.createTopicIfNotExists(client, ctx)
			if err != nil {
				return nil, err
			}
			sub, err = client.CreateSubscription(
				ctx,
				m.SubscriptionName,
				pubsub.SubscriptionConfig{
					Topic:       topic,
					AckDeadline: 60 * time.Second,
				},
			)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New(errGCPSubscriptionDontExist)
		}
	}
	return sub, nil
}
