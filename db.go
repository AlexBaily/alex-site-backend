package main

import (
	"fmt"
	"log"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//Record struct that will house the DynamoDB records.
type Record struct {
        UserID       string
        ExerciseName string
        ExerciseDate string
        Weight       int
        Reps         int
}

func queryTable(UserID string, exrtable string)(queryJson []byte) {
        sess := session.Must(session.NewSessionWithOptions(session.Options{
                SharedConfigState: session.SharedConfigEnable,
        }))

        svc := dynamodb.New(sess)

        //Filter and Projection are required for the expression builder.
        filter := expression.Name("UserID").Equal(expression.Value(UserID))
        projection := expression.NamesList(
                expression.Name("ExerciseName"),
                expression.Name("ExerciseDate"),
                expression.Name("Weight"),
                expression.Name("Reps"),
        )
        expr, err := expression.NewBuilder().
                WithFilter(filter).
                WithProjection(projection).
                Build()
        if err != nil {
                fmt.Println(err)
        }
        //Load up the parameters into a struct
        params := &dynamodb.ScanInput{
                ExpressionAttributeNames:  expr.Names(),
                ExpressionAttributeValues: expr.Values(),
                FilterExpression:          expr.Filter(),
                ProjectionExpression:      expr.Projection(),
                TableName:                 aws.String(exrtable),
        }

        //Complete a scan of the table with the params from above
        result, err := svc.Scan(params)
        if err != nil {
                fmt.Println(err)
        }
        //Used to check what keys are returned from the table scan
        /* 
        for keys := range tablesMap {
                tables += keys
        }*/

        recs := []Record{}

        err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &recs)
        if err != nil {
                panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
        }

        //Marshal the records into JSON
        queryJson, err = json.Marshal(recs[0])
        if err != nil {
                panic(fmt.Sprintf("failed to marshal records, %v", err))
        }
        log.Printf("records %+v", recs[0])
        return

}
