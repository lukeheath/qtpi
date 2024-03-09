package main

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/sashabaranov/go-openai"
)

var openaiAuthToken = os.Getenv("OPENAI_AUTH_TOKEN")
var defaultCaption = "Much like the cuteness of animals is constant, the number π is a mathematical constant that is the ratio of a circle's circumference to its diameter."
var chatgptPrompt = "You receive captions of photos of cute animals. You respond with a caption for the photo that includes interesting facts and history about pi. The target audience are people who enjoy math and are familiar with pi. Try to include a fun fact about the animal described. Try to keep captions sucinct. Emojis are encouraged!"

func getCaption(caption string) (string, error) {
	client := openai.NewClient(openaiAuthToken)
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: chatgptPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: caption,
				},
			},
		},
	)

	if err != nil {
		return defaultCaption, fmt.Errorf("error creating chat completion: %w", err)
	}

	// Check for an empty response.Choices array or an empty response.Choices[0].Message.Content
	if len(response.Choices) == 0 || len(response.Choices[0].Message.Content) <= 1 {
		return defaultCaption, nil
	}

	customCaption := response.Choices[0].Message.Content

	// Use regex to trim quotes from the start and end of the customCaption
	regex := regexp.MustCompile(`^["']+|["']+$`)
	customCaption = regex.ReplaceAllString(customCaption, "")

	return customCaption, nil
}
