package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/business/create"
	"github.com/patriciabonaldy/bequest_challenge/internal/business/find"
	"github.com/patriciabonaldy/bequest_challenge/internal/business/update"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/command"
)

// Create godoc
// @Summary      Create an event
// @Description  if ID params is not empty will update the record other cases will create a new answer
// @Tags         answers
// @Accept       json
// @Produce      plain
// @Param        message  body  CreateRequest  true  "Request"
// @Success      200  {string}  string         "success"
// @Success      201  {string}  string         "success"
// @Failure      400  {string}  string         "bad Request"
// @Failure      500  {string}  string         "fail"
// @Router       /answers [post]
func Create(commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var err error
		if len(req.ID) == 0 {
			resul, err := commandBus.Dispatch(ctx, create.NewAnswerCommand("", req.Data))
			if err != nil {
				switch err {
				case internal.ErrInvalidEvent,
					internal.ErrInvalidEventStatus,
					internal.ErrIDIsEmpty,
					internal.ErrInvalidData,
					internal.ErrAnswerNotFound:
					ctx.JSON(http.StatusBadRequest, err.Error())
					return

				default:
					ctx.JSON(http.StatusInternalServerError, err.Error())
					return
				}
			}

			ans, ok := resul.(internal.Answer)
			if !ok {
				ctx.JSON(http.StatusInternalServerError, nil)
				return
			}

			ctx.JSON(http.StatusCreated, ans.ID)
			return
		}

		//err = a.service.UpdateAnswer(ctx, req.ID, "create", req.Data)
		_, err = commandBus.Dispatch(ctx, update.NewAnswerCommand(req.ID, "create", req.Data))
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

	}
}

// GetAnswer godoc
// @Summary      Show an event
// @Description  get event by ID example:"0bfce8da-bdc9-11ec-b9f3-acde48001122"
// @Tags         answers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /answers/{id} [get]
func GetAnswer(commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RequestID
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		//ans, err := a.service.GetAnswerByID(ctx, req.ID)
		resul, err := commandBus.Dispatch(ctx, find.NewAnswerCommand(req.ID))
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

		ans, ok := resul.(internal.Answer)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}

		resp, err := toResponse(&ans)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

// GetHistory godoc
// @Summary      Show a history event
// @Description  get history of event by ID example:"0bfce8da-bdc9-11ec-b9f3-acde48001122"
// @Tags         answers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /answers/{id}/history [get]
func GetHistory(commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RequestID
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		//ans, err := a.service.GetAnswerByID(ctx, req.ID)
		resul, err := commandBus.Dispatch(ctx, find.NewAnswerCommand(req.ID))
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

		ans, ok := resul.(internal.Answer)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}

		resp, err := toHistoryResponse(&ans)
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
// @Tags         answers
// @Accept       json
// @Produce      plain
// @Param        id   path      string  true  "ID"
// @Param        message  body  Request  true  "Request"
// @Success      200  {string}  string         "success"
// @Failure      400  {string}  string         "bad Request"
// @Failure      500  {string}  string         "fail"
// @Router       /answers/{id} [put]
func Update(commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqID RequestID
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

		//err := a.service.UpdateAnswer(ctx, reqID.ID, req.Event, req.Data)
		_, err := commandBus.Dispatch(ctx, update.NewAnswerCommand(reqID.ID, req.Event, req.Data))
		if err != nil {
			switch err {
			case internal.ErrInvalidEvent,
				internal.ErrInvalidEventStatus,
				internal.ErrIDIsEmpty,
				internal.ErrInvalidData,
				internal.ErrAnswerNotFound:
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
// @Description  delete event by ID example:"0bfce8da-bdc9-11ec-b9f3-acde48001122"
// @Tags         answers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /answers/{id} [delete]
func Delete(commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" || id == "0" {
			ctx.Status(http.StatusBadRequest)
			return
		}

		//err := a.service.UpdateAnswer(ctx, id, "delete", nil)
		_, err := commandBus.Dispatch(ctx, update.NewAnswerCommand(id, "delete", nil))
		if err != nil {
			switch err {
			case internal.ErrInvalidEvent,
				internal.ErrInvalidEventStatus,
				internal.ErrIDIsEmpty,
				internal.ErrInvalidData,
				internal.ErrAnswerNotFound:
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

func toHistoryResponse(answ *internal.Answer) (HistoryResponse, error) {
	resp := HistoryResponse{
		ID:     answ.ID,
		Events: []event{},
	}

	for _, ev := range answ.Events {
		var data map[string]string

		err := json.Unmarshal(ev.RawData, &data)
		if err != nil {
			return HistoryResponse{}, err
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
