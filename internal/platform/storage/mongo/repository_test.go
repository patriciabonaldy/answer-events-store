package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/config"
)

var db *mongo.Client
var uri string

func TestMain(m *testing.M) {
	// setup
	pool, resource := mongoContainer()
	// run tests
	exitCode := m.Run()
	// kill and remove the container
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	// disconnect mongodb client
	if err := db.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

func TestNewDBStorage(t *testing.T) {
	_, err := NewDBStorage(context.Background(), &config.Config{})
	require.Error(t, err)

	_, err = NewDBStorage(context.Background(), &config.Config{
		URI:          uri,
		DatabaseName: "",
		User:         "root",
		Password:     "password",
	})
	require.NoError(t, err)
}

func TestRepository_Save(t *testing.T) {
	repo := &Repository{
		databaseName: "test",
		db:           db,
	}

	createdAt := primitive.Timestamp{
		T: uint32(time.Now().Unix()),
	}
	ev1, err := internal.NewEvent("create", []byte("{}"), 0)
	require.NoError(t, err)
	answer := internal.Answer{
		CreateAt: createdAt,
		Events:   []internal.Event{ev1},
	}

	ctx := context.Background()
	answer, err = repo.Save(ctx, answer)
	assert.NoError(t, err)

	got, err := repo.GetByID(ctx, answer.AnswerID.Hex())
	assert.NoError(t, err)

	assert.Equal(t, &answer, got)

	fmt.Println(got)
}

func TestRepository_Update(t *testing.T) {
	repo := &Repository{
		databaseName: "test",
		db:           db,
	}

	ctx := context.Background()
	cases := []struct {
		name        string
		fn          func() internal.Answer
		wantErr     bool
		expectedErr error
	}{
		{
			name: "error ID is empty",
			fn: func() internal.Answer {
				return mockAnswer(t)
			},
			wantErr:     true,
			expectedErr: ErrIDIsEmpty,
		},
		{
			name: "success",
			fn: func() internal.Answer {
				answer, err := repo.Save(ctx, mockAnswer(t))
				require.NoError(t, err)

				answer.Events = append(answer.Events, mockEvent(t))
				return answer
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			answer := tt.fn()
			err := repo.Update(ctx, answer)
			assert.Equal(t, tt.expectedErr, err)

			if !tt.wantErr {
				got, err := repo.GetByID(ctx, answer.AnswerID.Hex())
				require.NoError(t, err)

				assert.Equal(t, &answer, got)
			}
		})
	}
}

func mockAnswer(t *testing.T) internal.Answer {
	createdAt := primitive.Timestamp{
		T: uint32(time.Now().Unix()),
	}

	answer := internal.Answer{
		CreateAt: createdAt,
		Events: []internal.Event{
			mockEvent(t),
		},
	}
	return answer
}

func mockEvent(t *testing.T) internal.Event {
	ev, err := internal.NewEvent("create", []byte("{}"), 0)
	require.NoError(t, err)

	return ev
}

func mongoContainer() (*dockertest.Pool, *dockertest.Resource) {
	const MONGO_INITDB_ROOT_USERNAME = "root"
	const MONGO_INITDB_ROOT_PASSWORD = "password"

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	environmentVariables := []string{
		"MONGO_INITDB_ROOT_USERNAME=" + MONGO_INITDB_ROOT_USERNAME,
		"MONGO_INITDB_ROOT_PASSWORD=" + MONGO_INITDB_ROOT_PASSWORD,
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env:        environmentVariables,
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	err = pool.Retry(func() error {
		var err error
		uri = fmt.Sprintf("mongodb://%s:%s@localhost:%s",
			MONGO_INITDB_ROOT_USERNAME, MONGO_INITDB_ROOT_PASSWORD,
			resource.GetPort("27017/tcp"))
		db, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				uri,
			),
		)
		if err != nil {
			return err
		}
		return db.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}
	return pool, resource
}