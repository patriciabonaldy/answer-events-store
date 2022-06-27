package update

import (
	"context"

	"github.com/patriciabonaldy/bequest_challenge/internal/platform/command"
	"github.com/pkg/errors"
)

const AnswerCommandType command.Type = "command.update.answer"

// AnswerCommand is the command dispatched to create a new answer.
type AnswerCommand struct {
	eventID   string
	eventType string
	data      map[string]string
}

// NewAnswerCommand creates a new AnswerCommand.
func NewAnswerCommand(id, eventType string, data map[string]string) AnswerCommand {
	return AnswerCommand{
		eventID:   id,
		eventType: eventType,
		data:      data,
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
	updateAnswerCmd, ok := cmd.(AnswerCommand)
	if !ok {
		return nil, errors.New("unexpected command")
	}

	return nil, h.service.UpdateAnswer(
		ctx,
		updateAnswerCmd.eventID,
		updateAnswerCmd.eventType,
		updateAnswerCmd.data,
	)
}
