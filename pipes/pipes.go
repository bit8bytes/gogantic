package pipes

import (
	"context"

	"github.com/bit8bytes/gogantic/llms"
)

type parser[T any] interface {
	Parse(text string) (T, error)
	Instructions() string
}

type llm interface {
	Generate(ctx context.Context, messages []llms.Message) (*llms.ContentResponse, error)
}

type Pipe[T any] struct {
	messages []llms.Message
	model    llm
	parser   parser[T]
}

func New[T any](messages []llms.Message, model llm, parser parser[T]) *Pipe[T] {
	return &Pipe[T]{
		messages: messages,
		model:    model,
		parser:   parser,
	}
}

func (pipe *Pipe[T]) Invoke(ctx context.Context) (*T, error) {
	instructions := pipe.parser.Instructions()
	pipe.messages[0].Content += " " + instructions

	output, err := pipe.model.Generate(ctx, pipe.messages)
	if err != nil {
		return nil, err
	}

	parsed, err := pipe.parser.Parse(output.Result)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
