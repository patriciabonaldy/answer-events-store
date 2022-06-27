package update

import (
	"context"
	"testing"
	"time"

	"github.com/patriciabonaldy/bequest_challenge/internal/platform/pubsub"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/pubsub/pubsubMock"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/storage/storagemocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func Test_service_UpdateAnswer(t *testing.T) {
	tests := []struct {
		name      string
		eventID   string
		eventType string
		data      map[string]string
		producer  func() pubsub.Producer
		repo      func() internal.Storage
		wantErr   bool
	}{
		{
			name:      "error getting last event",
			eventID:   "b47915e6-bd66-11ec-aaa8-acde48001122",
			eventType: "update",
			data: map[string]string{
				"key": "9999",
			},
			producer: func() pubsub.Producer {
				productMock := new(pubsubMock.Producer)
				productMock.On("Produce", mock.Anything, mock.Anything).
					Return(nil)

				return productMock
			},
			repo: func() internal.Storage {
				repoMock := new(storagemocks.Storage)
				repoMock.On("GetByID", mock.Anything, mock.Anything).
					Return(internal.Answer{}, errors.New("something unexpected happened"))

				return repoMock

			},
			wantErr: true,
		},
		{
			name:      "error invalid event state",
			eventID:   "b47915e6-bd66-11ec-aaa8-acde48001122",
			eventType: "create",
			data: map[string]string{
				"key": "9999",
			},
			producer: func() pubsub.Producer {
				productMock := new(pubsubMock.Producer)
				productMock.On("Produce", mock.Anything, mock.Anything).
					Return(nil)

				return productMock
			},
			repo: func() internal.Storage {
				mockA := mockAnswer()
				mockA.AddEvent(mockEvent("update"))

				repoMock := new(storagemocks.Storage)
				repoMock.On("GetByID", mock.Anything, mock.Anything).
					Return(mockA, nil)

				return repoMock

			},
			wantErr: true,
		},
		{
			name:      "error invalid event state",
			eventID:   "b47915e6-bd66-11ec-aaa8-acde48001122",
			eventType: "create",
			data: map[string]string{
				"key": "9999",
			},
			producer: func() pubsub.Producer {
				productMock := new(pubsubMock.Producer)
				productMock.On("Produce", mock.Anything, mock.Anything).
					Return(nil)

				return productMock
			},
			repo: func() internal.Storage {
				mockA := mockAnswer()
				mockA.AddEvent(mockEvent("update"))

				repoMock := new(storagemocks.Storage)
				repoMock.On("GetByID", mock.Anything, mock.Anything).
					Return(mockA, nil)

				return repoMock

			},
			wantErr: true,
		},
		{
			name:      "error update answers",
			eventID:   "b47915e6-bd66-11ec-aaa8-acde48001122",
			eventType: "delete",
			data: map[string]string{
				"key": "9999",
			},
			producer: func() pubsub.Producer {
				productMock := new(pubsubMock.Producer)
				productMock.On("Produce", mock.Anything, mock.Anything).
					Return(errors.New("something unexpected happened"))

				return productMock
			},
			repo: func() internal.Storage {
				mockA := mockAnswer()
				mockA.AddEvent(mockEvent("update"))

				repoMock := new(storagemocks.Storage)
				repoMock.On("GetByID", mock.Anything, mock.Anything).
					Return(mockA, nil)

				return repoMock

			},
			wantErr: true,
		},
		{
			name:      "success",
			eventID:   "b47915e6-bd66-11ec-aaa8-acde48001122",
			eventType: "delete",
			data: map[string]string{
				"key": "9999",
			},
			producer: func() pubsub.Producer {
				productMock := new(pubsubMock.Producer)
				productMock.On("Produce", mock.Anything, mock.Anything).
					Return(nil)

				return productMock
			},
			repo: func() internal.Storage {
				mockA := mockAnswer()
				mockA.AddEvent(mockEvent("update"))

				repoMock := new(storagemocks.Storage)
				repoMock.On("GetByID", mock.Anything, mock.Anything).
					Return(mockA, nil)

				return repoMock

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.producer(), tt.repo(), logger.New())
			if err := s.UpdateAnswer(context.Background(), tt.eventID, tt.eventType, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAnswer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_checkEvent(t *testing.T) {
	tests := []struct {
		name    string
		answer  internal.Answer
		events  []internal.Event
		wantErr bool
	}{
		{
			name:   "invalid events states",
			answer: mockAnswer(),
			events: []internal.Event{
				mockEvent("delete"),
				mockEvent("update"),
			},
			wantErr: true,
		},
		{
			name:   "valid events states",
			answer: mockAnswer(),
			events: []internal.Event{
				mockEvent("delete"),
				mockEvent("create"),
				mockEvent("update"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{}
			var err error
			for _, event := range tt.events {
				err = s.checkEvent(tt.answer, string(event.Type))
				if err != nil {
					break
				}

				tt.answer.AddEvent(event)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("checkEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func mockAnswer() internal.Answer {
	return internal.Answer{
		ID: "b47915e6-bd66-11ec-aaa8-acde48001122",
		Events: []internal.Event{
			mockEvent("create"),
		},
		CreateAt: time.Time{},
		UpdateAt: time.Time{},
	}
}

func mockEvent(eventType string) internal.Event {
	return internal.Event{
		EventID: "b477344c-bd66-11ec-aaa8-acde48001122",
		Type:    internal.EventType(eventType),
		RawData: []byte("{}"),
	}
}
