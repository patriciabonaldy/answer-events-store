package handler

import (
	"time"
)

type createRequest struct {
	ID   string            `uri:"id,omitempty"`
	Data map[string]string `json:"data" binding:"required"`
}

type requestID struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type request struct {
	Event string            `json:"event" binding:"required"`
	Data  map[string]string `json:"data" binding:"required"`
}

type response struct {
	ID       string            `json:"answer_id"`
	Event    string            `json:"event"`
	Data     map[string]string `json:"data"`
	CreateAt time.Time         `json:"createdAt,omitempty"`
}

type historyResponse struct {
	ID     string  `json:"id"`
	Events []event `json:"events"`
}

type event struct {
	EventID   string            `json:"event_id"`
	Type      string            `json:"event_type"`
	Data      map[string]string `json:"data""`
	Timestamp time.Time         `timestamp:"createdAt"`
	Version   int               `json:"version"`
}
