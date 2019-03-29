package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//Http handler for responding to http/s requests.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	//Set response headers.
	w.Header().Add("statusDescription", "200 OK")
	w.Header().Set("statusDescription", "200 OK")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	URISegments := strings.Split(r.URL.Path, "/")
	w.Write([]byte(URISegments[1]))
}

func fullScan() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
	    SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	// Get the list of tables 
	//Lifted from the AWS docs - This just gets a list of tables
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
	 fmt.Println(err)
	}

	fmt.Println("Tables:")
	fmt.Println("")

	for _, n := range result.TableNames {
	    fmt.Println(*n)
	}

}

func main() {
	//Create a new mux router.
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
