package pipe

import (
	"context"

	"github.com/bit8bytes/gogantic/llm"
	"github.com/bit8bytes/gogantic/output"
)

type Pipe[T any] struct {
	Messages     []llm.Message
	Model        llm.LLM
	OutputParser output.Parser[T]
}

func New[T any](messages []llm.Message, model llm.LLM, outputParser output.Parser[T]) *Pipe[T] {
	return &Pipe[T]{
		Messages:     messages,
		Model:        model,
		OutputParser: outputParser,
	}
}

func (p *Pipe[T]) Invoke(ctx context.Context) (*T, error) {
	formatInstructions := p.OutputParser.GetFormatInstructions()
	p.Messages[0].Content += " " + formatInstructions

	output, err := p.Model.GenerateContent(ctx, p.Messages)
	if err != nil {
		return nil, err
	}

	parsedOutput, err := p.OutputParser.Parse(output.Result)
	if err != nil {
		return nil, err
	}
	return &parsedOutput, nil
}
