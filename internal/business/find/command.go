package find

import (
	"context"

	"github.com/patriciabonaldy/bequest_challenge/internal"
	"github.com/patriciabonaldy/bequest_challenge/internal/platform/command"
	"github.com/pkg/errors"
)

const AnswerCommandType command.Type = "command.find.answer"

// AnswerCommand is the command dispatched to create a new answer.
type AnswerCommand struct {
	eventID string
}

// NewAnswerCommand creates a new AnswerCommand.
func NewAnswerCommand(id string) AnswerCommand {
	return AnswerCommand{
		eventID: id,
	}
}

func (c AnswerCommand) Type() command.Type {
	return AnswerCommandType
}

// AnswerCommandHandler is the command handler
// responsible for creating answers.
type AnswerCommandHandler struct {
	service Service
}

// NewAnswerCommandHandler initializes a new AnswerCommandHandler.
func NewAnswerCommandHandler(service Service) AnswerCommandHandler {
	return AnswerCommandHandler{
		service: service,
	}
}

// Handle implements the command.Handler interface.
func (h AnswerCommandHandler) Handle(ctx context.Context, cmd command.Command) (interface{}, error) {
	createAnswerCmd, ok := cmd.(AnswerCommand)
	if !ok {
		return internal.Answer{}, errors.New("unexpected command")
	}

	answer, err := h.service.GetAnswerByID(
		ctx,
		createAnswerCmd.eventID,
	)

	if err != nil {
		return internal.Answer{}, err
	}

	return *answer, nil
}
