# Examples

## 1. Core (Input, LLM, Output)

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

## 2. Agent

`Tools -> Agent -> Executor -> Final Result`

Agents can interact with tools and get informationen from the outside world. The following example shows an agent that can access the current temperature.

```go
// This is not a full example. See 'examples/agent/ollama/temperature'
tools := map[string]tool.Tool{
		"CurrentTemperatureInFahrenheit": CurrentTemperatureInFahrenheit{},
		"FormatFahrenheitToCelsius":      FormatFahrenheitToCelsius{},
	}

weatherAgent := agents.New(llm, tools)
weatherAgent.Task("1. What is the temperature outside? 2. What is the temperature in Celsius?")

executor := agents.NewExecutor(weatherAgent)
executor.Run(context.TODO())

finalAnswer, _ := weatherAgent.GetFinalAnswer()
fmt.Println(finalAnswer)
```

Run `go get github.com/bit8bytes/gogantic/agent` to install the agent.

## 3. Director Agents (Coming soon)

Now that Agents (2) can call tools, we are able to create an Director Agent that can call Agents.
