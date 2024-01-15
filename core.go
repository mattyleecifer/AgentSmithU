package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// This contains all the core functions - it is designed so it can be copied into any other project to create new agents that can receive/send back to the original agent or to any other program through json input/string output

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

var homeDir string // Home directory for storing agent files/folders /Prompts /Functions /Saves

var guiFlag bool = false
var consoleFlag bool = false
var savechatName string

// var model string = "gpt-3.5-turbo"
var defaultmodel string = "mistral-tiny"
var callcost float64 = 0.002
var maxtokens int = 2048

var authstring string
var allowedIps []string
var allowAllIps bool = false
var port string = ":49327"

var pubkey = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
jNOy+K77PpegsYbpFJ7B2+EBkQ==
=TZBa
-----END PGP PUBLIC KEY BLOCK-----`

var privkey = `-----BEGIN PGP PRIVATE KEY BLOCK-----

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
=XYLc
-----END PGP PRIVATE KEY BLOCK-----`

func (agent *Agent) getmodelURL() string {
	// to be expanded
	var url string
	switch {
	case strings.HasPrefix(agent.model, "mistral"):
		url = "https://api.mistral.ai/v1/chat/completions"
	case strings.HasPrefix(agent.model, "gpt"):
		url = "https://api.openai.com/v1/chat/completions"
	default:
		// handle invalid model here
		fmt.Println("Error: Invalid model")
	}
	return url
}

func newAgent(key ...string) Agent {
	agent := Agent{}
	agent.prompt = defaultprompt
	agent.setprompt()
	agent.model = defaultmodel
	agent.tokencount = 0
	agent.getflags()
	if agent.api_key == "" {
		if len(key) == 0 {
			agent.getkey()
		}
	}
	return agent
}

func (agent *Agent) getflags() {
	// Set default home dir
	homeDir, _ = gethomedir()
	if homeDir != "" {
		homeDir = filepath.Join(homeDir, "AgentSmith")
	}

	// range over args to get flags
	for index, flag := range os.Args {
		var arg string
		if index < len(os.Args)-1 {
			item := os.Args[index+1]
			if !strings.HasPrefix(item, "-") {
				arg = item
			}
		}

		switch flag {
		case "-key":
			// Set API key
			agent.api_key = arg
		case "-home":
			// Set home directory
			homeDir = arg
		case "-save":
			// chats save to homeDir/Saves
			savechatName = arg
		case "-load":
			// load chat from homeDir/Saves
			agent.loadfile("Chats", arg)
		case "-prompt":
			// Set prompt
			agent.setprompt(arg)
		case "-model":
			// Set model
			defaultmodel = arg
		case "-maxtokens":
			// Change setting variable
			maxtokens, _ = strconv.Atoi(arg)
		case "-message":
			// Get the argument after the flag]
			// Set messages for the agent/create chat history
			agent.setmessage(RoleUser, arg)
		case "-messageassistant":
			// Allows multiple messages with different users to be loaded in order
			agent.setmessage(RoleAssistant, arg)
		case "--gui":
			// Run GUI
			guiFlag = true
		case "-ip":
			// allow ip
			if arg == "all" {
				allowAllIps = true
			} else {
				allowedIps = append(allowedIps, arg)
			}
		case "-auth":
			authstring = arg
		case "-port":
			// change port
			port = ":" + arg
		case "-allowallips":
			// allow all ips
			fmt.Println("Warning: Allowing all incoming connections.")
			allowAllIps = true
		case "--console":
			// Run as console
			consoleFlag = true
		}
	}
}

func gettextinput() string {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if len(input) == 0 {
			return ""
		}
		return input
	}
	return ""
}

func (agent *Agent) reset() {
	*agent = newAgent()
	callcost = 0.002
	maxtokens = 2048
}

func (agent *Agent) setmessage(role, content string) {
	message := Message{
		Role:    role,
		Content: content,
	}
	agent.Messages = append(agent.Messages, message)
}

func (agent *Agent) setprompt(prompt ...string) {
	agent.Messages = []Message{}
	if len(prompt) == 0 {
		agent.setmessage(RoleSystem, agent.prompt.Parameters)
	} else {
		agent.setmessage(RoleSystem, prompt[0])
	}
	agent.tokencount = 0
}

func (agent *Agent) getresponse() (Message, error) {
	var response Message

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

	if len(chatresponse.Choices) == 0 {
		fmt.Println("Error with response:", chatresponse)
		return response, err
	}

	fmt.Println(chatresponse)

	// Print the decoded message
	fmt.Println("Decoded message:", chatresponse.Choices[0].Message.Content)

	agent.tokencount = chatresponse.Usage.TotalTokens

	// Add message to chain for Agent
	agent.Messages = append(agent.Messages, chatresponse.Choices[0].Message)

	return chatresponse.Choices[0].Message, nil
}

func gethomedir() (string, error) {
	for _, item := range os.Args {
		if item == "-homedir" {
			homeDir = item
		} else {
			usr, err := user.Current()
			if err != nil {
				fmt.Println("Failed to get current user:", err)
				return "", err
			}

			// Retrieve the path to user's home directory
			homeDir = usr.HomeDir
		}
	}
	return homeDir, nil
}

func (agent *Agent) getkey() {
	filePath := filepath.Join(homeDir, "apikey.txt")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("\nEnter OpenAI key: ")
		key := gettextinput()

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Failed to create file:", err)
			os.Exit(0)
		}
		defer file.Close()

		// fmt.Println("File created successfully!")

		armor, _ := helper.EncryptMessageArmored(pubkey, key)

		_, err = file.WriteString(armor)
		if err != nil {
			fmt.Println("Failed to write to file:", err)
			os.Exit(0)
		}

		agent.api_key = key

		// fmt.Println("API key set.")
	} else {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Failed to read file:", err)
			os.Exit(0)
		}

		decrypted, err := helper.DecryptMessageArmored(privkey, nil, string(content))
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		agent.api_key = decrypted

		// fmt.Println("API key set.")
	}
}

func getrequest() map[string]string {
	// receive request from assistant
	// receives {"key": "string"} argument and converts it to map[string]string
	var req map[string]string
	args := os.Args[1]
	_ = json.Unmarshal([]byte(args), &req)
	return req
}

func getsubrequest(input string) map[string]string {
	// receives request from another function
	// receives {"key": "string"} argument and converts it to map[string]string
	var req map[string]string
	args := input
	_ = json.Unmarshal([]byte(args), &req)
	return req
}

func (agent *Agent) savefile(data interface{}, filetype string, input ...string) (string, error) {
	// savetype must be Chats, Prompts, or Functions

	var filename string
	if len(input) == 0 {
		currentTime := time.Now()
		filename = currentTime.Format("20060102150405")
	} else {
		filename = strings.Replace(input[0], " ", "_", -1)
	}

	var filedir string
	if strings.HasSuffix(filename, ".json") {
		filedir = filepath.Join(homeDir, filetype, filename)
	} else {
		filedir = filepath.Join(homeDir, filetype, filename+".json")
	}
	appDir := filepath.Join(homeDir, filetype)
	err := os.MkdirAll(appDir, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to create app directory:", err)
		return "", err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(filedir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return "", err
	}

	return filedir, nil
}

func (agent *Agent) loadfile(filetype string, filename string) ([]byte, error) {

	var filedir string
	if strings.HasSuffix(filename, ".json") {
		filedir = filepath.Join(homeDir, filetype, filename)
	} else {
		filedir = filepath.Join(homeDir, filetype, filename+".json")
	}

	file, err := os.Open(filedir)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	switch filetype {
	case "Chats":
		agent.reset()
		newmessages := []Message{}
		err = json.Unmarshal(data, &newmessages)
		if err != nil {
			return nil, err
		}
		agent.Messages = newmessages
		return nil, err
	case "Functions":
		return data, nil
	case "Prompts":
		return data, nil
	}
	return nil, nil
}

func deletefile(filetype, filename string) error {
	var filedir string
	if strings.HasSuffix(filename, ".json") {
		filedir = filepath.Join(homeDir, filetype, filename)
	} else {
		filedir = filepath.Join(homeDir, filetype, filename+".json")
	}

	err := os.Remove(filedir)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return err
	}

	fmt.Println("File deleted successfully.")

	return nil
}

func getsavefilelist(filetype string) ([]string, error) {
	// Create a directory for your app
	savepath := filepath.Join(homeDir, filetype)
	files, err := os.ReadDir(savepath)
	if err != nil {
		return nil, err
	}
	var res []string

	fmt.Println("\nFiles:")

	for _, file := range files {
		filename := strings.ReplaceAll(file.Name(), ".json", "")
		res = append(res, filename)
		fmt.Println(file.Name())
	}

	return res, nil
}

func (agent *Agent) deletelines(editchoice string) error {
	// Use regular expression to find all numerical segments in the input string
	reg := regexp.MustCompile("[0-9]+")
	nums := reg.FindAllString(editchoice, -1)

	var sortednums []int
	// Convert each numerical string to integer and sort
	for _, numStr := range nums {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return err
		}
		sortednums = append(sortednums, num)
	}

	sort.Ints(sortednums)

	fmt.Println("Deleting lines: ", sortednums)

	// go from highest to lowest to not fu the order
	for i := len(sortednums) - 1; i >= 0; i-- {
		agent.Messages = append(agent.Messages[:sortednums[i]], agent.Messages[sortednums[i]+1:]...)
	}

	return nil
}
