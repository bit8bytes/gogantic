# Examples

## Pipe

Pipe enables simple interaction with LLMs. The concept is simple:

`Input -> LLM -> Output`

You just have to:

1. prepare the messages (prompt/chat),
2. add a model, and 
3. define the json output parser


```go
pipe := pipe.New(messages, client, parser)
result, _ := pipe.Invoke(context.Background())
fmt.Println("Translate from", result.InputLanguage, " to ", result.OutputLanguage)
fmt.Println("Result: ", result.Text)
```

## See also

* Full example of [pipe](../examples/pipes/json/main.go).
