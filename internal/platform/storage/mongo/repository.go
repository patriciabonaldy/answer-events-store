package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/patriciabonaldy/bequest_challenge/internal"
)

const eventCollectionName = "event"

var ErrCollectionNotFound = errors.New("collection not found")

// Repository is a mongo EventRepository implementation.
type Repository struct {
	databaseName string
	db           *mongo.Client
}

type Config struct {
	URI             string
	databaseName    string
	User            string
	Password        string
	connectTimeout  time.Duration
	minPoolSize     uint64
	maxPoolSize     uint64
	maxConnIdleTime time.Duration
}

// NewDBStorage initializes a mongo-based implementation of Storage.
func NewDBStorage(ctx context.Context, cfg *Config) (*Repository, error) {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(cfg.URI).
			SetAuth(options.Credential{Username: cfg.User, Password: cfg.Password}).
			SetConnectTimeout(cfg.connectTimeout).
			SetMaxConnIdleTime(cfg.maxConnIdleTime).
			SetMinPoolSize(cfg.minPoolSize).
			SetMaxPoolSize(cfg.maxPoolSize))
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
		databaseName: cfg.databaseName,
		db:           client,
	}, nil
}

func (r *Repository) GetByID(ctx context.Context, answerID string) (*internal.Answer, error) {
	collection := r.getCollection(eventCollectionName)
	if collection == nil {
		return nil, ErrCollectionNotFound
	}

	var result internal.Answer
	err := collection.FindOne(ctx, bson.D{{"_id", answerID}}, nil).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &internal.Answer{}, nil
		}

		return nil, err
	}

	return &result, nil
}

func (r *Repository) Save(ctx context.Context, answer internal.Answer) error {
	collection := r.getCollection(eventCollectionName)
	if collection == nil {
		return ErrCollectionNotFound
	}

	_, err := collection.InsertOne(ctx, answer, &options.InsertOneOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, answer internal.Answer) error {
	collection := r.getCollection(eventCollectionName)
	if collection == nil {
		return ErrCollectionNotFound
	}

	opts := options.Replace().SetUpsert(true)
	filter := bson.D{{"_id", answer.ID}}
	result, err := collection.ReplaceOne(ctx, filter, answer, opts)
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
