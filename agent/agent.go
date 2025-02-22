// Package agent provides core components to create and run an agent
package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

var defaultprompt = PromptDefinition{
	Name:        "Default",
	Description: "Default Prompt",
	Parameters:  "You are a helpful assistant. Please generate truthful, accurate, and honest responses while also keeping your answers succinct and to-the-point. Do no assume the user is correct. If the user is wrong then state so plainly along with reasoning. Use step-by-step reasoning when generating a response. Show your working. Today's date is: ",
}

var defaultmodel string = "llama3.2"

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

type Agent struct {
	Prompt     PromptDefinition
	Tokencount int
	Api_key    string
	Model      string
	Modelurl   string
	Maxtokens  int
	Messages   Messages
	Functions  Functions
}

type RequestBody struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	Stream     bool      `json:"stream"`
	Max_tokens int       `json:"max_tokens"`
}

type ChatResponse struct {
	// ID string `json:"id"`
	// Object  string   `json:"object"`
	// Created int64    `json:"created"`
	// Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type ChatResponseOllama struct {
	Message    Message `json:"message"`
	Eval_count int     `json:"eval_count"`
}

type ChatResponseAnthropic struct {
	Content []ContentAnthropic `json:"content"`
	Usage   UsageAnthropic     `json:"usage"`
}

type ContentAnthropic struct {
	Text string `json:"text"`
}

type UsageAnthropic struct {
	Input_tokens  int `json:"input_tokens"`
	Output_tokens int `json:"output_tokens"`
}

type Choice struct {
	// Index   int     `json:"index"`
	Message Message `json:"message"`
	// FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	// PromptTokens     int `json:"prompt_tokens"`
	TotalTokens int `json:"total_tokens"`
	// CompletionTokens int `json:"completion_tokens"`
}

func (agent *Agent) GetmodelURL() string {
	// to be expanded
	var url string
	if agent.Modelurl == "" {
		switch {
		case strings.Contains(agent.Model, "mistral"):
			url = "https://api.mistral.ai/v1/chat/completions"
		case strings.Contains(agent.Model, "gpt"):
			url = "https://api.openai.com/v1/chat/completions"
		case strings.Contains(agent.Model, "claude"):
			url = "https://api.anthropic.com/v1/messages"
		default:
			// handle local models here
			url = "http://localhost:11434/api/chat"
		}
		return url
	}
	return agent.Modelurl
}

// Creates new Agent with default settings
func New() *Agent {
	var today = time.Now().Format("January 2, 2006")
	agent := Agent{}

	// Set prompt
	agent.Prompt = defaultprompt
	agent.Prompt.Parameters += today
	agent.Setprompt()

	// Set max tokens
	agent.Maxtokens = 2048

	// Set model
	agent.Model = defaultmodel

	// Set Tokencount
	agent.Tokencount = 0

	return &agent
}

// Sets prompt - note that this does not change the rest of the messages in a conversation
func (agent *Agent) Setprompt(prompt ...string) {
	if len(agent.Messages) == 0 {

		// RoleAssistant, not RoleSystem here because some models can't handle it
		agent.Messages.Add(RoleAssistant, "")
	}
	if len(prompt) == 0 {
		agent.Messages[0].Content = agent.Prompt.Parameters
	} else {
		agent.Messages[0].Content = prompt[0]
	}
}

func (agent *Agent) Getresponse() (Message, error) {
	var response Message

	modelurl := agent.GetmodelURL()
	parsedURL, err := url.Parse(modelurl)
	if err != nil {
		fmt.Println("Error parsing URL:", err) // Handle error accordingly
	}

	if strings.Contains(parsedURL.Host, "anthropic") {
		// Anthropic doesn't allow system role and roles must alternate between user/assistant
		// This breaks things so this snippet changes the system to user and adds a dummy assistant message
		if len(agent.Messages) == 2 {
			agent.Messages[0].Role = RoleUser
		}
		// checks for double role occurences and adds a dummy message in between
		// works backwards cause poppin in values probably isn't healthy going upwards
		for index := len(agent.Messages) - 1; index >= 1; index-- {
			if agent.Messages[index].Role == agent.Messages[index-1].Role {
				dummyMessage := Message{
					Role:    RoleAssistant,
					Content: "_",
				}
				agent.Messages = slices.Insert(agent.Messages, index, dummyMessage)
			}
		}
	}

	// Create the request body
	requestBody := &RequestBody{
		Model:      agent.Model,
		Messages:   agent.Messages,
		Stream:     false,
		Max_tokens: agent.Maxtokens,
	}

	// Encode the request body as JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error encoding request body:", err)
		return response, err
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, modelurl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return response, err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", agent.Api_key))

	// Anthropic-specific headers
	if strings.Contains(parsedURL.Host, "anthropic") {
		req.Header["x-api-key"] = []string{agent.Api_key}
		req.Header["content-type"] = []string{"application/json"}
		req.Header["anthropic-version"] = []string{"2023-06-01"}
	}

	// fmt.Println(req)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return response, err
	}

	// Handle the HTTP response
	defer resp.Body.Close()

	// fmt.Println(resp)

	// process the prompt and get response

	// For ollama based models
	if strings.Contains(parsedURL.Host, "localhost") {
		var chatresponse ChatResponseOllama
		err = json.NewDecoder(resp.Body).Decode(&chatresponse)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)
			return response, err
		}

		fmt.Println(chatresponse)

		response = chatresponse.Message
		// Print the decoded message
		fmt.Println("Decoded message:", response.Content)

		agent.Tokencount = chatresponse.Eval_count

		// Add message to chain for Agent
		agent.Messages = append(agent.Messages, response)
	} else if strings.Contains(parsedURL.Host, "anthropic") {
		var chatresponse ChatResponseAnthropic
		err = json.NewDecoder(resp.Body).Decode(&chatresponse)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)
			return response, err
		}

		fmt.Println(chatresponse)

		// Print the decoded message
		fmt.Println("Decoded message:", chatresponse.Content[0].Text)

		agent.Tokencount = chatresponse.Usage.Input_tokens + chatresponse.Usage.Output_tokens

		response = Message{
			Role:    RoleAssistant,
			Content: chatresponse.Content[0].Text,
		}

		// Add message to chain for Agent
		agent.Messages = append(agent.Messages, response)
	} else {
		var chatresponse ChatResponse

		// copy resp.body so can use it multiple times
		body, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(body))

		err = json.NewDecoder(resp.Body).Decode(&chatresponse)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)

		}

		if len(chatresponse.Choices) == 0 {
			fmt.Println("Error with response:", chatresponse)
			// attempt to use local llm to decode
			// convert the JSON to string
			// but first turn it into a map
			var jsonResponse interface{}

			// revive resp.body
			resp.Body = io.NopCloser(bytes.NewBuffer(body))

			err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
			if err != nil {
				fmt.Println("Error decoding JSON response:", err)
			}
			jsonData, err := json.Marshal(jsonResponse)
			if err != nil {
				panic(err)
			}
			jsonStr := string(jsonData)

			fmt.Println("jsonStr", jsonStr)

			// send the string to the converter and receive chatresponse
			chatresponse, err = agentAPIConverter(jsonStr)
			if err != nil {
				return response, err
			}
		}

		fmt.Println(chatresponse)

		response = chatresponse.Choices[0].Message

		// Print the decoded message
		fmt.Println("Decoded message:", response.Content)

		agent.Tokencount = chatresponse.Usage.TotalTokens

		// Add message to chain for Agent
		agent.Messages = append(agent.Messages, response)
	}

	// Check if there is a function call and then deal with it
	if strings.HasPrefix(response.Content, "**functioncall") {
		fmt.Println("functioncall detected", response.Content)
	}

	return response, nil
}

// experimental
func agentAPIConverter(jsonStr string) (ChatResponse, error) {
	var chatresponse ChatResponse // convert response to text

	// create local converter agent and set variables
	converter := New()
	converter.Model = "llama3.2" // any ollama llm should work, can even convert this to openai/mistral/anthropic
	converter.Modelurl = "http://localhost:11434/api/chat"
	converter.Maxtokens = 2048
	converter.Setprompt(`Extract the text/message data from any inputs. Output only the text/message data without any commentary. Do not change anything. Output the text/message data exactly as it is written in the original data`)

	// attempt to get response convertered
	converter.Messages.Add(RoleUser, jsonStr)

	response, err := converter.Getresponse()
	if err != nil {
		fmt.Println("failed to convert", err)
		return chatresponse, err
	}

	// put the extracted response into a new message and return
	newMessage := Message{
		Content: response.Content,
		Role:    RoleAssistant,
	}
	newChoice := Choice{
		Message: newMessage,
	}
	chatresponse.Choices = append(chatresponse.Choices, newChoice)

	return chatresponse, nil
}
