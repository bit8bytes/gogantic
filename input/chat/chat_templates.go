package chat

import (
	"bytes"
	"text/template"

	"github.com/bit8bytes/gogantic/llm"
)

type PromptTemplate struct {
	Messages []llm.Message
}

func New(messages []llm.Message) (*PromptTemplate, error) {
	return &PromptTemplate{Messages: messages}, nil
}

func (cpt *PromptTemplate) Format(data interface{}) ([]llm.Message, error) {
	var formattedMessages []llm.Message

	for _, templat := range cpt.Messages {
		tmpl, err := template.New("prompt").Parse(templat.Content)
		if err != nil {
			return nil, err
		}

		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, data)
		if err != nil {
			return nil, err
		}

		formattedMessages = append(formattedMessages, llm.Message{
			Role:    templat.Role,
			Content: buffer.String(),
		})
	}

	return formattedMessages, nil
}
