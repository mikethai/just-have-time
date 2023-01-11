package firestoreClient

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/mikethai/just-have-time/internal/model"
	"google.golang.org/api/iterator"
)

const gcpProjectID = "kkchack22-just-have-time"

type FirestoreClient interface {
	Get(collection string, docID string) (*firestore.DocumentSnapshot, error)
	Set(collection string, docID string, anyStructure any)
	BatchGETSontInfo(docIDs []string) (map[string]*model.SongInfo, error)
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

func (firestoreClient *firestoreClient) BatchGETSontInfo(docIDs []string) (map[string]*model.SongInfo, error) {

	songInfos := make(map[string]*model.SongInfo)
	iter := firestoreClient.client.Collection("track").Where("id", "in", docIDs).Documents(firestoreClient.context)
	defer iter.Stop() // add this line to ensure resources cleaned up
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var songInfo *model.SongInfo
		if err := doc.DataTo(&songInfo); err != nil {
			return nil, err
		}

		songInfos[songInfo.Id] = songInfo
	}

	return songInfos, nil
}

func UniqueStings(stringSlices ...[]string) []string {
	uniqueMap := map[string]bool{}

	for _, stringSlice := range stringSlices {
		for _, number := range stringSlice {
			uniqueMap[number] = true
		}
	}

	// Create a slice with the capacity of unique items
	// This capacity make appending flow much more efficient
	result := make([]string, 0, len(uniqueMap))

	for key := range uniqueMap {
		result = append(result, key)
	}

	return result
}
