package handler

import (
	"encoding/json"
	"github.com/nexus-telegram/NexusSDK/types"
	"os"
)

// LoadConfig reads the configuration file (config.json) from the specified file path
// and parses its contents into a Config struct.
//
// This function ensures the configuration settings required for the application,
// such as proxy details and API key, are properly loaded and ready for use.
//
// # Parameters:
//   - filePath: The path to the configuration file (e.g., "config.json").
//
// # Returns:
//   - types.Config: A struct containing the parsed configuration data.
//   - error: An error if the file cannot be opened, read, or parsed.
//
// # Example config.json:
//
//	{
//		"proxy": {
//			"ip": "192.168.1.100",
//			"port": 8080,
//			"username": "proxyuser",
//			"password": "proxypass",
//			"socksType": 5,
//			"timeout": 30
//		},
//		"api_key": "your-api-key-here"
//	}
//
// # Example Usage:
//
//	config, err := LoadConfig("config.json")
//	if err != nil {
//		log.Fatalf("Failed to load configuration: %v", err)
//	}
//	fmt.Println(config.APIKey) // Output: your-api-key-here
//
// # Notes:
//   - Ensure the file at the specified path exists and is properly formatted as JSON.
//   - If the file contains invalid JSON or cannot be accessed, an error is returned.
func LoadConfig(filePath string) (types.Config, error) {
	var config types.Config
	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

// LoadAccounts reads the accounts.json file and parses it into a slice of Account structs.
//
// # Parameters:
//   - filePath: The path to the accounts.json file.
//
// # Returns:
//   - []types.Account: A slice of Account structs parsed from the file.
//   - error: An error if the file cannot be opened, read, or parsed.
//
// # Example accounts.json:
//
//	[
//		{
//			"game-data": "user=%7B%22id%22%3A78894796...",
//			"telegram": {
//				"tdataStringSession": "session-data-string",
//				"appId": "123456",
//				"appHash": "abcdef123456",
//				"telegramId": "987654321"
//			}
//		},
//		{
//			"game-data": "user=%7B%22id%22%3A78894796...",
//			"telegram": {
//				"tdataStringSession": "another-session-data",
//				"appId": "654321",
//				"appHash": "fedcba654321",
//				"telegramId": "123456789"
//			}
//		}
//	]
//
// # Example Usage:
//
//	accounts, err := LoadAccounts("accounts.json")
//	if err != nil {
//		log.Fatalf("Failed to load accounts: %v", err)
//	}
//	fmt.Println(accounts[0].GameData) // Output: user=%7B%22id%22%3A78894796...
func LoadAccounts(filePath string) ([]types.Account, error) {
	var accounts []types.Account
	file, err := os.Open(filePath)
	if err != nil {
		return accounts, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&accounts)
	return accounts, err
}

// LoadTasks reads the tasks.json file and parses it into a TaskCollection struct.
//
// # Parameters:
//   - filePath: The path to the tasks.json file.
//
// # Returns:
//   - types.TaskCollection: A struct containing the parsed task data.
//   - error: An error if the file cannot be opened, read, or parsed.
//
// # Example tasks.json:
//
//	{
//		"tasks": [
//			{
//				"name": "Task 1",
//				"payload": "Payload for Task 1"
//			},
//			{
//				"name": "Task 2",
//				"payload": "Payload for Task 2",
//				"interval": "Interval for Task 2"
//			}
//		]
//	}
//
// # Example Usage:
//
//	tasks, err := LoadTasks("tasks.json")
//	if err != nil {
//		log.Fatalf("Failed to load tasks: %v", err)
//	}
//	fmt.Println(tasks.Tasks[0].Name) // Output: Task 1
func LoadTasks(filePath string) (types.TaskCollection, error) {
	var tasks types.TaskCollection
	file, err := os.Open(filePath)
	if err != nil {
		return tasks, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	return tasks, err
}
