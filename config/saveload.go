package config

import (
	"agentsmithu/agent"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Save(data interface{}, filetype string, input ...string) (string, error) {
	// savetype must be Chats, Prompts, or Functions

	var filename string
	if len(input) == 0 {
		currentTime := time.Now()
		filename = currentTime.Format("20060102150405")
	} else {
		filename = strings.Replace(input[0], " ", "_", -1)
	}

	var filedir string
	if strings.HasSuffix(filename, ".json") {
		filedir = filepath.Join(HomeDir, filetype, filename)
	} else {
		filedir = filepath.Join(HomeDir, filetype, filename+".json")
	}
	appDir := filepath.Join(HomeDir, filetype)
	err := os.MkdirAll(appDir, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to create app directory:", err)
		return "", err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(filedir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return "", err
	}

	return filedir, nil
}

func Load(ag *agent.Agent, filetype string, filename string) ([]byte, error) {

	var filedir string
	if strings.HasSuffix(filename, ".json") {
		filedir = filepath.Join(HomeDir, filetype, filename)
	} else {
		filedir = filepath.Join(HomeDir, filetype, filename+".json")
	}

	file, err := os.Open(filedir)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	switch filetype {
	case "Chats":
		Reset(ag)
		newmessages := agent.Messages{}
		err = json.Unmarshal(data, &newmessages)
		if err != nil {
			return nil, err
		}
		ag.Messages = newmessages
		return nil, err
	case "Functions":
		return data, nil
	case "Prompts":
		return data, nil
	}
	return nil, nil
}

func Delete(filetype, filename string) error {
	var filedir string
	if strings.HasSuffix(filename, ".json") {
		filedir = filepath.Join(HomeDir, filetype, filename)
	} else {
		filedir = filepath.Join(HomeDir, filetype, filename+".json")
	}

	err := os.Remove(filedir)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return err
	}

	fmt.Println("File deleted successfully.")

	return nil
}

func GetSaveFileList(filetype string) ([]string, error) {
	// Create a directory for your app
	savepath := filepath.Join(HomeDir, filetype)
	files, err := os.ReadDir(savepath)
	if err != nil {
		return nil, err
	}
	var res []string

	fmt.Println("\nFiles:")

	for _, file := range files {
		filename := strings.ReplaceAll(file.Name(), ".json", "")
		res = append(res, filename)
		fmt.Println(file.Name())
	}

	return res, nil
}
