package find

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/storage/storagemocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

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
