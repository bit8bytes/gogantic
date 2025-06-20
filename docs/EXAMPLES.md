# Examples

## Usage of Pipe: Input, LLM, Output

Core enables simple interaction with LLMs. The concept is simple:

`Input -> LLM -> Output`

You just have to prepare the messages, add a model and define the output (e.g. json)

```go
// This is not the full example. See 'examples/pipe'
pipe := pipe.New(messages, ollamaClient, parser)
result, _ := pipe.Invoke(context.Background())
fmt.Println("Translate from", result.InputLanguage, " to ", result.OutputLanguage)
fmt.Println("Result: ", result.Text)
```

Install using `go get github.com/bit8bytes/gogantic/pipe`
