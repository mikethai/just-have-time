package firestoreClient

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

const gcpProjectID = "kkchack22-just-have-time"

type FirestoreClient interface {
	Get(collection string, docID string) (*firestore.DocumentSnapshot, error)
	Set(collection string, docID string, anyStructure any)
	CloseConnection()
}

type firestoreClient struct {
	client  *firestore.Client
	context context.Context
}

func NewFirestoreClient() *firestoreClient {
	// Sets your Google Cloud Platform project ID.
	projectID := gcpProjectID

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

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

func (firestoreClient *firestoreClient) CloseConnection() {
	firestoreClient.client.Close()
}
