package utility

import (
	"context"
	"os"
	"strings"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/awssnssqs"
)

func SendSNSMessage(ctx context.Context, env_topic string, message []byte, metadata map[string]string) error {
	topicARN := os.Getenv(env_topic)
	topic, err := pubsub.OpenTopic(ctx, "awssns:///"+topicARN+"?region=us-east-1")
	if err != nil {
		return err
	}
	defer topic.Shutdown(ctx)

	err = topic.Send(ctx, &pubsub.Message{
		Body: message,
		// Metadata is optional and can be nil.
		Metadata: metadata,
	})
	if err != nil {
		return err
	}

	return nil
}

func SendSQSMessage(ctx context.Context, env_topic string, message []byte, metadata map[string]string) error {
	queueURL := os.Getenv(env_topic)
	if strings.HasPrefix(queueURL, "https://") {
		queueURL = queueURL[8:]
	}
	topic, err := pubsub.OpenTopic(ctx, "awssqs://"+queueURL+"?region=us-east-1")
	if err != nil {
		return err
	}
	defer topic.Shutdown(ctx)

	err = topic.Send(ctx, &pubsub.Message{
		Body: message,
		// Metadata is optional and can be nil.
		Metadata: metadata,
	})
	if err != nil {
		return err
	}

	return nil
}
