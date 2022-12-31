package firestoreClient

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func Get(collection string, docID string) (*firestore.DocumentSnapshot, error) {
	ctx := context.Background()
	fsc := createClient(ctx)
	defer fsc.Close()

	dsnap, err := fsc.Collection(collection).Doc(docID).Get(ctx)
	return dsnap, err
}

func Set(collection string, docID string, anyStructure any) {
	ctx := context.Background()
	fsc := createClient(ctx)
	defer fsc.Close()

	fsc.Collection(collection).Doc(docID).Set(ctx, anyStructure)
}

func createClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "kkchack22-just-have-time"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}
