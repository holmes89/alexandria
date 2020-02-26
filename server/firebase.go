package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"firebase.google.com/go"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

const (
	documentCollection = "documents"
	userCollection = "users"
)

func NewFirebaseApp() *firebase.App {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		logrus.WithError(err).Fatal("unable to create firebase applicaiton")
	}
	return app
}

type DocumentsFirestoreDatabase struct {
	client *firestore.Client
}

func NewDocumentsFirestoreDatabase(app *firebase.App) *DocumentsFirestoreDatabase {
	client, err := app.Firestore(context.Background())
	if err != nil {
		logrus.WithError(err).Fatal("unable to create firestore client")
	}
	return &DocumentsFirestoreDatabase{
		client: client,
	}
}

func (r *DocumentsFirestoreDatabase) FindAll(ctx context.Context, filter map[string]interface{}) (docs []*Document, err error) {
	docIter := r.client.Collection(documentCollection).Documents(ctx)
	for {
		d, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logrus.WithError(err).Error("unable to read results")
			return nil, errors.New("unable to fetch documents")
		}
		entity := &Document{}
		if err := d.DataTo(entity); err != nil {
			logrus.WithError(err).Error("unable to convert entity")
			return nil, errors.New("unable to convert entity")
		}
		docs = append(docs, entity)
	}
	return docs, nil
}

func (r *DocumentsFirestoreDatabase) FindByID(ctx context.Context, id string) (*Document, error) {
	d, err := r.client.Collection(documentCollection).Doc(id).Get(ctx)
	if err != nil {
		logrus.WithError(err).Error("unable to fetch document")
		return nil, errors.New("unable to fetch document")
	}
	entity := &Document{}
	if err := d.DataTo(entity); err != nil {
		logrus.WithError(err).Error("unable to convert entity")
		return nil, errors.New("unable to convert entity")
	}
	return entity, nil
}

func (r *DocumentsFirestoreDatabase) Insert(ctx context.Context, doc *Document) error {
	_, err := r.client.Collection(documentCollection).Doc(doc.ID).Set(ctx, doc)
	if err != nil {
		logrus.WithError(err).Error("unable to insert document")
		return errors.New("failed ot insert document")
	}
	return nil
}

func (r *DocumentsFirestoreDatabase) UpsertStream(ctx context.Context, input <-chan *Document) error {
	count := 0
	for doc := range input {
		//Because I'm lazy I'm going to just add
		if err := r.Insert(context.Background(), doc); err != nil {
			logrus.WithError(err).Info("unable to upsert document")
			return errors.New("unable to upsert document")
		}
		count++

	}
	logrus.WithField("count", count).Info("documents added")
	return nil
}

func (r *DocumentsFirestoreDatabase) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection(documentCollection).Doc(id).Delete(ctx)
	if err != nil {
		logrus.WithError(err).Error("unable to delete document")
		return errors.New("failed ot delete document")
	}
	return nil
}

type UserFirestoreDatabase struct {
	client *firestore.Client
}

func NewUserFirestoreDatabase(app *firebase.App) *UserFirestoreDatabase {
	client, err := app.Firestore(context.Background())
	if err != nil {
		logrus.WithError(err).Fatal("unable to create firestore client")
	}
	return &UserFirestoreDatabase{
		client: client,
	}
}

func (r *UserFirestoreDatabase) FindUserByUsername(ctx context.Context, username string) (*User, error){
	d, err := r.client.Collection(userCollection).Where("username", "==", username).Documents(ctx).Next()
	if err != nil {
		logrus.WithError(err).Error("failed to find user")
		return nil, errors.New("unable of find user")
	}

	entity := &User{}
	if err := d.DataTo(entity); err != nil {
		logrus.WithError(err).Error("unable to convert entity")
		return nil, errors.New("unable to convert entity")
	}
	return entity, nil
}