package business

import (
	"context"
	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/storage/storagemocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
	"time"
)

func Test_service_CreateAnswer(t *testing.T) {
	tests := []struct {
		name    string
		repo    func() internal.Storage
		data    map[string]string
		wantErr bool
	}{
		{
			name: "error creating a new event",
			repo: func() internal.Storage {
				repoMock := new(storagemocks.Storage)
				repoMock.On("Save", mock.Anything, mock.Anything).
					Return(errors.New("something unexpected happened"))

				return repoMock

			},
			wantErr: true,
		},
		{
			name: "success",
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
			s := NewService(repo, logger.New())
			if _, err := s.CreateAnswer(context.Background(), tt.data); (err != nil) != tt.wantErr {
				t.Errorf("CreateAnswer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_GetAnswers(t *testing.T) {
	tests := []struct {
		name    string
		repo    func() internal.Storage
		eventID string
		want    func() *internal.Answer
		wantErr bool
	}{
		{
			name: "error getting answer",
			repo: func() internal.Storage {
				repoMock := new(storagemocks.Storage)
				repoMock.On("GetByID", mock.Anything, mock.Anything).
					Return(internal.Answer{}, errors.New("something unexpected happened"))

				return repoMock

			},
			eventID: "b47915e6-bd66-11ec-aaa8-acde48001122",
			want: func() *internal.Answer {
				return nil
			},
			wantErr: true,
		},
		{
			name: "success",
			repo: func() internal.Storage {
				mockA := mockAnswer()
				mockA.AddEvent(mockEvent("update"))
				repoMock := new(storagemocks.Storage)
				repoMock.On("GetByID", mock.Anything, mock.Anything).
					Return(mockA, nil)

				return repoMock

			},
			eventID: "b47915e6-bd66-11ec-aaa8-acde48001122",
			want: func() *internal.Answer {
				mockA := mockAnswer()
				mockA.AddEvent(mockEvent("update"))

				return &mockA
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.repo(), logger.New())
			got, err := s.GetAnswerByID(context.Background(), tt.eventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnswerByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want := tt.want()
			if !reflect.DeepEqual(got, want) {
				t.Errorf("GetAnswerByID() got = %v, want %v", got, want)
			}
		})
	}
}

func Test_service_UpdateAnswer(t *testing.T) {
	tests := []struct {
		name      string
		eventID   string
		eventType string
		data      map[string]string
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
			repo: func() internal.Storage {
				mockA := mockAnswer()
				mockA.AddEvent(mockEvent("update"))

				repoMock := new(storagemocks.Storage)
				repoMock.On("GetByID", mock.Anything, mock.Anything).
					Return(mockA, nil)

				repoMock.On("Update", mock.Anything, mock.Anything).
					Return(internal.Answer{}, errors.New("something unexpected happened"))

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
			repo: func() internal.Storage {
				mockA := mockAnswer()
				mockA.AddEvent(mockEvent("update"))

				repoMock := new(storagemocks.Storage)
				repoMock.On("GetByID", mock.Anything, mock.Anything).
					Return(mockA, nil)

				repoMock.On("Update", mock.Anything, mock.Anything).
					Return(mockA, nil)

				return repoMock

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.repo(), logger.New())
			if _, err := s.UpdateAnswer(context.Background(), tt.eventID, tt.eventType, tt.data); (err != nil) != tt.wantErr {
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
