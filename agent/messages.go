package agent

// Package messages contains functions to edit an agent's message chain

import (
	// "AgentSmithU/agent"

	"fmt"
	"regexp"
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

// Add new message to message chain
func (m *Messages) Add(role, content string) {
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

	// Convert each numerical string to integer and turn the corresponding message to '_'
	for _, numStr := range nums {
		if num, err := strconv.Atoi(numStr); err == nil && num < len(*m) {
			(*m)[num].Content = "_"
			fmt.Println("Clearing line: ", num)
		}
	}

	return nil
}

// Deletes lines marked for deletion
// Clearlines marks lines for deletion with a '_' - Deletelines
// is used to actually remove them on page reload for gui
func (m *Messages) Deletelines() {
	messages := *m
	// Create a new slice to store non-empty messages
	nonEmptyMessages := make([]Message, 0, len(messages))
	// Append all non-empty messages to the new slice
	for _, msg := range messages {
		if msg.Content != "_" {
			nonEmptyMessages = append(nonEmptyMessages, msg)
		}
	}
	// Replace the original slice with the new slice
	*m = nonEmptyMessages
}
