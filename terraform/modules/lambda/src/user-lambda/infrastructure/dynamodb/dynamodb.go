package dynamodb

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/test-lambda/application/domain/model"
	"github.com/test-lambda/application/domain/repository"
)

type UserRepositoryImpl struct {
	client *dynamodb.Client
}

func NewUserRepository() (repository.UserRepository, error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {

		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           os.Getenv("DYNAMODB_ENDPOINT"),
			SigningRegion: os.Getenv("AWS_REGION"),
		}, nil
	})
	dynamoDBConf, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, err
	}

	return &UserRepositoryImpl{
		client: dynamodb.NewFromConfig(dynamoDBConf),
	}, nil
}

type UserTable struct {
	ID        string `dynamodbav:"id"`
	Name      string `dynamodbav:"name"`
	BirthDate string `dynamodbav:"birth_date"`
}

func (r *UserRepositoryImpl) Save(user model.User) error {
	item := UserTable{
		ID:        user.ID,
		Name:      user.Name,
		BirthDate: user.BirthDate.String(),
	}
	av, marshalErr := attributevalue.MarshalMap(item)
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal user item, %w", marshalErr)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE_NAME")),
	}

	_, err := r.client.PutItem(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to put item, %w", err)
	}

	return nil
}
