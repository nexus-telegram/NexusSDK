package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/nexus-telegram/NexusSDK/types"
	"time"
)

// RecurrentTask represents a task that runs repeatedly at a set interval.
type RecurrentTask struct {
	BaseTask
	Interval time.Duration // Interval between executions
}

// NewRecurrentTask creates a new recurrent task.
func NewRecurrentTask(name string, payload map[string]interface{}, interval time.Duration) *RecurrentTask {
	return &RecurrentTask{
		BaseTask: BaseTask{Name: name, Payload: payload},
		Interval: interval,
	}
}

// Run executes the task for a given account.
func (task *RecurrentTask) Run(account types.Account, handler Handler) error {
	fmt.Printf("Running recurrent task '%s' for account %s with payload %v\n", task.Name, account.TelegramData.TelegramId, task.Payload)
	payload := task.Payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload for recurrent task '%s': %w", task.Name, err)
	}
	response, err := handler.Post(handler.GetBaseURL(), payloadBytes)
	if err != nil {
		return fmt.Errorf("failed to execute recurrent task '%s' for account %s: %w", task.Name, account.TelegramData.TelegramId, err)
	}
	fmt.Printf("Successfully executed recurrent task '%s' for account %s with response: %v\n", task.Name, account.TelegramData.TelegramId, response)
	return nil
}
