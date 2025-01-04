package agent

// Package messages contains functions to edit an agent's message chain

import (
	// "AgentSmithU/agent"

	"fmt"
	"regexp"
	"sort"
	"strconv"
)

type Messages []Message

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PromptDefinition struct {
	Name        string
	Description string
	Parameters  string
}

func (m *Messages) Add(role, content string) {
	// messages := *m
	// messages = append(messages, Message{
	// 	Role:    role,
	// 	Content: content,
	// })
	// m = &messages
	*m = append(*m, Message{Role: role, Content: content})
}

// Makes numbered messages empty
// This 'clears' lines to '_'
// Deleting lines is a two-step process because I needed a way to
// keep track of what was deleted in the gui - it was basically a choice
// between reloading the page each delete (ie rewriting/resyncing html/
// backend index) or keeping a record of cleared lines and refreshing
// only on reload
func (m *Messages) Clearlines(editchoice string) error {
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
	// for i := len(sortednums) - 1; i >= 0; i-- {
	// 	agent.Messages = append(agent.Messages[:sortednums[i]], agent.Messages[sortednums[i]+1:]...)
	// }
	newmessages := *m

	for _, num := range sortednums {
		newmessages[num].Content = "_"
	}

	*m = newmessages

	return nil
}

// Deletes lines marked for deletion
// Clearlines marks lines for deletion with a '_' - Deletelines
// is used to actually remove them on page reload for gui
func (m *Messages) Deletelines() {
	messages := *m
	// remove empty messages
	// figure out what they are first
	var emptymessages []int
	// for i, item := range messages[1:] {
	for i, item := range messages {
		if item.Content == "_" {
			emptymessages = append(emptymessages, i)
		}
	}
	// sort the numbers and start from top
	sort.Ints(emptymessages)
	for i := len(emptymessages) - 1; i >= 0; i-- {
		messages = append(messages[:emptymessages[i]], messages[emptymessages[i]+1:]...)
		// messages = append(messages[:emptymessages[i]+1], messages[emptymessages[i]+2:]...)
	}
	*m = messages
}
