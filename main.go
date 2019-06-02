package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
)

func main() {
	client := createFirestoreClient()
	defer client.Close()

	controller := createController(client)

	router := NewRouter(controller)

	fmt.Println("listening on port 8080")
	http.ListenAndServe(":8080", router)
}

func createController(client *firestore.Client) *Controller {
	repo := createRepository(client)
	return &Controller{repo}
}

func createRepository(client *firestore.Client) *Repository {
	return &Repository{client}
}

func createFirestoreClient() *firestore.Client {
	ctx := context.Background()
	serviceAccount := option.WithCredentialsFile("/Users/nikhilkulkarni/downloads/i-have-friends-7145215c635d.json")
	app, err := firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		log.Fatalln("Failed to create firebase app", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("Failed to create firestore client", err)
	}

	return client
}
