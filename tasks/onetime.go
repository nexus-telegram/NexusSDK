package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/nexus-telegram/NexusSDK/types"
)

// OneTimeTask represents a task that runs once per account.
type OneTimeTask struct {
	BaseTask
}

// NewOneTimeTask creates a new one-time task.
func NewOneTimeTask(name string, payload map[string]interface{}) *OneTimeTask {
	return &OneTimeTask{
		BaseTask: BaseTask{Name: name, Payload: payload},
	}
}

// Run executes the task for a given account.
func (task *OneTimeTask) Run(account types.Account, handler Handler) error {
	fmt.Printf("Running one-time task '%s' for account %s with payload %v\n", task.Name, account.TelegramData.TelegramId, task.Payload)
	payload := task.Payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload for one-time task '%s': %w", task.Name, err)
	}
	response, err := handler.Post(handler.GetBaseURL(), payloadBytes)
	if err != nil {
		return fmt.Errorf("failed to execute one-time task '%s' for account %s: %w", task.Name, account.TelegramData.TelegramId, err)
	}
	fmt.Printf("Successfully executed one-time task '%s' for account %s with response: %v\n", task.Name, account.TelegramData.TelegramId, response)
	return nil
}
