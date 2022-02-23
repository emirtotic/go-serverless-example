package handlers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/emirtotic/go-serverless-example/pkg/user"
	"net/http"
)

var ErrorMethodNotAllowed = "Method not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error_msg,omitempty"`
}

func GetUser(request events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {

	email := request.QueryStringParameters["email"]

	if len(email) > 0 {
		result, err := user.FetchUser(email, tableName, dynaClient)

		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}

	result, err := user.FetchUsers(tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusOK, result)

}

func CreateUser(request events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {

	result, err := user.CreateUser(request, tableName, dynaClient)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusCreated, result)
}

func UpdateUser(request events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {

	result, err := user.UpdateUser(request, tableName, dynaClient)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusOK, result)

}

func DeleteUser(request events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {

	err := user.DeleteUser(request, tableName, dynaClient)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return apiResponse(http.StatusOK, nil)

}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
