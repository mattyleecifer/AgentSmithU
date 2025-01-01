package main

import (
	"AgentSmithU/agent"
	"fmt"
)

// type Agent agent.Agent

func main() {
	// agent := agent.Agent{}
	a := agent.New()
	getflags(&a)

	if guiFlag {
		fmt.Println("Running GUI...")
		go console(&a)
		gui(&a)
	} else if consoleFlag {
		fmt.Println("Console only mode.")
		console(&a)
	} else {
		response, err := a.Getresponse()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response.Content)
		if savechatName != "" {
			savefile(a.Messages, "Chats")
		}
	}
}
