package handler

import (
	"github.com/nexus-telegram/NexusSDK/httpclient"
	"github.com/nexus-telegram/NexusSDK/tasks"
	"github.com/nexus-telegram/NexusSDK/types"
	"io"
	"log"
	"sync"
	"time"
)

// GameHandler manages tasks, accounts, and API configuration.
//
// This struct holds the configuration and state necessary for managing tasks and accounts
// for a specific game. It includes the base API URL, proxy settings, API key, accounts,
// tasks, and an HTTP client for sending requests.
//
// # Fields:
//   - BaseURL: The base API URL for the specific game.
//   - Proxy: The proxy configuration for all requests.
//   - APIKey: The API key used for authentication.
//   - Accounts: A list of accounts to process.
//   - Tasks: A list of tasks, both one-time and recurrent.
//   - HttpClient: The HTTP client used for sending requests.
//   - mu: A mutex for thread-safe operations.
type GameHandler struct {
	GameName   string                 // Name of the game
	BaseURL    string                 // Base API URL for the specific game
	Proxy      types.Proxy            // Proxy configuration for all requests
	APIKey     string                 // API key for authentication
	Accounts   []types.Account        // List of accounts to process
	Tasks      []tasks.Task           // List of tasks (both one-time and recurrent)
	HttpClient *httpclient.HTTPClient // HTTP client for sending requests
	mu         sync.Mutex             // Mutex for thread-safe operations
}

// Post sends a POST request using the HTTP client.
func (handler *GameHandler) Post(url string, payload []byte) ([]byte, error) {
	resp, err := handler.HttpClient.Post(url, payload)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetBaseURL returns the base URL.
func (handler *GameHandler) GetBaseURL() string {
	return handler.BaseURL
}

// GetAccounts returns the list of accounts.
func (handler *GameHandler) GetAccounts() []types.Account {
	return handler.Accounts
}

// SetBaseURL sets the base API URL for the handler.
//
// This method updates the BaseURL field of the GameHandler with the provided URL.
//
// # Parameters:
//   - url: The new base API URL to be set.
//
// # Example:
//
//	handler.SetBaseURL("https://api.example.com")
func (handler *GameHandler) SetBaseURL(url string) {
	handler.BaseURL = url
}

// AddTask adds a new task to the handler.
//
// This method locks the handler's mutex to ensure thread-safe access to the tasks slice,
// then appends the new task to the slice. The mutex is unlocked after the task is added.
//
// # Parameters:
//   - task: The task to be added to the handler's list of tasks.
//
// # Example:
//
//	task := &tasks.OneTimeTask{}
//	handler.AddTask(task)
func (handler *GameHandler) AddTask(task tasks.Task) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.Tasks = append(handler.Tasks, task)
}

// RunTasks executes all tasks for all accounts.
//
// This method iterates over all accounts and tasks, executing each task in a separate goroutine.
// It uses a WaitGroup to ensure all tasks are completed before returning. One-time tasks are
// executed immediately, while recurrent tasks are executed at regular intervals specified by
// the task's interval.
//
// # Notes:
//   - One-time tasks are executed once per account.
//   - Recurrent tasks executes at regular intervals until the program stops.
//   - Errors during task execution do not stop the execution of other tasks.
//
// # Example:
//
//	handler.RunTasks()
func (handler *GameHandler) RunTasks() {
	var wg sync.WaitGroup
	for _, account := range handler.Accounts {
		wg.Add(1)
		go func(account types.Account) {
			defer wg.Done()
			for _, task := range handler.Tasks {
				switch t := task.(type) {
				case *tasks.OneTimeTask:
					if err := handler.runTaskWithRetry(account, t); err != nil {
						log.Printf("Error executing one-time task for account %s: %v\n", account.TelegramData.TelegramId, err)
					}
				case *tasks.RecurrentTask:
					func(task *tasks.RecurrentTask) {
						ticker := time.NewTicker(task.Interval)
						defer ticker.Stop()
						for {
							select {
							case <-ticker.C:
								if err := handler.runTaskWithRetry(account, task); err != nil {
									log.Printf("Error executing recurrent task for account %s: %v\n", account.TelegramData.TelegramId, err)
								}
							}
						}
					}(t)
				}
			}
		}(account)
	}
	wg.Wait()
}

// runTaskWithRetry attempts to run a task for a given account, retrying once if it fails.
//
// This method first tries to execute the task. If the task fails, it attempts to refresh
// the game data and retries the task once. If the retry also fails, the error is ignored.
//
// # Parameters:
//   - account: The account for which the task is being executed.
//   - task: The task to be executed.
//
// # Returns:
//   - error: An error if the task fails on both the initial attempt and the retry.
//
// # Notes:
//   - Errors during the initial task execution trigger a refresh of the game data.
//   - If the refresh fails, the method returns without retrying the task.
//   - If the retry fails, the method returns without further action.
func (handler *GameHandler) runTaskWithRetry(account types.Account, task tasks.Task) error {
	if err := task.Run(account, handler); err != nil {
		refreshErr, _ := refreshGameData(
			handler.HttpClient,
			handler.GameName,
			handler.APIKey,
			account.TelegramData,
			handler.Proxy,
		)
		if refreshErr != nil {
			return nil
		}
		if retryErr := task.Run(account, handler); retryErr != nil {
			return nil
		}
	}
	return nil
}

// NewGameHandler creates a new instance of GameHandler by loading the necessary configuration
// and game data from the specified file paths.
//
// It reads the configuration and game data, initializes a GameHandler instance with the
// loaded values, and prepares an HTTP client configured with the provided proxy settings.
//
// # Example:
//
//	handler, err := NewGameHandler("config.json", "game_data.json")
//	if err != nil {
//		log.Fatalf("Failed to create GameHandler: %v", err)
//	}
//	fmt.Println(handler)
//
// # Parameters:
//   - configPath: The file path to the configuration file (e.g., "config.json").
//   - gameDataPath: The file path to the game data file (e.g., "game_data.json").
//
// # Returns:
//   - *types.GameHandler: A pointer to the initialized GameHandler instance.
//   - error: An error if the configuration or game data cannot be loaded.
//
// # Notes:
//   - Ensure the provided file paths are valid and accessible.
//   - The `BaseURL` field of the GameHandler is left empty and should be set manually
//     before making API requests.
func NewGameHandler(configPath, gameDataPath string) (*GameHandler, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	accounts, err := LoadAccounts(gameDataPath)
	if err != nil {
		return nil, err
	}
	httpClient, err := httpclient.NewHTTPClient(config.Proxy)
	if err != nil {
		return nil, err
	}
	handler := &GameHandler{
		BaseURL:    "",
		Proxy:      config.Proxy,
		APIKey:     config.APIKey,
		Accounts:   accounts,
		HttpClient: httpClient,
	}
	return handler, nil
}
