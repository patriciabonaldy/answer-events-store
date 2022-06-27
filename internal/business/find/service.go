package find

import (
	"context"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/logger"
)

type Service interface {
	GetAnswerByID(ctx context.Context, eventID string) (*internal.Answer, error)
}

type service struct {
	repository internal.Storage
	log        logger.Logger
}

// NewService returns the default Service interface implementation.
func NewService(repository internal.Storage, log logger.Logger) Service {
	return &service{
		repository: repository,
		log:        log,
	}
}

func (s service) GetAnswerByID(ctx context.Context, eventID string) (*internal.Answer, error) {
	answer, err := s.repository.GetByID(ctx, eventID)
	if err != nil {
		s.log.Errorf("error GetAnswerByID ID:%s:%s", eventID, err.Error())
		return nil, err
	}

	return &answer, nil
}
