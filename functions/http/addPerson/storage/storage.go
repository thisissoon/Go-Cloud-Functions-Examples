package storage

import (
	"context"

	"cloud.google.com/go/firestore"
)

type Storage interface {
	SaveDoc(ctx context.Context, collection string, data interface{}) (*firestore.DocumentRef, error)
}

type Store struct {
	Client *firestore.Client
}

func (s Store) SaveDoc(ctx context.Context, collection string, data interface{}) (*firestore.DocumentRef, error) {
	ref, _, err := s.Client.Collection(collection).Add(ctx, data)
	if err != nil {
		return nil, err
	}
	return ref, nil
}
