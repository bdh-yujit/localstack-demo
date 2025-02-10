package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/test-async-lambda/application/domain/model"
)

type S3Storage interface {
	SaveUsers(users []model.User) error
}

func NewS3Storage() (S3Storage, error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           os.Getenv("S3_ENDPOINT"),
			SigningRegion: os.Getenv("AWS_REGION"),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, err
	}

	return &s3Storage{
		client: s3.NewFromConfig(cfg),
	}, nil
}

type s3Storage struct {
	client *s3.Client
}

func (s *s3Storage) SaveUsers(users []model.User) error {

	if os.Getenv("USER_BUCKET_NAME") == "" {
		return fmt.Errorf("USER_BUCKET_NAME is not set")
	}

	fileName := fmt.Sprintf("%s_users.json", time.Now().Format("2006-01-02T15:04:05.000Z"))

	body, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal users, %w", err)
	}

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("USER_BUCKET_NAME")),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(body),
	})

	return err
}
