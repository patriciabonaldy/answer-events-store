package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/storage/storagemocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
	"github.com/patriciabonaldy/bequest_challenge/internal/business"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/stretchr/testify/assert"
)

var timeN = time.Now()

func TestHandler_Get(t *testing.T) {
	repositoryMock := new(storagemocks.Storage)
	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(internal.Answer{}, errors.New("something unexpected happened")).Once()
	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(mockAnswer(), nil).Once()
	log := logger.New()
	svc := business.NewService(nil, repositoryMock, log)
	handler := New(svc, log)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/answers/:id", handler.GetAnswer())

	t.Run("given a invalid request it returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/answers/0", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a error it returns 500", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/answers/b47915e6-bd66-11ec-aaa8-acde48001122", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("given a valid request it returns 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/answers/b47915e6-bd66-11ec-aaa8-acde48001122", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		var resp Response
		err = json.NewDecoder(res.Body).Decode(&resp)
		require.NoError(t, err)

		want := Response{
			ID:    "b47915e6-bd66-11ec-aaa8-acde48001122",
			Event: "create",
			Data: map[string]string{
				"en": "Map",
			},
			CreateAt: timeN,
		}
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, resp.ID, want.ID)
		assert.Equal(t, resp.Event, want.Event)
		assert.Equal(t, reflect.DeepEqual(resp.Data, want.Data), true)
	})
}

func TestHandler_GetHistory(t *testing.T) {
	repositoryMock := new(storagemocks.Storage)
	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(internal.Answer{}, errors.New("something unexpected happened")).Once()
	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(mockAnswer(), nil).Once()
	log := logger.New()
	svc := business.NewService(nil, repositoryMock, log)
	handler := New(svc, log)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/answers/:id/history", handler.GetHistory())

	t.Run("given a invalid request it returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/answers/0/history", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a error it returns 500", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/answers/b47915e6-bd66-11ec-aaa8-acde48001122/history", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("given a valid request it returns 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/answers/b47915e6-bd66-11ec-aaa8-acde48001122/history", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		var resp historyResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		require.NoError(t, err)

		want := historyResponse{
			ID: "b47915e6-bd66-11ec-aaa8-acde48001122",
			Events: []event{
				{
					EventID: "b477344c-bd66-11ec-aaa8-acde48001122",
					Type:    "create",
					Data: map[string]string{
						"en": "Map",
					},
					Timestamp: timeN,
					Version:   0,
				},
			},
		}
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, resp.ID, want.ID)

		got, err := json.Marshal(resp.Events)
		require.NoError(t, err)

		_want, err := json.Marshal(want.Events)
		require.NoError(t, err)

		assert.Equal(t, string(got), string(_want))
	})
}

func TestHandler_Create(t *testing.T) {
	repositoryMock := new(storagemocks.Storage)
	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(internal.Answer{}, internal.ErrAnswerNotFound).Once()

	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(mockAnswer(), nil).Once()

	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(internal.Answer{}, errors.New("something unexpected happened")).Once()
	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(mockAnswer(), nil).Once()
	repositoryMock.On("Save", mock.Anything, mock.Anything).
		Return(nil)
	log := logger.New()
	svc := business.NewService(nil, repositoryMock, log)
	handler := New(svc, log)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/answer", handler.Create())

	t.Run("given an invalid request it returns 400", func(t *testing.T) {
		createReq := CreateRequest{}

		b, err := json.Marshal(createReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/answer", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given an invalid id it returns 400", func(t *testing.T) {
		createReq := CreateRequest{
			ID: "string",
			Data: map[string]string{
				"key": "value",
			},
		}

		b, err := json.Marshal(createReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/answer", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a invalid previous state (delete) it returns 400", func(t *testing.T) {
		createCustomerReq := CreateRequest{
			ID: "b47915e6-bd66-11ec-aaa8-acde48001122",
			Data: map[string]string{
				"key": "value",
			},
		}

		b, err := json.Marshal(createCustomerReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/answer", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a valid request it returns 201", func(t *testing.T) {
		createCustomerReq := CreateRequest{
			Data: map[string]string{
				"key": "value",
			},
		}

		b, err := json.Marshal(createCustomerReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/answer", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}

func TestHandler_Update(t *testing.T) {
	repositoryMock := new(storagemocks.Storage)
	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(internal.Answer{}, errors.New("something unexpected happened")).Once()
	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(mockAnswer(), nil).Once()

	repositoryMock.On("Update", mock.Anything, mock.Anything).
		Return(nil)
	log := logger.New()
	svc := business.NewService(nil, repositoryMock, log)
	handler := New(svc, log)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.PUT("/answers/:id", handler.Update())

	t.Run("given a invalid request it returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, "/answers/0", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a empty body it returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, "/answers/b47915e6-bd66-11ec-aaa8-acde48001122", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a invalid event type it returns 400", func(t *testing.T) {
		createReq := Request{
			Event: "delete",
			Data: map[string]string{
				"key": "value",
			},
		}

		b, err := json.Marshal(createReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "/answers/b47915e6-bd66-11ec-aaa8-acde48001122", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a error it returns 500", func(t *testing.T) {
		createReq := Request{
			Event: "update",
			Data: map[string]string{
				"key": "value",
			},
		}

		b, err := json.Marshal(createReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "/answers/b47915e6-bd66-11ec-aaa8-acde48001122", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("given a valid request it returns 200", func(t *testing.T) {
		createReq := Request{
			Event: "update",
			Data: map[string]string{
				"key": "value",
			},
		}

		b, err := json.Marshal(createReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "/answers/b47915e6-bd66-11ec-aaa8-acde48001122", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestHandler_Delete(t *testing.T) {
	repositoryMock := new(storagemocks.Storage)
	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(internal.Answer{}, errors.New("something unexpected happened")).Once()
	repositoryMock.On("GetByID", mock.Anything, mock.Anything).
		Return(mockAnswer(), nil).Once()

	repositoryMock.On("Update", mock.Anything, mock.Anything).
		Return(nil)
	log := logger.New()
	svc := business.NewService(nil, repositoryMock, log)
	handler := New(svc, log)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.DELETE("/answers/:id", handler.Delete())

	t.Run("given a invalid request it returns 400", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/answers/0", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a error it returns 500", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/answers/b47915e6-bd66-11ec-aaa8-acde48001122", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("given a valid request it returns 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/answers/b47915e6-bd66-11ec-aaa8-acde48001122", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func mockAnswer() internal.Answer {
	return internal.Answer{
		ID: "b47915e6-bd66-11ec-aaa8-acde48001122",
		Events: []internal.Event{
			mockEvent("create"),
		},
		CreateAt: timeN,
		UpdateAt: timeN,
	}
}

func mockEvent(eventType string) internal.Event {
	return internal.Event{
		EventID:   "b477344c-bd66-11ec-aaa8-acde48001122",
		Type:      internal.EventType(eventType),
		RawData:   []byte(`{"en": "Map"}`),
		Timestamp: timeN,
	}
}
