package utility

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/awssnssqs"
)

type Message struct {
	Body     []byte
	Metadata map[string]string
}

func HandlerSNSWithLogError(ctx context.Context, snsEvent events.SNSEvent, handler func(context.Context, *Message) error) error {
	message, err := ConvertSNSEventToMessage(snsEvent)
	if err != nil {
		SendSQSLogError(ctx, message, err)
		return err
	}
	err = handler(ctx, &message)
	if err != nil {
		SendSQSLogError(ctx, message, err)
		return err
	}
	return err
}

func HandlerSQSWithLogError(ctx context.Context, sqsEvent events.SQSEvent, handler func(context.Context, *Message) error) error {
	message, err := ConvertSQSEventToMessage(sqsEvent)
	if err != nil {
		SendSQSLogError(ctx, message, err)
		return err
	}
	err = handler(ctx, &message)
	if err != nil {
		SendSQSLogError(ctx, message, err)
		return err
	}
	return nil
}

func ConvertSQSEventToMessage(sqsEvent events.SQSEvent) (Message, error) {
	if len(sqsEvent.Records) == 0 {
		return Message{}, errors.New("NO SQS Event")
	}
	if len(sqsEvent.Records) > 1 {
		event, _ := json.Marshal(sqsEvent)
		return Message{}, errors.New("More SQS Event: " + string(event))
	}

	sqsMessage := sqsEvent.Records[0]
	return Message{
		Body:     []byte(sqsMessage.Body),
		Metadata: convertSQSAttributeToMapString(sqsMessage.MessageAttributes),
	}, nil
}

func ConvertSNSEventToMessage(snsEvent events.SNSEvent) (Message, error) {
	if len(snsEvent.Records) == 0 {
		return Message{}, errors.New("NO SNS Event")
	}
	if len(snsEvent.Records) > 1 {
		event, _ := json.Marshal(snsEvent)
		return Message{}, errors.New("More SNS Event: " + string(event))
	}

	snsEntity := snsEvent.Records[0].SNS
	return Message{
		Body:     []byte(snsEntity.Message),
		Metadata: convertSNSAttributeToMapString(snsEntity.MessageAttributes),
	}, nil
}

func convertSNSAttributeToMapString(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		vmap := v.(map[string]interface{})
		if vmap["Type"] == "String" {
			result[k] = vmap["Value"].(string)
		}
	}
	return result
}

func convertSQSAttributeToMapString(m map[string]events.SQSMessageAttribute) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if v.DataType == "String" {
			result[k] = *v.StringValue
		}
	}
	return result
}

func SendSNSMessage(ctx context.Context, env_topic string, message Message) error {
	BuildLogMetadata("SENDING", message.Metadata["function"], env_topic, "SNS", &message)
	if err := SendSQSLog(ctx, message); err != nil {
		return err
	}

	topicARN, ok := os.LookupEnv(env_topic)
	if !ok {
		return errors.New(fmt.Sprintf("Environment variable %s not found", env_topic))
	}

	topic, err := pubsub.OpenTopic(ctx, "awssns:///"+topicARN+"?region=us-east-1")
	if err != nil {
		return err
	}
	defer topic.Shutdown(ctx)

	err = topic.Send(ctx, &pubsub.Message{
		Body: message.Body,
		// Metadata is optional and can be nil.
		Metadata: message.Metadata,
	})
	if err != nil {
		return err
	}

	return nil
}

func SendSQSMessage(ctx context.Context, env_topic string, message Message, fifo bool) error {
	BuildLogMetadata("SENDING", message.Metadata["function"], env_topic, "SQS", &message)
	if err := SendSQSLog(ctx, message); err != nil {
		return err
	}
	return sendSQSMessageWithoutLog(ctx, env_topic, message, fifo)
}

func sendSQSMessageWithoutLog(ctx context.Context, env_topic string, message Message, fifo bool) error {
	queueURL, ok := os.LookupEnv(env_topic)
	if !ok {
		return errors.New(fmt.Sprintf("Environment variable %s not found", env_topic))
	}

	if strings.HasPrefix(queueURL, "https://") {
		queueURL = queueURL[8:]
	}
	topic, err := pubsub.OpenTopic(ctx, "awssqs://"+queueURL+"?region=us-east-1")
	if err != nil {
		return err
	}
	defer topic.Shutdown(ctx)

	var beforeSend func(asFunc func(interface{}) bool) error
	if fifo {
		beforeSend = func(asFunc func(interface{}) bool) error {
			var smi *sqs.SendMessageInput
			if asFunc(&smi) {
				str := "default"
				(*smi).MessageGroupId = &str
			}
			return nil
		}
	}

	err = topic.Send(ctx, &pubsub.Message{
		Body: message.Body,
		// Metadata is optional and can be nil.
		Metadata:   message.Metadata,
		BeforeSend: beforeSend,
	})
	if err != nil {
		return err
	}

	return nil
}

// LOGGING
func BuildLogMetadata(status string, function string, event string, service string, message *Message) {
	message.Metadata["status"] = status
	message.Metadata["function"] = function
	message.Metadata["event"] = event
	message.Metadata["service"] = service
	if uuid_data, ok := message.Metadata["uuid"]; !ok || uuid_data == "" {
		message.Metadata["uuid"] = uuid.New().String()
	}
}

func SendSQSLogError(ctx context.Context, message Message, err error) error {
	message.Metadata["status"] = "ERROR"
	message.Body = []byte(err.Error())
	// fmt.Printf("LOG ERROR: %s", err.Error())
	return SendSQSLog(ctx, message)
}

func SendSQSLog(ctx context.Context, message Message) error {
	if _, ok := os.LookupEnv("LogMessageQueueUrl"); ok {
		body := map[string]interface{}{}
		oldBody := map[string]interface{}{}
		if err := json.Unmarshal(message.Body, &oldBody); err != nil {
			body["body"] = string(message.Body)
		} else {
			body["body"] = oldBody
		}
		body["metadata"] = message.Metadata
		message.Body, _ = json.Marshal(body)
		if err := sendSQSMessageWithoutLog(ctx, "LogMessageQueueUrl", message, true); err != nil {
			return err
		}
	}
	return nil
}
