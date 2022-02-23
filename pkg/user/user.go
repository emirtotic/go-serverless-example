package user

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/emirtotic/go-serverless-example/pkg/validators"
)

var (
	ErrorFailedToFetchUser       = "Failed to fetch record!"
	ErrorFailedToUnmarshalRecord = "Failed to unmarshal record!"
	ErrorInvalidUserData         = "Error Invalid User Data"
	ErrorInvalidEmail            = "Error Invalid Email"
	ErrorFailedToMarshalRecord   = "Failed to Marshal record"
	ErrorCouldNotDeleteItem      = "Failed to Delete Item"
	ErrorCouldNotDynamoPutItem   = "Failed to dynamo put Item"
	ErrorUserAlreadyExist        = "Error User Already Exist"
	ErrorUserDoesNotExist        = "Error User Does Not Exist"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FetchUser(email string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := dynaClient.GetItem(input)

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchUser)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)

	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil

}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	res, err := dynaClient.Scan(input)

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchUser)
	}
	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, item)

	return item, err

}

func CreateUser(request events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {

	var u User

	if err := json.Unmarshal([]byte(request.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	//Checking if the user already exist

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)

	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExist)
	}

	av, err := dynamodbattribute.MarshalMap(u)

	if err != nil {
		return nil, errors.New(ErrorFailedToMarshalRecord)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)

	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &u, nil

}

func UpdateUser(request events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {

	var u User

	if err := json.Unmarshal([]byte(request.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient)

	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExist)
	}

	av, err := dynamodbattribute.MarshalMap(u)

	if err != nil {
		return nil, errors.New(ErrorFailedToMarshalRecord)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)

	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &u, nil

}

func DeleteUser(request events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {

	email := request.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},

		TableName: aws.String(tableName),
	}

	_, err := dynaClient.DeleteItem(input)

	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil
}
