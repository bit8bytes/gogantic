package input

import (
	"bytes"
	"text/template"

	"github.com/bit8bytes/gogantic/core/models"
)

type ChatPromptTemplate struct {
	Messages []models.MessageContent
}

func NewChatPromptTemplate(messages []models.MessageContent) (*ChatPromptTemplate, error) {
	return &ChatPromptTemplate{Messages: messages}, nil
}

func (cpt *ChatPromptTemplate) FormatMessages(data interface{}) ([]models.MessageContent, error) {
	var formattedMessages []models.MessageContent

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

		formattedMessages = append(formattedMessages, models.MessageContent{
			Role:    templat.Role,
			Content: buffer.String(),
		})
	}

	return formattedMessages, nil
}
