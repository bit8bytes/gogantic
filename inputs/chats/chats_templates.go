// Package chats provides templated chat message handling for LLM interactions.
package chats

import (
	"bytes"
	"text/template"

	"github.com/bit8bytes/gogantic/llms"
)

type chats struct {
	messages  []llms.Message
	templates []*template.Template
}

// New creates a new chats instance by parsing the content of each message as a
// Go text/template. It panics if any message content fails to parse.
func New(messages []llms.Message) *chats {
	templates := make([]*template.Template, 0, len(messages))
	for _, message := range messages {
		tmpl, err := template.New("prompts").Parse(message.Content)
		if err != nil {
			panic(err)
		}

		templates = append(templates, tmpl)
	}

	return &chats{
		messages:  messages,
		templates: templates,
	}
}

// Execute applies the given data to each template and returns the resulting
// messages with their content replaced by the executed template output.
func (chat *chats) Execute(data any) ([]llms.Message, error) {
	for i, template := range chat.templates {
		buffer := new(bytes.Buffer)
		if err := template.Execute(buffer, data); err != nil {
			return nil, err
		}

		chat.messages[i] = llms.Message{
			Role:    chat.messages[i].Role,
			Content: buffer.String(),
		}
	}

	return chat.messages, nil
}

// Messages returns the current list of messages.
func (chat *chats) Messages() []llms.Message {
	return chat.messages
}
