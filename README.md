# Interacting with LLMs in Go has never been easier.

Meet Gogo the Giant Gopher.

<p align="center"> <img src="gogantic-mascot.png" alt="Gogantic Mascot" width="250"/></p>

Gogo helps you work with LLMs in Go(lang) — without external dependencies.
Gogo speeds up your interactions with LLMs while keeping your stack lean and efficient.

## 🚴🏽‍♂️ Roadmap

- ✅ Build core components to interact with large language models.
- ✅ Create an Agent that can interact with the outside world using tools.
- 🔜 Develop a Director Agent that manages complex tasks by coordinating multiple Agents.

Next up: We’ll experiment with intraction between the host system and a local LLM.

Bonus: Gogantic includes a simple interface for adding documents to a Qdrant vector store.

## Examples:

Usage of the `core/pipe`

```go
// This is not the full example. See 'examples/core/pipe'
pipe := pipe.New(messages, ollamaClient, parser)
result, _ := pipe.Invoke(context.Background())
fmt.Println("Translate from", result.InputLanguage, " to ", result.OutputLanguage)
fmt.Println("Result: ", result.Text)
```

Go to [Examples](/EXAMPLES.md) for more info.
You also can fork the repo and run `make examples/core/pipe` (requires ollama and llama3:8b model).

## 📥 Contributions Welcome

Contributions are very welcome! Whether it's fixing a typo, suggesting an improvement, or adding new simplified streamlined features — feel free to open an issue or a pull request.

## 📚 Sources and Inspiration

- [tmc/langchaingo](https://github.com/tmc/langchaingo)

## ✨ Contributors

<table>
  <tr>
    <td align="center"><a href="https://github.com/tobiasgleiter"><img src="https://avatars.githubusercontent.com/tobiasgleiter" width="100px;" alt=""/><br /><sub><b>Tobias Gleiter</b></sub></a></td>
  </tr>
</table>

Contributions of any kind are welcome! 🙌
