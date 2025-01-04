package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"AgentSmithU/agent"
	"AgentSmithU/config"

	"github.com/atotto/clipboard"
)

func console(ag *agent.Agent) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	_, err := fmt.Println("Welcome")
	if err != nil {
		fmt.Println(err)
	} else {
		for {
			select {
			case <-interrupt:
				fmt.Println("Enter 'q' or 'quit' to exit!")
				continue
			default:
				fmt.Print("\nUser:\n")

				input := gettextinput()

				text := process_text(ag, input)

				if text == "" {
					continue
				}

				query := agent.Message{
					Role:    agent.RoleUser,
					Content: text,
				}

				ag.Messages = append(ag.Messages, query)

				for {
					// retries response until it works
					response, err := ag.Getresponse()
					if err != nil {
						fmt.Println(err)
						continue
					}

					estcost := (float64(ag.Tokencount) / 1000) * config.CallCost

					fmt.Println("\nAssistant:")
					fmt.Println(response.Content)

					fmt.Println("\nTokencount: ", ag.Tokencount, " Est. Cost: ", estcost)
					break
				}
			}
		}
	}
}

func process_text(ag *agent.Agent, text string) string {
	switch text {
	case "q", "quit":
		fmt.Println("\nQuitting...")
		os.Exit(0)
	case "del", "delete", "!":
		ag.Setprompt()
		fmt.Println("\nChat cleared!")
		return ""
	case "reset":
		config.Reset(ag)
		fmt.Println("\nChat reset!")
		return ""
	case "paste":
		text, err := clipboard.ReadAll()
		if err != nil {
			fmt.Println(err)
			return ""
		}
		fmt.Println("\nPasted text!")
		return text
	case "copy":
		response := ag.Messages[len(ag.Messages)].Content
		err := clipboard.WriteAll(response)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("\nCopied text!")
		return ""
	case "@", "sel", "select":
		printnumberlines(ag)
		fmt.Println("\nWhich lines would you like to delete?")
		editchoice := gettextinput()
		if editchoice == "" {
			return ""
		}
		ag.Messages.Clearlines(editchoice)
		// messages.Clearlines(agent, editchoice)
		// agent.Clearlines(editchoice)
		ag.Messages.Deletelines()
		// messages.Deletelines(agent)
		// agent.Deletelines()
		printnumberlines(ag)
		fmt.Println("Lines deleted!")
		return ""
	case "save":
		_, err := config.Save(ag.Messages, "Chats")
		if err != nil {
			fmt.Println(err)
		}
		return ""
	case "load":
		_, err := config.GetSaveFileList("Chats")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("What file would you like to load?")
		filename := gettextinput()
		if filename == "" {
			return ""
		}
		_, err = config.Load(ag, "Chats", filename)
		if err != nil {
			fmt.Println(err)
		}
		return ""
	case "prompt":
		fmt.Println("\nEnter new prompt:")
		input := gettextinput()
		if input == "paste" {
			text, err := clipboard.ReadAll()
			if err != nil {
				fmt.Println(err)
				return ""
			}
			ag.Prompt.Parameters = text
		} else {
			text := input
			ag.Prompt.Parameters = text

		}
		ag.Setprompt()
		fmt.Println("\nPrompt edited!")
		return ""
	case "help":
		fmt.Println("• Typing 'copy' will copy the last output from the bot\n• Typing 'paste' will paste your clipboard as a query - this way you can craft prompts in a text editor for multi-line queries\n• 'prompt' will let you enter in a new prompt ('paste' command works here)\n• 'save' will save the chat into a json file with the filename YYYYMMDDHHMM.txt\n• 'load <filename>' will load files\n• '@', 'sel', or 'select' will allow you to select lines to delete (handy if the chat is getting a bit long and you want to save on costs)\n• '!', 'del', or 'delete' will clear the chat log and start fresh\n'q' or 'quit' will quit the program")
		return ""
	default:
		// Nothing will happen
	}
	return text
}

func printnumberlines(ag *agent.Agent) {
	for i, msg := range ag.Messages {
		if msg.Role == agent.RoleUser {
			fmt.Printf("%d. User: %s\n", i, msg.Content)
		} else if msg.Role == agent.RoleAssistant {
			fmt.Printf("%d. Assistant: %s\n", i, msg.Content)
		}
	}
}
