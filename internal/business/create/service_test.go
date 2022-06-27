package create

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/patriciabonaldy/bequest_challenge/internal/platform/pubsub"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/pubsub/pubsubMock"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/storage/storagemocks"
)

func Test_service_CreateAnswer(t *testing.T) {
	tests := []struct {
		name     string
		producer func() pubsub.Producer
		repo     func() internal.Storage
		data     map[string]string
		wantErr  bool
	}{
		{
			name: "error creating a new event",
			producer: func() pubsub.Producer {
				productMock := new(pubsubMock.Producer)
				productMock.On("Produce", mock.Anything, mock.Anything).
					Return(errors.New("something unexpected happened"))

				return productMock
			},
			repo: func() internal.Storage {
				repoMock := new(storagemocks.Storage)
				return repoMock

			},
			wantErr: true,
		},
		{
			name: "success",
			producer: func() pubsub.Producer {
				productMock := new(pubsubMock.Producer)
				productMock.On("Produce", mock.Anything, mock.Anything).
					Return(nil)

				return productMock
			},
			repo: func() internal.Storage {
				repoMock := new(storagemocks.Storage)
				repoMock.On("Save", mock.Anything, mock.Anything).
					Return(nil)

				return repoMock

			},
			data: map[string]string{
				"key": "9999",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.repo()
			s := NewService(tt.producer(), repo, logger.New())
			if _, err := s.CreateAnswer(context.Background(), "", tt.data); (err != nil) != tt.wantErr {
				t.Errorf("CreateAnswer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
