package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/business"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
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

// Create returns an HTTP handler for answer creation.
func (a *AnswerHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var ans *internal.Answer
		var err error
		if len(req.ID) == 0 {
			ans, err = a.service.CreateAnswer(ctx, req.Data)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			resp, err := toResponse(ans)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			ctx.JSON(http.StatusCreated, resp)
		}

		err = a.service.UpdateAnswer(ctx, req.ID, "create", req.Data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

	}
}

func (a *AnswerHandler) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requestID
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		ans, err := a.service.GetAnswerByID(ctx, req.ID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				ctx.Status(http.StatusBadRequest)
				return
			}

			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		resp, err := toResponse(ans)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func (a *AnswerHandler) GetHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requestID
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		ans, err := a.service.GetAnswerByID(ctx, req.ID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				ctx.Status(http.StatusBadRequest)
				return
			}

			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		resp, err := toHistoryResponse(ans)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// Update returns an HTTP handler for answer creation.
func (a *AnswerHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqID requestID
		if err := ctx.ShouldBindUri(&reqID); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		var req request
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if req.Event != "update" {
			ctx.JSON(http.StatusBadRequest, "invalid event state")
			return
		}

		err := a.service.UpdateAnswer(ctx, reqID.ID, req.Event, req.Data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Status(http.StatusOK)
	}
}

// Delete returns an HTTP handler for answer creation.
func (a *AnswerHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.Status(http.StatusBadRequest)
			return
		}

		err := a.service.UpdateAnswer(ctx, id, "delete", nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func toResponse(answ *internal.Answer) (response, error) {
	ev := answ.Events[len(answ.Events)-1]
	var data map[string]string

	err := json.Unmarshal(ev.RawData, &data)
	if err != nil {
		return response{}, err
	}

	return response{
		ID:       answ.ID,
		Event:    string(ev.Type),
		Data:     data,
		CreateAt: ev.Timestamp,
	}, nil
}

func toHistoryResponse(answ *internal.Answer) (historyResponse, error) {
	resp := historyResponse{
		ID:     answ.ID,
		Events: []event{},
	}

	for _, ev := range answ.Events {
		var data map[string]string

		err := json.Unmarshal(ev.RawData, &data)
		if err != nil {
			return historyResponse{}, err
		}

		resp.Events = append(resp.Events, event{
			EventID:   ev.EventID,
			Type:      string(ev.Type),
			Data:      data,
			Timestamp: ev.Timestamp,
			Version:   ev.Version,
		})
	}

	return resp, nil
}
