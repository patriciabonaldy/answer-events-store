package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/business"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"net/http"
)

type AnswerHandler struct {
	service business.Service
	log     logger.Logger
}

func New(service business.Service, log logger.Logger) AnswerHandler {
	return AnswerHandler{
		service: service,
		log:     log,
	}
}

type createRequest struct {
	Data map[string]string `json:"data" binding:"required"`
}

type createResponse struct {
	ID      string            `json:"answer_id" binding:"required"`
	EventID string            `json:"event_id" binding:"required"`
	Event   string            `json:"event" binding:"required"`
	Data    map[string]string `json:"data" binding:"required"`
}

// Create returns an HTTP handler for answer creation.
func (a *AnswerHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		ans, err := a.service.CreateAnswer(ctx, req.Data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		resp, err := toCreateResponse(ans)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, resp)
	}
}

func toCreateResponse(answ *internal.Answer) (createResponse, error) {
	ev := answ.Events[len(answ.Events)-1]
	var data map[string]string

	err := json.Unmarshal(ev.RawData, &data)
	if err != nil {
		return createResponse{}, err
	}

	return createResponse{
		ID:      answ.ID,
		EventID: ev.EventID,
		Event:   string(ev.Type),
		Data:    data,
	}, nil
}

type request struct {
	ID    string            `json:"id" binding:"required"`
	Event string            `json:"event" binding:"required"`
	Data  map[string]string `json:"data" binding:"required"`
}
