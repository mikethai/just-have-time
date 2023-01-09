package firestoreClient

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

type firestoreClient struct {
	client  *firestore.Client
	context context.Context
}

func NewFirestoreClient() *firestoreClient {
	// Sets your Google Cloud Platform project ID.
	projectID := "kkchack22-just-have-time"

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	return &firestoreClient{
		client:  client,
		context: ctx,
	}
}

func (firestoreClient *firestoreClient) Get(collection string, docID string) (*firestore.DocumentSnapshot, error) {

	dsnap, err := firestoreClient.client.Collection(collection).Doc(docID).Get(firestoreClient.context)
	return dsnap, err
}

func (firestoreClient *firestoreClient) Set(collection string, docID string, anyStructure any) {
	firestoreClient.client.Collection(collection).Doc(docID).Set(firestoreClient.context, anyStructure)
}
