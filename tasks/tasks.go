package tasks

import (
	"github.com/nexus-telegram/NexusSDK/types"
)

// Task is the interface implemented by all tasks (one-time and recurrent).
//
// This interface defines the methods that must be implemented by any task type,
// including one-time and recurrent tasks. It ensures that all tasks can be
// configured, have a name, and can be executed for a given account.
//
// # Methods:
//   - Configure(payload map[string]interface{}) error: Configures the task's payload.
//   - GetName() string: Returns the name of the task.
//   - Run(account types.Account, client *httpclient.HTTPClient) error: Executes the task for a given account.
type Task interface {
	Run(account types.Account, handler Handler) error // Execute the task for a given account
}

// Handler is an interface that abstracts the GameHandler functionality.
type Handler interface {
	Post(url string, payload []byte) ([]byte, error)
	GetBaseURL() string
	GetAccounts() []types.Account
}

// BaseTask provides shared functionality for all tasks.
//
// This struct includes common fields and methods that are shared among different task types,
// such as the task's name and payload.
//
// # Fields:
//   - Name: The name of the task.
//   - Payload: A map containing the task's payload data.
type BaseTask struct {
	Name    string                 // Name of the task
	Payload map[string]interface{} // Payload for the task
}
