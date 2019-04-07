package main

import (
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

)

var (
	exrtable string = os.Getenv("EXRTABLE")
)

type Record struct {
            Exercise   string
}

//Http handler for responding to http/s requests.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	//Set response headers.
	w.Header().Add("statusDescription", "200 OK")
	w.Header().Set("statusDescription", "200 OK")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fullScan()))
}


func fullScan()(tables []byte) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
	    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	//Filter and Projection are required for the expression builder.
	filter := expression.Name("UserID").Equal(expression.Value("1"))
	projection := expression.NamesList(
		expression.Name("WorkoutTimeStamp"),
		expression.Name("Exercise"),
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
	tables, err = json.Marshal(recs[0])
	if err != nil {
	 panic(fmt.Sprintf("failed to marshal records, %v", err))
	}
	//log.Printf("records %+v", result)
	//log.Printf("records %+v", recs[0].Exercise)
	return

}

func main() {
	//Create a new mux router.
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	fullScan()
	log.Fatal(http.ListenAndServe(":8080", mux))
}
