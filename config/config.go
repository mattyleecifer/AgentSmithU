package config

import (
	"agentsmithu/agent"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var HomeDir string // Home directory for storing agent files/folders /Prompts /Functions /Saves

var GuiFlag bool = false
var ConsoleFlag bool = false
var SaveChatName string

// var model string = "gpt-3.5-turbo"
var CallCost float64 = 0.002

var AuthString string
var AllowedIps []string
var AllowAllIps bool = false
var Port string = ":49327"

func GetFlags(ag *agent.Agent) {
	// Set default home dir
	HomeDir, _ = gethomedir()
	if HomeDir != "" {
		HomeDir = filepath.Join(HomeDir, "AgentSmith")
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
			ag.Api_key = arg
		case "-home":
			// Set home directory
			HomeDir = arg
		case "-save":
			// chats save to homeDir/Saves
			SaveChatName = arg
		case "-load":
			// load chat from homeDir/Saves
			Load(ag, "Chats", arg)
		case "-prompt":
			// Set prompt
			ag.Setprompt(arg)
		case "-model":
			// Set model
			ag.Model = arg
		case "-modelurl":
			// Set model
			ag.Modelurl = arg
		case "-maxtokens":
			// Change setting variable
			ag.Maxtokens, _ = strconv.Atoi(arg)
		case "-message":
			// Get the argument after the flag]
			// Set messages for the agent/create chat history
			ag.Messages.Add(agent.RoleUser, arg)
		case "-messageassistant":
			// Allows multiple messages with different users to be loaded in order
			ag.Messages.Add(agent.RoleAssistant, arg)
		case "--gui":
			// Run GUI
			GuiFlag = true
		case "-ip":
			// allow ip
			if arg == "all" {
				AllowAllIps = true
			} else {
				AllowedIps = append(AllowedIps, arg)
			}
		case "-auth":
			AuthString = arg
		case "-port":
			// change port
			Port = ":" + arg
		case "-allowallips":
			// allow all ips
			fmt.Println("Warning: Allowing all incoming connections.")
			AllowAllIps = true
		case "--console":
			// Run as console
			ConsoleFlag = true
		}
	}
}

func gethomedir() (string, error) {
	for _, item := range os.Args {
		if item == "-homedir" {
			HomeDir = item
		} else {
			usr, err := user.Current()
			if err != nil {
				fmt.Println("Failed to get current user:", err)
				return "", err
			}

			// Retrieve the path to user's home directory
			HomeDir = usr.HomeDir
		}
	}
	return HomeDir, nil
}

func Reset(ag *agent.Agent) {
	ag = agent.New()
	CallCost = 0.002
}
