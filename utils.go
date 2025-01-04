package main

// mostly defunct
// used to contain functions that have been split into other packages

import (
	"encoding/json"
	"os"
)

func getrequest() map[string]string {
	// receive request from assistant
	// receives {"key": "string"} argument and converts it to map[string]string
	var req map[string]string
	args := os.Args[1]
	_ = json.Unmarshal([]byte(args), &req)
	return req
}

func getsubrequest(input string) map[string]string {
	// receives request from another function
	// receives {"key": "string"} argument and converts it to map[string]string
	var req map[string]string
	args := input
	_ = json.Unmarshal([]byte(args), &req)
	return req
}
