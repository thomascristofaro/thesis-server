# Thesis Server
Repository for the server code of my MSc thesis.

Detailed description [here](https://thomascristofaro.github.io/docs/tesi/intro)

## SNS

### Who Publish

Put the resource in serverless file:

```yaml
resources:
  Resources:
    NameSNSResourceTopic:
      Type: AWS::SNS::Topic
      Properties:
        TopicName: NameSNSResource
```

Put the ARN in the environment variable:

```yaml
  environment:
    NameSNSResourceTopicArn: !GetAtt NameSNSResourceTopic.TopicArn
```
You will use the ARN to publish the message.

### Who Subscribe

Create the function with the linked event:

```yaml
functions:
  NameFunction:
    handler: handlers/nameFunction
    events:
      - sns:
          arn: NameSNSResourceTopicArn
```

## SQS

### Who Receive Message

Put the resource in serverless file:

```yaml
resources:
  Resources:
    NameSQSResourceQueue:
      Type: AWS::SQS::Queue
      Properties:
        QueueName: NameSQSResource
```

Create a function with the linked event:

```yaml
functions:
  NameFunction:
    handler: handlers/nameFunction
    events:
      - sqs:
          arn: !GetAtt NameSQSResourceQueue.Arn
```

### Who Send Message

Put the ARN in the environment variable:

```yaml
  environment:
    NameSQSResourceQueueUrl: NameSQSResourceQueueUrl
```