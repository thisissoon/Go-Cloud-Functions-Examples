package addPerson

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
)

// GCLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GCLOUD_PROJECT")

var client *firestore.Client

func init() {
	c, err := firestore.NewClient(
		context.Background(),
		projectID,
	)
	if err != nil {
		log.Fatalf("error creating firestore client: %v", err)
	}
	client = c
}

func AddPerson(w http.ResponseWriter, r *http.Request) {
	decode := json.NewDecoder(r.Body)
	var p Person
	err := decode.Decode(&p)
	if err != nil {
		log.Printf("erroring decoding json: %v", err)
	}
	dbP, err := p.ToDBPerson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error creating person: %v", err)
		return
	}
	ref, err := dbP.SavePerson(r.Context(), client)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error saving user: %v", err)
	} else {
		fmt.Fprintf(w, "successfully add user id: %v", ref.ID)
	}
}
