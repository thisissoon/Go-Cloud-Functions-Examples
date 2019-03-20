package addPerson

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

type PersonDB struct {
	Firstname string    `firestore:"firstname"`
	Lastname  string    `firestore:"lastname"`
	Dob       time.Time `firestore:"dob"`
	Postcode  string    `firestore:"postcode"`
	CreatedAt time.Time `firestore:"createdAt"`
}

func (p *PersonDB) SavePerson(ctx context.Context, client *firestore.Client) (*firestore.DocumentRef, error) {
	ref, _, err := client.Collection("people").Add(ctx, p)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

type Person struct {
	Firstname string `json:"fname"`
	Lastname  string `json:"lname"`
	Dob       string `json:"dob"`
	Postcode  string `json:"postcode"`
}

func (p *Person) DobTimestamp() (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, p.Dob)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (p *Person) ToDBPerson() (*PersonDB, error) {
	ts, err := p.DobTimestamp()
	if err != nil {
		return nil, err
	}
	return &PersonDB{
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Dob:       ts,
		Postcode:  p.Postcode,
		CreatedAt: time.Now(),
	}, nil
}
