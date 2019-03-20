package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var (
	projectId             = os.Getenv("GCLOUD_PROJECT")
	serviceAccountKeyPath = os.Getenv("GCLOUD_SERVICE_ACCOUNT_PATH")
)

func main() {
	fmt.Println(len(os.Args))
	if len(os.Args) == 1 {
		log.Fatalf("please specify an ID to query for")
	}
	qId := os.Args[1]
	var qLastName = "Warren"
	if len(os.Args) == 3 {
		qLastName = os.Args[2]
	}
	ctx := context.Background()
	client, err := firestore.NewClient(
		ctx,
		projectId,
		option.WithCredentialsFile(serviceAccountKeyPath),
	)
	if err != nil {
		log.Fatalf("error creating firestore client: %v", err)
	}
	people := NewPeople(client)

	// by id
	log.Println("-----")
	log.Printf("Fetching person with id: %v", qId)
	person := people.GetPersonById(qId)
	p1, err := NewPerson(ctx, person)
	if err != nil {
		log.Printf("error getting users: %v", err)
	} else {
		p1.Print()
	}

	// query
	log.Println("-----")
	log.Printf("Querying for lastname == %v", qLastName)
	refs, err := people.GetPeopleByLastName(ctx, qLastName)
	log.Printf("total query size: %v", len(refs))
	if err != nil {
		log.Printf("error querying people: %v", err)
	}
	for _, r := range refs {
		rp, err := NewPerson(ctx, r)
		if err != nil {
			log.Printf("error getting users: %v", err)
		}
		rp.Print()
	}
}
