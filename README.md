# Interacting with LLMs in Go has never been easier.

## Core

Core enables simple interaction with LLMs. The concept is simple:

`Input -> Model -> Output`

You just have to prepare the messages, add a model and define the output (e.g. json)

```go
pipe := pipe.NewPipe(messages, ollamaClient, parser)
result, _ := pipe.Invoke(context.Background())
fmt.Println("Translate from", result.InputLanguage, " to ", result.OutputLanguage)
fmt.Println("Result: ", result.Text)
```
