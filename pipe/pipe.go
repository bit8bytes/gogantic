package pipe

import (
	"context"
	"log/slog"

	"github.com/bit8bytes/gogantic/llm"
	"github.com/bit8bytes/gogantic/output"
)

type Pipe[T any] struct {
	Messages     []llm.Message
	Model        llm.LLM
	OutputParser output.Parser[T]
}

type Options struct {
	logger *slog.Logger
}

func WithLogger(logger *slog.Logger) func(*Options) {
	return func(o *Options) {
		o.logger = logger
	}
}

func New[T any](messages []llm.Message, model llm.LLM, outputParser output.Parser[T], optsFunc ...func(*Options)) *Pipe[T] {
	opts := Options{}

	for _, fn := range optsFunc {
		fn(&opts)
	}

	return &Pipe[T]{
		Messages:     messages,
		Model:        model,
		OutputParser: outputParser,
	}
}

func (p *Pipe[T]) Invoke(ctx context.Context) (*T, error) {
	// 1. Input
	instructions := p.OutputParser.GetFormatInstructions()
	p.Messages[0].Content += " " + instructions

	// 2. LLM
	output, err := p.Model.GenerateContent(ctx, p.Messages)
	if err != nil {
		return nil, err
	}

	// 3. Output
	parsed, err := p.OutputParser.Parse(output.Result)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
