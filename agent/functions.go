package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Functions []Function

type Function struct {
	Name        string
	Description string
	Parameters  string
}

// adds function to []Function
func (agent *Agent) AddFunction(function Function) error {
	for _, name := range agent.Functions {
		if name.Name == function.Name {
			return fmt.Errorf("Function with same name already exists")
		}
	}
	agent.Functions = append(agent.Functions, function)
	agent.SetFunctionPrompt()
	return nil
}

// removes function from []Function
func (agent *Agent) RemoveFunction(function string) {
	for index, item := range agent.Functions {
		if item.Name == function {
			agent.Functions = append(agent.Functions[:index], agent.Functions[index+1:]...)
		}
	}
	agent.SetFunctionPrompt()
}

// detects if function is being called and then extracts the function and runs it if approved
func (agent *Agent) RunFunction(function Function) Message {
	// runs function on system
	data, err := json.Marshal(function.Parameters)
	if err != nil {
		fmt.Println(err)
	}
	cmd := strings.ToLower(function.Name)
	arg1, _ := strconv.Unquote(string(data))
	// unq := strconv.Unquote(string(data))
	// arg1 := string(data)

	// fmt.Println("\nFunction call: ", functiondef.Name)
	fmt.Println("\nCommand: ", arg1)

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
	}

	runPath := filepath.Join(currentDir, cmd)

	exec := exec.Command(runPath, arg1)
	output, err := exec.CombinedOutput()
	if err != nil {
		log.Println(err)
		output = []byte(err.Error())
	}

	fmt.Println("Function Output:\n", string(output))

	var response Message
	response.Content = string(output)
	response.Role = RoleAssistant
	return response
}

// func (agent *Agent) loadFunction(filename string) (Function, error) {
// 	var newfunction Function

// 	filedata, err := loadfile("Functions", filename)
// 	if err != nil {
// 		return newfunction, err
// 	}

// 	err = json.Unmarshal(filedata, &newfunction)
// 	if err != nil {
// 		return newfunction, err
// 	}

// 	return newfunction, nil
// }
