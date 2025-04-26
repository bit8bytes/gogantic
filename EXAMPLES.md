# Examples

## 1. Core

Core enables simple interaction with LLMs. The concept is simple:

`Input -> Model -> Output`

You just have to prepare the messages, add a model and define the output (e.g. json)

```go
// This is not the full example. See 'examples/core/pipe'
pipe := pipe.New(messages, ollamaClient, parser)
result, _ := pipe.Invoke(context.Background())
fmt.Println("Translate from", result.InputLanguage, " to ", result.OutputLanguage)
fmt.Println("Result: ", result.Text)
```

Install using `go get github.com/bit8bytes/gogantic/core/pipe`

## 2. Agents

`Tools -> Agent -> Executor -> Final Result`

Agents can interact with tools and get informationen from the outside world. The following example shows an agent that can access the current temperature.

```go
// This is not a full example. See 'examples/agents/temperature'
tools := map[string]agents.Tool{
    "CurrentTemperatureInFahrenheit": CurrentTemperatureInFahrenheit{},
}

weatherAgent := agents.New(llm, tools)
weatherAgent.Task("What is the temperature outside?")

executor := agents.NewExecutor(weatherAgent)
executor.Run(context.TODO())

finalAnswer, _ := weatherAgent.GetFinalAnswer()
fmt.Println(finalAnswer)
```

Run `go get github.com/bit8bytes/gogantic/agents` to install the agents.

## 3. Director Agents (Coming soon)

Now that Agents (2) can call tools, we are able to create an Director Agent that can call Agents.
