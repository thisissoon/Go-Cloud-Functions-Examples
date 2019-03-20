package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type People struct {
	ColRef *firestore.CollectionRef
}

func (p *People) GetPersonById(id string) *firestore.DocumentRef {
	return p.ColRef.Doc(id)
}

func (p *People) GetPeopleByLastName(ctx context.Context, name string) ([]*firestore.DocumentRef, error) {
	q := p.ColRef.Where("lastname", "==", name)
	iter := q.Documents(ctx)
	defer iter.Stop()
	var refs []*firestore.DocumentRef
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating: %v", err)
		}
		if len(refs) == 0 {
			refs = []*firestore.DocumentRef{
				doc.Ref,
			}
		} else {
			refs = append(refs, doc.Ref)
		}
	}
	return refs, nil
}

func (p *People) CreatePerson(ctx context.Context, person *Person) (*firestore.DocumentRef, error) {
	ref, _, err := p.ColRef.Add(ctx, person)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

func (p *People) DeletePerson(ctx context.Context, id string) error {
	_, err := p.ColRef.Doc(id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewPeople(client *firestore.Client) *People {
	return &People{
		ColRef: client.Collection("people"),
	}
}
