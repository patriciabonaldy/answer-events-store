package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"

	"github.com/patriciabonaldy/bequest_challenge/internal"
)

// AnswerDB is a structure of answers to be stored
type AnswerDB struct {
	AnswerID string              `bson:"answer_id"`
	Events   []EventDB           `bson:"events"`
	CreateAt primitive.Timestamp `bson:"createdAt,omitempty"`
	UpdateAt primitive.Timestamp `bson:"updatedAt,omitempty"`
}

func (a *AnswerDB) createAt() time.Time {
	return time.Unix(int64(a.CreateAt.T), 0).UTC()
}

func (a *AnswerDB) updateAt() time.Time {
	return time.Unix(int64(a.UpdateAt.T), 0).UTC()
}

// EventDB is a structure of events to be stored
type EventDB struct {
	EventID   string              `bson:"event_id"`
	Type      string              `bson:"event_type"`
	RawData   bson.Raw            `bson:"data"`
	Timestamp primitive.Timestamp `timestamp:"createdAt" bson:"timestamp"`
	Version   int                 `bson:"version"`
}

func NewEvent(id, eventType string, data []byte, version int) (EventDB, error) {
	createdAt := primitive.Timestamp{
		T: uint32(time.Now().Unix()),
	}

	return EventDB{
		EventID:   id,
		Type:      eventType,
		RawData:   bson.Raw{Kind: data[0], Data: data},
		Timestamp: createdAt,
		Version:   version,
	}, nil
}

func (c *EventDB) ID() string {
	return c.EventID
}

func (c *EventDB) Name() string {
	return c.Type
}

func (c *EventDB) At() time.Time {
	return time.Unix(int64(c.Timestamp.T), 0).UTC()
}

func (c *EventDB) Data() []byte {
	return c.RawData.Data
}

func (c *EventDB) Unmarshall(out interface{}) error {
	return c.RawData.Unmarshal(out)
}

func parseToBusinessAnswer(result AnswerDB) internal.Answer {
	answer := internal.Answer{
		ID:       result.AnswerID,
		CreateAt: result.createAt(),
		UpdateAt: result.updateAt(),
	}

	for _, ev := range result.Events {
		event := internal.NewEvent(ev.EventID, internal.EventType(ev.Type), ev.Data())
		event.Timestamp = ev.At()
		answer.AddEvent(event)
	}
	return answer
}

func parseToAnswerDB(result internal.Answer) AnswerDB {
	answer := AnswerDB{
		AnswerID: result.ID,
		CreateAt: primitive.Timestamp{
			T: uint32(result.CreateAt.Unix()),
		},
		UpdateAt: primitive.Timestamp{
			T: uint32(result.UpdateAt.Unix()),
		},
	}

	for _, ev := range result.Events {
		answer.Events = append(answer.Events, EventDB{
			EventID: ev.EventID,
			Type:    string(ev.Type),
			RawData: bson.Raw{Kind: ev.RawData[0], Data: ev.RawData},
			Timestamp: primitive.Timestamp{
				T: uint32(ev.Timestamp.Unix()),
			},
			Version: 0,
		})
	}
	return answer
}
