package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type Storage interface {
	UpdateDoc(ctx context.Context, collection string, doc string, update firestore.Update) error
}

type Store struct {
	Client *firestore.Client
}

func (s Store) UpdateDoc(ctx context.Context, collection string, doc string, update firestore.Update) error {
	_, err := s.Client.
		Collection(collection).
		Doc(doc).
		Update(ctx, []firestore.Update{update})
	if err != nil {
		return fmt.Errorf("Update error: %v", err)
	}
	return nil
}
