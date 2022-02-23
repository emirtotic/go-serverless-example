package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/emirtotic/go-serverless-example/pkg/handlers"
	"os"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {

	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)})

	if err != nil {
		return
	}
	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)

}

const tableName = "go-serverless-example"

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	switch request.HTTPMethod {
	case "GET":
		return handlers.GetUser(request, tableName, dynaClient)
	case "POST":
		return handlers.CreateUser(request, tableName, dynaClient)
	case "PUT":
		return handlers.UpdateUser(request, tableName, dynaClient)
	case "DELETE":
		return handlers.DeleteUser(request, tableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}
}
