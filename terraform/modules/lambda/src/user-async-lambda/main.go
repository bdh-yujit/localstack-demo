package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/test-async-lambda/application/usecase"
)

func main() {
	lambda.Start(ErrorHandler(HandleRequest))
}

type Payload struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	BirthDate string `json:"birth_date"`
}

type handler func(ctx context.Context, p events.SQSEvent) (events.SQSEventResponse, error)

func HandleRequest(ctx context.Context, p events.SQSEvent) (events.SQSEventResponse, error) {

	fmt.Println("Request: ", p)
	var input usecase.ExecuteInput
	input.Users = make([]usecase.ExecuteInputUser, 0, len(p.Records))
	for _, record := range p.Records {
		p := Payload{}
		err := json.Unmarshal([]byte(record.Body), &p)
		if err != nil {
			return events.SQSEventResponse{}, err
		}
		input.Users = append(input.Users, usecase.ExecuteInputUser{
			Name:      p.Name,
			BirthDate: p.BirthDate,
			ID:        p.ID,
		})
	}

	u, err := usecase.NewUsecase()
	if err != nil {
		return events.SQSEventResponse{}, err
	}
	err = u.Execute(input)
	if err != nil {
		return events.SQSEventResponse{}, err
	}

	return events.SQSEventResponse{}, nil
}

func ErrorHandler(h handler) handler {
	return func(ctx context.Context, p events.SQSEvent) (events.SQSEventResponse, error) {
		res, err := h(ctx, p)
		if err != nil {
			fmt.Println(err)
		}
		return res, err
	}
}
