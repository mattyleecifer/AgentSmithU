package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	. "AgentSmithU/agent"

	"github.com/atotto/clipboard"
)

func console(agent *Agent) {
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

				text := process_text(agent, input)

				if text == "" {
					continue
				}

				query := Message{
					Role:    RoleUser,
					Content: text,
				}

				agent.Messages = append(agent.Messages, query)

				for {
					// retries response until it works
					response, err := agent.Getresponse()
					if err != nil {
						fmt.Println(err)
						continue
					}

					estcost := (float64(agent.Tokencount) / 1000) * callcost

					fmt.Println("\nAssistant:")
					fmt.Println(response.Content)

					fmt.Println("\nTokencount: ", agent.Tokencount, " Est. Cost: ", estcost)
					break
				}
			}
		}
	}
}

func process_text(agent *Agent, text string) string {
	switch text {
	case "q", "quit":
		fmt.Println("\nQuitting...")
		os.Exit(0)
	case "del", "delete", "!":
		agent.Setprompt()
		fmt.Println("\nChat cleared!")
		return ""
	case "reset":
		Reset(agent)
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
		response := agent.Messages[len(agent.Messages)].Content
		err := clipboard.WriteAll(response)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("\nCopied text!")
		return ""
	case "@", "sel", "select":
		printnumberlines(agent)
		fmt.Println("\nWhich lines would you like to delete?")
		editchoice := gettextinput()
		if editchoice == "" {
			return ""
		}
		agent.Clearlines(editchoice)
		agent.Deletelines()
		printnumberlines(agent)
		fmt.Println("Lines deleted!")
		return ""
	case "save":
		_, err := savefile(agent.Messages, "Chats")
		if err != nil {
			fmt.Println(err)
		}
		return ""
	case "load":
		_, err := getsavefilelist("Chats")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("What file would you like to load?")
		filename := gettextinput()
		if filename == "" {
			return ""
		}
		_, err = loadfile(agent, "Chats", filename)
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
			agent.Prompt.Parameters = text
		} else {
			text := input
			agent.Prompt.Parameters = text

		}
		agent.Setprompt()
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

func printnumberlines(agent *Agent) {
	for i, msg := range agent.Messages {
		if msg.Role == RoleUser {
			fmt.Printf("%d. User: %s\n", i, msg.Content)
		} else if msg.Role == RoleAssistant {
			fmt.Printf("%d. Assistant: %s\n", i, msg.Content)
		}
	}
}
