package mongo

import (
	"context"
	"fmt"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	_, err := NewDBStorage(context.Background(), &config.Database{}, logger.New())
	require.Error(t, err)

	_, err = NewDBStorage(context.Background(), &config.Database{
		URI:          uri,
		DatabaseName: "",
		User:         "root",
		Password:     "password",
	},
		logger.New())
	require.NoError(t, err)
}

func TestRepository_Save(t *testing.T) {
	repo := &Repository{
		databaseName: "test",
		db:           db,
	}

	answer := mockAnswer()
	ctx := context.Background()
	got, err := repo.Save(ctx, answer)
	assert.NoError(t, err)

	want, err := repo.GetByID(ctx, got.ID)
	assert.NoError(t, err)
	assert.Equal(t, reflect.DeepEqual(want, got), true)
}

func TestRepository_Update(t *testing.T) {
	repo := &Repository{
		databaseName: "test",
		db:           db,
		log:          logger.New(),
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
				return mockAnswer()
			},
			wantErr:     true,
			expectedErr: internal.ErrIDIsEmpty,
		},
		{
			name: "success",
			fn: func() internal.Answer {
				answer, err := repo.Save(ctx, mockAnswer())
				require.NoError(t, err)
				answer.Events = append(answer.Events, mockEvent("update"))
				answer.Events = append(answer.Events, mockEvent("delete"))
				return answer
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			answer := tt.fn()
			answer, err := repo.Update(ctx, answer)
			assert.Equal(t, tt.expectedErr, err)

			if !tt.wantErr {
				got, err := repo.GetByID(ctx, answer.ID)
				require.NoError(t, err)

				assert.Equal(t, answer, got)
			}
		})
	}
}

func mockAnswer() internal.Answer {
	ev1 := internal.NewEvent("", internal.Create, []byte("{}"))
	answer := internal.NewAnswer(ev1)

	return answer
}

func mockEvent(event string) internal.Event {
	return internal.NewEvent("", internal.EventType(event), []byte("{}"))
}

func mockEventDB(t *testing.T) EventDB {
	id, _ := uuid.NewUUID()
	evn, err := NewEvent(id.String(), "create", []byte("{}"), 0)
	require.NoError(t, err)

	return evn
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
