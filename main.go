package main

import (
	"AgentSmithU/agent"
	"AgentSmithU/config"
	"AgentSmithU/gui"
	"fmt"
)

// type Agent agent.Agent

func main() {
	// agent := agent.Agent{}
	a := agent.New()
	config.GetFlags(a)

	if config.GuiFlag {
		fmt.Println("Running GUI...")
		go console(a)
		gui.Gui(a)
	} else if config.ConsoleFlag {
		fmt.Println("Console only mode.")
		console(a)
	} else {
		response, err := a.Getresponse()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response.Content)
		if config.SaveChatName != "" {
			config.Save(a.Messages, "Chats")
		}
	}
}
