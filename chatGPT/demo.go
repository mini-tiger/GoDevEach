package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/0x9ef/openai-go"
	"os"
)

// go version 1.19
func main() {
	e := openai.New(os.Getenv("OPENAI_KEY"))
	r, err := e.Completion(context.Background(), &openai.CompletionOptions{
		// Choose model, you can see list of available models in models.go file
		Model: openai.ModelTextDavinci001,
		// Text to completion
		Prompt: []string{"Write a little bit of Wikipedia. What is that?"},
	})
	if err != nil {
		panic(err)
	}
	if b, err := json.MarshalIndent(r, "", "  "); err != nil {
		panic(err)
	} else {
		fmt.Println(string(b))
	}

	// Wikipedia is a free online encyclopedia, created and edited by volunteers.
	fmt.Println("What is the Wikipedia?", r.Choices[0].Text)
}
