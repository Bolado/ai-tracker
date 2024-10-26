package ai

import (
	"context"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

func Summarize(text string) (string, error) {
	token := os.Getenv("OPENAI_API_TOKEN")

	client := openai.NewClient(token)

	//limit text to 50000 characters
	if len(text) > 50000 {
		text = text[:50000]
	}

	// Create a context with a 1-minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	resp, err := client.CreateChatCompletion(
		ctx,
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
