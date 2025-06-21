package pipe

import (
	"context"
	"log/slog"
	"time"

	"github.com/bit8bytes/gogantic/llm"
	"github.com/bit8bytes/gogantic/output"
)

type Pipe[T any] struct {
	Messages     []llm.Message
	Model        llm.LLM
	OutputParser output.Parser[T]
	logger       *slog.Logger
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
	opts := Options{
		logger: slog.Default(),
	}

	for _, fn := range optsFunc {
		fn(&opts)
	}

	return &Pipe[T]{
		Messages:     messages,
		Model:        model,
		OutputParser: outputParser,
		logger:       opts.logger,
	}
}

func (p *Pipe[T]) Invoke(ctx context.Context) (*T, error) {
	// 1. Input
	start := time.Now()
	instructions := p.OutputParser.GetFormatInstructions()
	p.Messages[0].Content += " " + instructions
	inputDebugGroup := slog.Group("input",
		"messages", p.Messages,
		"instructions", instructions,
		"duration_ms", time.Since(start).Microseconds(),
	)
	p.logger.DebugContext(ctx, "Input", inputDebugGroup)

	// 2. LLM
	start = time.Now()
	output, err := p.Model.GenerateContent(ctx, p.Messages)
	if err != nil {
		return nil, err
	}

	p.logger.InfoContext(ctx, "Generated content from LLM", "preview", output.Result)
	debugOutputGroup := slog.Group("output",
		"length", len(output.Result),
		"duration_ms", time.Since(start).Milliseconds(),
	)
	p.logger.DebugContext(ctx, "Generated content from LLM", debugOutputGroup)

	// 3. Output
	start = time.Now()
	parsed, err := p.OutputParser.Parse(output.Result)
	if err != nil {
		return nil, err
	}

	parsedOutputGroup := slog.Group("parsed",
		"generic", parsed,
		"duration_ms", time.Since(start).Milliseconds(),
	)
	p.logger.DebugContext(ctx, "Successfully parsed output", parsedOutputGroup)

	return &parsed, nil
}
