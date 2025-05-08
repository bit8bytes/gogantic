# Interacting with LLMs in Go has never been easier.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Meet Gogo the Giant Gopher.

<p align="center"> <img src="/docs/img/gogantic-mascot.png" alt="Gogantic Mascot" width="250"/></p>

Gogo helps you work with LLMs in Go(lang) — without external dependencies.
Gogo speeds up your interactions with LLMs while keeping your stack lean and efficient.

## 🚴🏽‍♂️ Roadmap

- ✅ Build core components to interact with large language models.
- ✅ Create an Agent that can interact with the outside world using tools.
- 🔜 Develop a Director Agent that manages complex tasks by coordinating multiple Agents.

Next up: We’ll experiment with intraction between the host system and a local LLM.

Bonus: Gogantic includes a simple interface for adding documents to a Qdrant vector store.

## Example:

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

## 📚 Sources and Inspiration

- [tmc/langchaingo](https://github.com/tmc/langchaingo)

## ✨ Contributors

<table>
  <tr>
    <td align="center">
      <a href="https://github.com/tobiasgleiter">
        <img src="https://avatars.githubusercontent.com/tobiasgleiter" width="100px" style="border-radius: 50%;" alt="Contributor Avatar"/>
      </a>
    </td>
  </tr>
</table>

Contributions of any kind are welcome! 🙌 See [Get Involved](/GET-INVOLVED.md) to get started.
