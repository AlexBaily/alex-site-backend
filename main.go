package main

import (
	"os"
	"log"
	"fmt"
	"strings"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

)

//"Global" variable for the exerise table.
var (
	exrtable string = os.Getenv("EXRTABLE")
)

//Record struct that will house the DynamoDB records.
type Record struct {
	UserID       string
	ExerciseName string
	ExerciseDate string
	Weight       int
	Reps         int
}

//Middleware to read the Authorization header for the Cognito JWT token
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Don't bother checking the Auth header if we are just going to route
		if r.URL.Path == "/" {
			next.ServeHTTP(w,r)
		} else {
			//Get the token
			token := r.Header.Get("Authorization")
			//log.Printf("token %+v", token)
			jwtToken := strings.Split(token, " ")
			//Check is the jwtToken contains an actual token
			if len(jwtToken) <= 1 {
				//Return a 403 if no token is found
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				log.Printf("jwt %+v", jwtToken[1])
				jwtTokenArray := strings.Split(jwtToken[1], ".")
				//Checks to make sure that the jwtTokenArray has 3 parts.
				//Thsi will be the header, payload and signature of the jwt token
				if len(jwtTokenArray) <= 2 {
					http.Error(w, "Forbidden", http.StatusForbidden)
				} else {
					next.ServeHTTP(w,r)
				}
			}
		}
	})
}

//Http handler for responding to http/s requests.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	//Set response headers.
	w.Header().Add("statusDescription", "200 OK")
	w.Header().Set("statusDescription", "200 OK")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("{\"Route\":\"Test\"}"))
}

func exerciseHandler(w http.ResponseWriter, r *http.Request) {
        //Set response headers.
        w.Header().Add("statusDescription", "200 OK")
        w.Header().Set("statusDescription", "200 OK")
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
	log.Printf("records %+v", queryTable("0000"))
	//retrieve the UserID variable

        w.Write(queryTable("0000"))
}

func queryTable(UserID string)(queryJson []byte) {
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

func main() {
	//Create a new mux router.
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/exercises", exerciseHandler)
	r.Use(authMiddleware)
	log.Fatal(http.ListenAndServe(":8080", r))
}
