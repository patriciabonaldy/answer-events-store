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

// Create godoc
// @Summary      Create an event
// @Description  if ID params is not empty will update the record other cases will create a new answer
// @Tags         answer
// @Accept       json
// @Produce      plain
// @Param        message  body  CreateRequest  true  "Request"
// @Success      200  {string}  string         "success"
// @Success      201  {string}  string         "success"
// @Failure      400  {string}  string         "bad Request"
// @Failure      500  {string}  string         "fail"
// @Router       /answer [post]
func (a *AnswerHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var ans *internal.Answer
		var err error
		if len(req.ID) == 0 {
			ans, err = a.service.CreateAnswer(ctx, req.Data)
			if err != nil {
				switch err {
				case internal.ErrInvalidEvent,
					internal.ErrInvalidEventStatus,
					internal.ErrIDIsEmpty,
					internal.ErrInvalidData:
					ctx.JSON(http.StatusBadRequest, err.Error())
					return

				default:
					ctx.JSON(http.StatusInternalServerError, err.Error())
					return
				}
			}

			resp, err := toResponse(ans)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			ctx.JSON(http.StatusCreated, resp)
			return
		}

		err = a.service.UpdateAnswer(ctx, req.ID, "create", req.Data)
		if err != nil {
			switch err {
			case internal.ErrInvalidEvent,
				internal.ErrInvalidEventStatus,
				internal.ErrIDIsEmpty,
				internal.ErrInvalidData:
				ctx.JSON(http.StatusBadRequest, err.Error())
				return

			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

	}
}

// Get godoc
// @Summary      Show an event
// @Description  get event by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /answer/{id} [get]
func (a *AnswerHandler) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requestID
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		ans, err := a.service.GetAnswerByID(ctx, req.ID)
		if err != nil {
			switch err {
			case internal.ErrInvalidEvent,
				internal.ErrInvalidEventStatus,
				internal.ErrIDIsEmpty,
				internal.ErrAnswerNotFound,
				internal.ErrInvalidData:
				ctx.JSON(http.StatusBadRequest, err.Error())
				return

			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		resp, err := toResponse(ans)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// GetHistory godoc
// @Summary      Show a history event
// @Description  get history of event by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /answer/{id}/history [get]
func (a *AnswerHandler) GetHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requestID
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		ans, err := a.service.GetAnswerByID(ctx, req.ID)
		if err != nil {
			switch err {
			case internal.ErrInvalidEvent,
				internal.ErrInvalidEventStatus,
				internal.ErrIDIsEmpty,
				internal.ErrAnswerNotFound,
				internal.ErrInvalidData:
				ctx.JSON(http.StatusBadRequest, err.Error())
				return

			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
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
// Update godoc
// @Summary      Update an event
// @Description
// @Tags         answer
// @Accept       json
// @Produce      plain
// @Param        id   path      string  true  "ID"
// @Param        message  body  Request  true  "Request"
// @Success      200  {string}  string         "success"
// @Failure      400  {string}  string         "bad Request"
// @Failure      500  {string}  string         "fail"
// @Router       /answer/{id} [put]
func (a *AnswerHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqID requestID
		if err := ctx.ShouldBindUri(&reqID); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		var req Request
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
			switch err {
			case internal.ErrInvalidEvent,
				internal.ErrInvalidEventStatus,
				internal.ErrIDIsEmpty,
				internal.ErrInvalidData:
				ctx.JSON(http.StatusBadRequest, err.Error())
				return

			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		ctx.Status(http.StatusOK)
	}
}

// Delete returns an HTTP handler for answer creation.
// Delete godoc
// @Summary      Delete an event
// @Description  delete event by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Success      200  {string}  string         "success"
// @Failure      400  {string}  string         "bad Request"
// @Failure      500  {string}  string         "fail"
// @Router       /answer/{id} [delete]
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

func toResponse(answ *internal.Answer) (Response, error) {
	ev := answ.Events[len(answ.Events)-1]
	var data map[string]string

	err := json.Unmarshal(ev.RawData, &data)
	if err != nil {
		return Response{}, err
	}

	return Response{
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
