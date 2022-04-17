package mongo

import (
	"context"

	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/pkg/errors"

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
	var result AnswerDB

	err := r.getCollection(eventCollectionName).
		FindOne(ctx, bson.M{"answer_id": answerID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return internal.Answer{}, internal.ErrAnswerNotFound
		}

		return internal.Answer{}, err
	}

	return parseToBusinessAnswer(result), nil
}

func (r *Repository) Save(ctx context.Context, answer internal.Answer) error {
	answerDB := parseToAnswerDB(answer)
	_, err := r.getCollection(eventCollectionName).InsertOne(ctx, answerDB)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, answer internal.Answer) error {
	if len(answer.ID) == 0 {
		return internal.ErrIDIsEmpty
	}

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"answer_id": answer.ID}
	answerDB := parseToAnswerDB(answer)
	result, err := r.getCollection(eventCollectionName).
		ReplaceOne(ctx, filter, answerDB, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount != 0 {
		r.log.Info("matched and replaced an existing document")
	}
	if result.UpsertedCount != 0 {
		r.log.Info("inserted a new document with ID %v\n", result.UpsertedID)
	}

	return nil
}

func (r *Repository) getCollection(collectionName string) *mongo.Collection {
	return r.db.Database(r.databaseName).Collection(collectionName, nil)
}
