// mostly defunct
package main

import (
	"encoding/json"
	"os"
)

// func getkey(agent *Agent) {
// 	filePath := filepath.Join(homeDir, "apikey.txt")

// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 		fmt.Println("\nEnter OpenAI key: ")
// 		key := gettextinput()

// 		file, err := os.Create(filePath)
// 		if err != nil {
// 			fmt.Println("Failed to create file:", err)
// 			os.Exit(0)
// 		}
// 		defer file.Close()

// 		// fmt.Println("File created successfully!")

// 		armor, _ := helper.EncryptMessageArmored(pubkey, key)

// 		_, err = file.WriteString(armor)
// 		if err != nil {
// 			fmt.Println("Failed to write to file:", err)
// 			os.Exit(0)
// 		}

// 		agent.Api_key = key

// 		// fmt.Println("API key set.")
// 	} else {
// 		content, err := os.ReadFile(filePath)
// 		if err != nil {
// 			fmt.Println("Failed to read file:", err)
// 			os.Exit(0)
// 		}

// 		decrypted, err := helper.DecryptMessageArmored(privkey, nil, string(content))
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(0)
// 		}

// 		agent.Api_key = decrypted

// 		// fmt.Println("API key set.")
// 	}
// }

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

// func savefile(data interface{}, filetype string, input ...string) (string, error) {
// 	// savetype must be Chats, Prompts, or Functions

// 	var filename string
// 	if len(input) == 0 {
// 		currentTime := time.Now()
// 		filename = currentTime.Format("20060102150405")
// 	} else {
// 		filename = strings.Replace(input[0], " ", "_", -1)
// 	}

// 	var filedir string
// 	if strings.HasSuffix(filename, ".json") {
// 		filedir = filepath.Join(HomeDir, filetype, filename)
// 	} else {
// 		filedir = filepath.Join(HomeDir, filetype, filename+".json")
// 	}
// 	appDir := filepath.Join(HomeDir, filetype)
// 	err := os.MkdirAll(appDir, os.ModePerm)
// 	if err != nil {
// 		fmt.Println("Failed to create app directory:", err)
// 		return "", err
// 	}

// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return "", err
// 	}

// 	file, err := os.OpenFile(filedir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer file.Close()

// 	_, err = file.Write(jsonData)
// 	if err != nil {
// 		return "", err
// 	}

// 	return filedir, nil
// }

// func loadfile(agent *Agent, filetype string, filename string) ([]byte, error) {

// 	var filedir string
// 	if strings.HasSuffix(filename, ".json") {
// 		filedir = filepath.Join(HomeDir, filetype, filename)
// 	} else {
// 		filedir = filepath.Join(HomeDir, filetype, filename+".json")
// 	}

// 	file, err := os.Open(filedir)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()
// 	data, err := io.ReadAll(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	switch filetype {
// 	case "Chats":
// 		Reset(agent)
// 		newmessages := []Message{}
// 		err = json.Unmarshal(data, &newmessages)
// 		if err != nil {
// 			return nil, err
// 		}
// 		agent.Messages = newmessages
// 		return nil, err
// 	case "Functions":
// 		return data, nil
// 	case "Prompts":
// 		return data, nil
// 	}
// 	return nil, nil
// }

// func deletefile(filetype, filename string) error {
// 	var filedir string
// 	if strings.HasSuffix(filename, ".json") {
// 		filedir = filepath.Join(HomeDir, filetype, filename)
// 	} else {
// 		filedir = filepath.Join(HomeDir, filetype, filename+".json")
// 	}

// 	err := os.Remove(filedir)
// 	if err != nil {
// 		fmt.Println("Error deleting file:", err)
// 		return err
// 	}

// 	fmt.Println("File deleted successfully.")

// 	return nil
// }

// func getsavefilelist(filetype string) ([]string, error) {
// 	// Create a directory for your app
// 	savepath := filepath.Join(HomeDir, filetype)
// 	files, err := os.ReadDir(savepath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res []string

// 	fmt.Println("\nFiles:")

// 	for _, file := range files {
// 		filename := strings.ReplaceAll(file.Name(), ".json", "")
// 		res = append(res, filename)
// 		fmt.Println(file.Name())
// 	}

// 	return res, nil
// }
