package ai

import (
	"context"
	"os"

	"github.com/sashabaranov/go-openai"
)

func Summarize(text string) (string, error) {
	token := os.Getenv("OPENAI_API_TOKEN")

	client := openai.NewClient(token)

	//limit text to 50000 characters
	if len(text) > 50000 {
		text = text[:50000]
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Summarize the following article using a maximum of 240 characters.",
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
