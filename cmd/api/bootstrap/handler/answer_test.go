package handler

import (
	"bytes"
	"encoding/json"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/storage/storagemocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/patriciabonaldy/bequest_challenge/internal/business"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Create(t *testing.T) {
	repositoryMock := new(storagemocks.Storage)
	repositoryMock.On("Save", mock.Anything, mock.Anything).
		Return(nil)
	log := logger.New()
	svc := business.NewService(repositoryMock, log)
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
