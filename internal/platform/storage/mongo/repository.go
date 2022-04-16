package mongo

import (
	"context"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/config"
)

const eventCollectionName = "event"

// Repository is a mongo EventRepository implementation.
type Repository struct {
	databaseName string
	db           *mongo.Client
	log          logger.Logger
}

// NewDBStorage initializes a mongo-based implementation of Storage.
func NewDBStorage(ctx context.Context, cfg *config.Database, log logger.Logger) (*Repository, error) {
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
		log:          log,
	}, nil
}

func (r *Repository) GetByID(ctx context.Context, answerID string) (internal.Answer, error) {
	objectID, err := primitive.ObjectIDFromHex(answerID)
	if err != nil {
		return internal.Answer{}, err
	}

	var result AnswerDB
	err = r.getCollection(eventCollectionName).FindOne(ctx, bson.M{
		"_id": objectID,
	}).Decode(&result)
	if err != nil {
		return internal.Answer{}, err
	}

	return parseToBusinessAnswer(result), nil
}

func (r *Repository) Save(ctx context.Context, answer internal.Answer) (internal.Answer, error) {
	answerDB, err := parseToAnswerDB(answer)
	if err != nil {
		return internal.Answer{}, err
	}

	result, err := r.getCollection(eventCollectionName).InsertOne(ctx, answerDB)
	if err != nil {
		return internal.Answer{}, err
	}

	id := result.InsertedID.(primitive.ObjectID)

	return r.GetByID(ctx, id.Hex())
}

func (r *Repository) Update(ctx context.Context, answer internal.Answer) (internal.Answer, error) {
	if len(answer.ID) == 0 {
		return internal.Answer{}, internal.ErrIDIsEmpty
	}

	if _, err := r.GetByID(ctx, answer.ID); err != nil {
		return internal.Answer{}, err
	}

	opts := options.Replace().SetUpsert(true)
	objectID, err := primitive.ObjectIDFromHex(answer.ID)
	if err != nil {
		return internal.Answer{}, err
	}

	filter := bson.M{"_id": objectID}
	answerDB, err := parseToAnswerDB(answer)
	if err != nil {
		return internal.Answer{}, err
	}

	result, err := r.getCollection(eventCollectionName).
		ReplaceOne(ctx, filter, answerDB, opts)
	if err != nil {
		return internal.Answer{}, err
	}

	if result.MatchedCount != 0 {
		r.log.Info("matched and replaced an existing document")
	}
	if result.UpsertedCount != 0 {
		r.log.Info("inserted a new document with ID %v\n", result.UpsertedID)
	}

	return r.GetByID(ctx, answer.ID)
}

func (r *Repository) getCollection(collectionName string) *mongo.Collection {
	return r.db.Database(r.databaseName).Collection(collectionName, nil)
}
