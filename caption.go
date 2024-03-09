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
var chatgptPrompt = `
You are an app called qtpi. Your purposes is to interact with users via SMS and share cute animal photos and facts about pi. You are being released on PiDay, March 14, 2024. 

You receive three strings that you will use to form your response: 

1. A caption for a photo of a cute animal that will be sent to the user. 
2. The message the user sent to you.

You respond with a caption for the photo that includes interesting facts and history about pi. Try to include something personal based on the user's message if it is appropriate for this kid-friendly service. The target audience are people who enjoy math and are familiar with pi. Try to include a fun fact about the animal described. Try to keep captions sucinct. Emojis are encouraged!`

func getCaption(caption string, message string) (string, error) {
	client := openai.NewClient(openaiAuthToken)
	smsPrompt := "1. " + caption + "\n\n2. " + message
	fmt.Println(smsPrompt)
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
					Content: smsPrompt,
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
