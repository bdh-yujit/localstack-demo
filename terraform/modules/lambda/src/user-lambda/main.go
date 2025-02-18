package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/test-lambda/adapter"
	"github.com/test-lambda/application/usecase"
	"github.com/test-lambda/infrastructure/mysql"
)

type CreateUserPayload struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

type GetUserPayload struct {
	Id int64 `json:"id"`
}

func init() {
	err := mysql.InitDBConn()
	if err != nil {
		panic(err)
	}
}

func main() {
	lambda.Start(ErrorHandler(HandleRequest))
}

func NewInvalidRequestError(msg string, err error) *InvalidRequestError {
	return &InvalidRequestError{
		Message: msg,
		Err:     err,
	}
}

type InvalidRequestError struct {
	Message string
	Err     error
}

func (e *InvalidRequestError) Error() string {
	return e.Message
}

func (e *InvalidRequestError) Unwrap() error {
	return e.Err
}

var invalidRequestError *InvalidRequestError

type handler func(ctx context.Context, p events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

func HandleRequest(ctx context.Context, p events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Request: %s\n", string(jsonBytes))

	switch p.HTTPMethod {
	case "POST":
		payload := CreateUserPayload{}
		err := json.Unmarshal([]byte(p.Body), &payload)
		if err != nil {
			return nil, NewInvalidRequestError(fmt.Sprintf("Failed to unmarshal request body: %s", err.Error()), err)
		}

		if payload.Name == "" {
			return nil, NewInvalidRequestError("Name is required", nil)
		}
		if payload.BirthDate == "" {
			return nil, NewInvalidRequestError("BirthDate is required", nil)
		}

		repo := adapter.NewUserRepository()

		service, err := adapter.NewUserProcessServiceAdapter()
		if err != nil {
			return nil, err
		}

		output, err := usecase.NewUsecase(repo, service).CreateUser(usecase.CreateUserInput{
			Name:      payload.Name,
			BirthDate: payload.BirthDate,
		})

		if err != nil {
			return nil, err
		}

		respBodyBytes, err := json.Marshal(output)
		if err != nil {
			return nil, err
		}

		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(respBodyBytes),
		}, nil
	case "GET":
		id := p.QueryStringParameters["id"]
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, NewInvalidRequestError("Invalid id", err)
		}
		output, err := usecase.NewUsecase(adapter.NewUserRepository(), nil).GetUser(idInt)
		if err != nil {
			return nil, err
		}

		respBodyBytes, err := json.Marshal(output)
		if err != nil {
			return nil, err
		}

		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(respBodyBytes),
		}, nil

	}

	return nil, NewInvalidRequestError(fmt.Sprintf("Unsupported HTTP method: %s", p.HTTPMethod), nil)
}

func ErrorHandler(h handler) handler {
	return func(ctx context.Context, p events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		res, err := h(ctx, p)
		if err != nil {
			fmt.Println(err)

			switch {
			case errors.As(err, &invalidRequestError):
				return &events.APIGatewayProxyResponse{
					StatusCode: 400,
					Body:       fmt.Sprintf("Bad Request: %s", err.Error()),
				}, nil
			}

			return &events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Internal Server Error",
			}, nil
		}
		return res, err
	}
}
