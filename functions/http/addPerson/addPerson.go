package addPerson

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github/thisissoon/Go-Cloud-Functions-Examples/http/addPerson/storage"

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

func runAddPerson(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	decode := json.NewDecoder(r.Body)
	var p Person
	err := decode.Decode(&p)
	if err != nil {
		log.Printf("erroring decoding json: %v", err)
	}
	dbP, err := p.ToDBPerson()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "incorrect time format use: 2006-01-02T15:04:05.000Z")
		return
	}

	ref, err := s.SaveDoc(r.Context(), "people", dbP)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error saving user: %v", err)
	} else {
		fmt.Fprintf(w, "successfully added user. Id: %v", ref.ID)
	}
}

func AddPerson(w http.ResponseWriter, r *http.Request) {
	store := storage.Store{
		Client: client,
	}
	runAddPerson(w, r, store)
}
