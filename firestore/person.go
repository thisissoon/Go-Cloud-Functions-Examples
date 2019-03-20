package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/genproto/googleapis/type/latlng"
)

type Person struct {
	Id        string         `firestore:"id"`
	Firstname string         `firestore:"firstname"`
	Lastname  string         `firestore:"lastname"`
	Dob       time.Time      `firestore:"dob"`
	CreatedAt time.Time      `firestore:"createdAt"`
	Postcode  string         `firestore:"postcode"`
	Location  *latlng.LatLng `firestore:"location"`
}

func (p *Person) Print() {
	fmt.Printf(
		"Person added at %v:\n  ID: %s\n  Name: %s %s\n  DOB: %v\n  Postcode: %v\n  Location: %v\n",
		p.CreatedAt,
		p.Id,
		p.Firstname,
		p.Lastname,
		p.Dob,
		p.Postcode,
		p.Location,
	)
}

func NewPerson(ctx context.Context, docRef *firestore.DocumentRef) (Person, error) {
	docsnap, err := docRef.Get(ctx)
	if err != nil {
		return Person{}, fmt.Errorf("error reading data: %v", err)
	}
	if !docsnap.Exists() {
		return Person{}, fmt.Errorf("person doesn't exist")
	}
	var person Person
	if err := docsnap.DataTo(&person); err != nil {
		return Person{}, fmt.Errorf("error marshalling data: %v", err)
	}
	person.Id = docRef.ID
	return person, nil
}
