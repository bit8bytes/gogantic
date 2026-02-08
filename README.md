# Interacting with LLMs in Go has never been easier.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Test](https://github.com/bit8bytes/gogantic/actions/workflows/test_cmd_app.yml/badge.svg)](https://github.com/bit8bytes/gogantic/actions/workflows/tests.yml)

Meet Gogo the standing chick. 🐥.

* Gogo helps you work with LLMs in Go(lang) — without external dependencies
* Gogo keeps your stack lean and efficient
* Gogo can be run everywhere using minimal resources pointing to a LLM

## Example:

Usage of the `pipe`

```go
// This is not the full example. See 'examples/pipe'
pipe := pipe.New(messages, ollamaClient, parser)
result, _ := pipe.Invoke(context.Background())
fmt.Println("Translate from", result.InputLanguage, " to ", result.OutputLanguage)
fmt.Println("Result: ", result.Text)
```

## 📚 Sources and Inspiration

Note, these inspirations have different goals then Gogo but are worth looking into.

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
    <td align="center">
      <img src="https://avatars.githubusercontent.com/u/173567119" width="64px" style="border-radius: 50%;" alt="Contributor Avatar"/>
    </td>
  </tr>
</table>

Contributions of any kind are welcome! 🙌 See [Get Involved](/docs/GET-INVOLVED.md) to get started.
