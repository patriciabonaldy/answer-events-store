package handler

import (
	"time"
)

// swagger:model CreateRequest
type CreateRequest struct {
	// the id for create a new event
	ID string `uri:"id,omitempty"`

	// data is answers list
	Data map[string]string `json:"data" binding:"required"`
}

type requestID struct {
	ID string `uri:"id" binding:"required,uuid" example:"0bfce8da-bdc9-11ec-b9f3-acde48001122"`
}

// swagger:model Request
type Request struct {
	// type of event
	Event string `json:"event" binding:"required"`
	// data is answers list
	Data map[string]string `json:"data" binding:"required" example:"en:Map,ru:Карта,kk:Карталар"`
}

// swagger:model Response
type Response struct {
	// event id
	ID string `json:"answer_id"`
	// type of event
	Event string `json:"event"`
	// data is answers list
	Data     map[string]string `json:"data" example:"en:Map,ru:Карта,kk:Карталар"`
	CreateAt time.Time         `json:"createdAt,omitempty"`
}

type historyResponse struct {
	ID     string  `json:"id"`
	Events []event `json:"events"`
}

type event struct {
	EventID   string            `json:"event_id"`
	Type      string            `json:"event_type"`
	Data      map[string]string `json:"data"`
	Timestamp time.Time         `timestamp:"createdAt"`
	Version   int               `json:"version"`
}
