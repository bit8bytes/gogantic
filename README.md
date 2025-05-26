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

## Example:

Usage of the `pipe`

```go
// This is not the full example. See 'examples/pipe'
pipe := pipe.New(messages, ollamaClient, parser)
result, _ := pipe.Invoke(context.Background())
fmt.Println("Translate from", result.InputLanguage, " to ", result.OutputLanguage)
fmt.Println("Result: ", result.Text)
```

Go to [Examples](/docs/EXAMPLES.md) for more info.

## 📚 Sources and Inspiration

- [tmc/langchaingo](https://github.com/tmc/langchaingo)

## ✨ Contributors

<table>
  <tr>
    <td align="center">
      <img src="https://avatars.githubusercontent.com/tobiasgleiter" width="64px" style="border-radius: 50%;" alt="Contributor Avatar"/>
    </td>
     <td align="center">
      <img src="https://avatars.githubusercontent.com/u/79313705" width="64px" style="border-radius: 50%;" alt="Contributor Avatar"/>
    </td>
    <td align="center">
      <img src="https://avatars.githubusercontent.com/u/184933573" width="64px" style="border-radius: 50%;" alt="Contributor Avatar"/>
    </td>
  
  </tr>
</table>

Contributions of any kind are welcome! 🙌 See [Get Involved](/docs/GET-INVOLVED.md) to get started.
