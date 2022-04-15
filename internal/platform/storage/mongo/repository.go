package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/config"
)

const eventCollectionName = "event"

var (
	ErrIDIsEmpty          = errors.New("invalid ID")
	ErrCollectionNotFound = errors.New("collection not found")
)

// Repository is a mongo EventRepository implementation.
type Repository struct {
	databaseName string
	db           *mongo.Client
}

// NewDBStorage initializes a mongo-based implementation of Storage.
func NewDBStorage(ctx context.Context, cfg *config.MongoConfig) (*Repository, error) {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(cfg.URI).
			SetAuth(options.Credential{Username: cfg.User,
				Password: cfg.Password}))
	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Repository{
		databaseName: cfg.DatabaseName,
		db:           client,
	}, nil
}

func (r *Repository) GetByID(ctx context.Context, answerID string) (*internal.Answer, error) {
	objectID, err := primitive.ObjectIDFromHex(answerID)
	if err != nil {
		return nil, err
	}

	var result internal.Answer
	err = r.getCollection(eventCollectionName).FindOne(ctx, bson.M{
		"_id": objectID,
	}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &internal.Answer{}, nil
		}

		return nil, err
	}

	return &result, nil
}

func (r *Repository) Save(ctx context.Context, answer internal.Answer) (internal.Answer, error) {
	result, err := r.getCollection(eventCollectionName).InsertOne(ctx, answer)
	if err != nil {
		return internal.Answer{}, err
	}

	answer.AnswerID = result.InsertedID.(primitive.ObjectID)

	return answer, nil
}

func (r *Repository) Update(ctx context.Context, answer internal.Answer) error {
	if answer.AnswerID.IsZero() {
		return ErrIDIsEmpty
	}

	opts := options.Replace().SetUpsert(true)
	filter := bson.D{{
		Key:   "_id",
		Value: answer.AnswerID,
	}}
	result, err := r.getCollection(eventCollectionName).
		ReplaceOne(ctx, filter, answer, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}

	return nil
}

func (r *Repository) getCollection(collectionName string) *mongo.Collection {
	return r.db.Database(r.databaseName).Collection(collectionName, nil)
}
