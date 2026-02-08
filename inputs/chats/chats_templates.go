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

func (chat *chats) Messages() []llms.Message {
	return chat.messages
}
