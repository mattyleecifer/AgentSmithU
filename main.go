package main

import (
	"agentsmithu/agent"
	"agentsmithu/config"
	"agentsmithu/console"
	"agentsmithu/gui"
	"fmt"
)

func main() {
	a := agent.New()
	config.GetFlags(a)

	if config.GuiFlag {
		fmt.Println("Running GUI...")
		go console.Console(a)
		gui.Gui(a)
	} else if config.ConsoleFlag {
		fmt.Println("Console only mode.")
		console.Console(a)
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
