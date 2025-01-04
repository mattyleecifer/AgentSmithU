package config

import (
	"AgentSmithU/agent"
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

var pubkey = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
jNOy+K77PpegsYbpFJ7B2+EBkQ==
=TZBa
-----END PGP PUBLIC KEY BLOCK-----`

var privkey = `-----BEGIN PGP PRIVATE KEY BLOCK-----

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
=XYLc
-----END PGP PRIVATE KEY BLOCK-----`

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
			ag.Messages.Set(agent.RoleUser, arg)
			// messages.Set(agent, RoleUser, arg)
			// agent.Setmessage(RoleUser, arg)
		case "-messageassistant":
			// Allows multiple messages with different users to be loaded in order
			ag.Messages.Set(agent.RoleAssistant, arg)
			// messages.Set(agent, RoleAssistant, arg)
			// agent.Setmessage(RoleAssistant, arg)
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
