package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/test-lambda/application/domain/model"
)

type UserAsyncTask interface {
	SaveToS3(user model.User) error
}

func NewUserAsyncTask() (UserAsyncTask, error) {

	if os.Getenv("SQS_ENDPOINT") == "" {
		return nil, fmt.Errorf("SQS_ENDPOINT is not set")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           os.Getenv("SQS_ENDPOINT"),
			SigningRegion: os.Getenv("AWS_REGION"),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, err
	}

	return &userAsyncTask{
		client: sqs.NewFromConfig(cfg),
	}, err
}

type userAsyncTask struct {
	client *sqs.Client
}

func (u *userAsyncTask) SaveToS3(user model.User) error {

	if os.Getenv("USER_SQS_URL") == "" {
		return fmt.Errorf("USER_SQS_URL is not set")
	}

	type userMessage struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		BirthDate string `json:"birth_date"`
	}

	message, err := json.Marshal(userMessage{
		ID:        user.ID,
		Name:      user.Name,
		BirthDate: user.BirthDate.Format("2006-01-02"),
	})

	if err != nil {
		return fmt.Errorf("failed to marshal user message, %w", err)
	}

	_, err = u.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(string(message)),
		QueueUrl:    aws.String(os.Getenv("USER_SQS_URL")),
	})

	if err != nil {
		return fmt.Errorf("failed to send message to sqs, %w", err)
	}

	return nil
}
