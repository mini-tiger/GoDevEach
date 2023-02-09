package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

const (
	openAIAPI = "https://api.openai.com/v1/completions"
	model     = "text-davinci-003"
	apikey    = "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"
)

// Request struct for making API request
type Request struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
	//TopP        float32  `json:"top_p"`
	//N           int      `json:"n"`
	//Stream      bool     `json:"stream"`
	//Stop        string   `json:"stop"`
	//Context     []string `json:"context"`
	Model string `json:"model"`
}

// Response struct to parse API response
type Response struct {
	ID         string   `json:"id"`
	Choices    []Choice `json:"choices"`
	StopTokens []string `json:"stop_tokens"`
}

// Choice struct to parse individual text completions
type Choice struct {
	Text string `json:"text"`
}

func main() {
	// Store context between turns
	var context []string

	// Start conversation
	fmt.Println("ChatGPT: Hello! How can I help you today?")

	// Read user input
	var input string
	fmt.Scanln(&input)

	// Continuously loop until user inputs "exit"
	for input != "exit" {
		// Make API request
		resp, err := makeRequest(input, context)
		if err != nil {
			fmt.Println("Error making API request:", err)
			break
		}
		fmt.Println(resp)
		// Update context with API response
		context = append(context, resp.Choices[0].Text)

		// Print API response
		fmt.Println("ChatGPT:", resp.Choices[0].Text)

		// Read next user input
		fmt.Scanln(&input)
	}

	fmt.Println("ChatGPT: Goodbye! Have a great day.")
}
func makeRequest(input string, context1 []string) (Response, error) {
	// Create API request body
	req := Request{
		Prompt:      input,
		MaxTokens:   1024,
		Model:       model,
		Temperature: 0.5,
		//TopP:        0.9,
		//N:           1,
		//Stream:      false,
		//Context:     context1,
	}
	fmt.Println(input)
	//var result = &Response{}
	client := resty.New()
	resp, err := client.R().
		SetBody(req).
		//SetResult(result). // or SetResult(AuthSuccess{}).
		SetAuthToken(apikey).
		//SetError(&AuthError{}).       // or SetError(AuthError{}).
		Post(openAIAPI)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	return Response{}, nil

	//return response, nil
}
