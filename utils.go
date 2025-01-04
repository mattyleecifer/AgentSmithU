package main

// mostly defunct
// used to contain functions that have been split into other packages

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
