package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

type promptDefinition struct {
	Name        string
	Description string
	Parameters  string
}

var today = time.Now().Format("January 2, 2006")

var defaultprompt = promptDefinition{
	Name:        "Default",
	Description: "Default Prompt",
	Parameters:  "You are a helpful assistant. Please generate truthful, accurate, and honest responses while also keeping your answers succinct and to-the-point. Today's date is: " + today,
}

type Agent struct {
	prompt     promptDefinition
	tokencount int
	api_key    string
	model      string
	Messages   []Message
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

func (agent *Agent) getmodelURL() string {
	var url string
	switch agent.model {
	case "mistral-tiny":
		url = "https://api.mistral.ai/v1/chat/completions"
	case "gpt-3.5-turbo":
		url = "https://api.openai.com/v1/chat/completions"
	default:
		// handle invalid model here
	}
	return url
}

func (agent *Agent) getresponse(role, message string) (Message, error) {
	var response Message

	newmessage := Message{
		Role:    role,
		Content: message,
	}

	// Attach to agent
	agent.Messages = append(agent.Messages, newmessage)

	// Create the request body
	requestBody := &RequestBody{
		Model:    agent.model,
		Messages: agent.Messages,
	}

	// Encode the request body as JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error encoding request body:", err)
		return response, err
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, agent.getmodelURL(), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return response, err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", agent.api_key))

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return response, err
	}

	// Handle the HTTP response
	defer resp.Body.Close()

	// Decode the response body into a Message object
	var chatresponse ChatResponse
	err = json.NewDecoder(resp.Body).Decode(&chatresponse)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return response, err
	}

	// Print the decoded message
	fmt.Println("Decoded message:", chatresponse.Choices[0].Message.Content)

	agent.tokencount = chatresponse.Usage.TotalTokens

	// Add message to chain for Agent
	agent.Messages = append(agent.Messages, chatresponse.Choices[0].Message)

	return chatresponse.Choices[0].Message, nil
}

func main() {
	// agent := Agent{
	// 	prompt:   defaultprompt,
	// 	api_key:  "NIBmDn0Ds7IUi985yIkmEHd6lixw9pss",
	// 	model:    "mistral-tiny",
	// 	Messages: []Message{},
	// }
	agent := Agent{
		prompt:   defaultprompt,
		api_key:  "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",
		model:    "gpt-3.5-turbo",
		Messages: []Message{},
	}
	// set prompt
	prompt := Message{
		Role:    RoleSystem,
		Content: defaultprompt.Parameters,
	}

	agent.Messages = append(agent.Messages, prompt)

	_, err := agent.getresponse("user", "write a very short response")
	if err != nil {
		// handle the error here
		fmt.Println("An error occurred:", err)
	}

	fmt.Println(agent)
}
