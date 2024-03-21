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
var chatgptPrompt = "You are an app called qtpi. Your purposes is to interact with users via SMS and share cute animal photos and facts about pi."
var searchPrompt = "You receive one string: the message the user sent to you. You respond with a single search term for a cute animal photo with no quotes. If the user mentions a specific animal, try to include that animal name in the search term. Your response will be used verbatim to search an API so it must be a single query. Your response must always start with an animal name. This is a kid-friendly service and all search terms should be kid-friendly. If the user sends an inappropriate message, simply ignore it and return a search term about a specific type of animal."
var captionPrompt = `You receive two strings that you will use to form your response: 

1. A caption for a photo of a cute animal that will be sent to the user. 
2. The message the user sent to you.

You respond with a caption for the photo that includes interesting facts and history about pi. Try to include something personal based on the user's message if it is appropriate for this kid-friendly service. If the message is inappropriate, tell the user you don't appreciate that. The target audience are people who enjoy math and are familiar with pi. Try to include a fun fact about the animal described. Try to keep captions sucinct. Emojis are encouraged! If the message from the user isn't about math or animals, prompt them for an animal they want to learn about.`

func getSearchTerm(message string) (string, error) {
	client := openai.NewClient(openaiAuthToken)
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: chatgptPrompt + "\n\n" + searchPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	if err != nil {
		return defaultCaption, fmt.Errorf("error creating search term completion: %w", err)
	}

	// Check for an empty response.Choices array or an empty response.Choices[0].Message.Content
	if len(response.Choices) == 0 || len(response.Choices[0].Message.Content) <= 1 {
		return defaultCaption, nil
	}

	searchTerm := response.Choices[0].Message.Content

	fmt.Println("-----\n\n" + "User mesage: " + message)
	fmt.Println("Search term: " + searchTerm)

	// Use regex to remove all quotes from the string
	regex := regexp.MustCompile(`["]+`)
	searchTerm = regex.ReplaceAllString(searchTerm, "")

	return searchTerm, nil
}

func getCaption(caption string, message string) (string, error) {
	client := openai.NewClient(openaiAuthToken)
	smsPrompt := captionPrompt + "\n\n1. " + caption + "\n\n2. " + message
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
		return defaultCaption, fmt.Errorf("error creating caption completion: %w", err)
	}

	// Check for an empty response.Choices array or an empty response.Choices[0].Message.Content
	if len(response.Choices) == 0 || len(response.Choices[0].Message.Content) <= 1 {
		return defaultCaption, nil
	}

	customCaption := response.Choices[0].Message.Content

	fmt.Println("Caption: " + customCaption + "\n\n-----")

	// Use regex to remove all quotes from the string
	regex := regexp.MustCompile(`["]+`)
	customCaption = regex.ReplaceAllString(customCaption, "")

	return customCaption, nil
}
