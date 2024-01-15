package main

import (
	"fmt"
)

func test() {
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

	agent.setmessage(RoleUser, "write a very short response")

	fmt.Println(agent)
}
