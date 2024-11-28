package types

// TaskConfig represents the configuration for a one-time task.
//
// # Fields:
//   - Name: The name of the task.
//   - Payload: A map containing task-specific payload data.
//
// # Example Usage:
//
//	taskConfig := TaskConfig{
//		Name: "Example Task",
//		Payload: map[string]interface{}{
//			"key1": "value1",
//			"key2": 2,
//		},
//	}
//	fmt.Println(taskConfig.Name) // Output: Example Task
type TaskConfig struct {
	Name    string                 `json:"name"`    // Name of the task
	Payload map[string]interface{} `json:"payload"` // Task-specific payload
}

// RecurrentTaskConfig represents the configuration for a recurrent task.
//
// # Fields:
//   - Name: The name of the task.
//   - Payload: A map containing task-specific payload data.
//   - IntervalMinutes: The interval in minutes between task executions.
//
// # Example Usage:
//
//	recurrentTaskConfig := RecurrentTaskConfig{
//		Name: "Recurrent Task",
//		Payload: map[string]interface{}{
//			"key1": "value1",
//			"key2": 2,
//		},
//		IntervalMinutes: 60,
//	}
//	fmt.Println(recurrentTaskConfig.Name) // Output: Recurrent Task
type RecurrentTaskConfig struct {
	Name            string                 `json:"name"`             // Name of the task
	Payload         map[string]interface{} `json:"payload"`          // Task-specific payload
	IntervalMinutes int                    `json:"interval_minutes"` // Interval in minutes between executions
}

// TaskCollection groups all tasks, both one-time and recurrent, for easier loading and management.
//
// # Fields:
//   - OneTimeTasks: A list of one-time tasks.
//   - RecurrentTasks: A list of recurrent tasks.
//
// # Example Usage:
//
//	taskCollection := TaskCollection{
//		OneTimeTasks: []TaskConfig{
//			{
//				Name: "One-Time Task",
//				Payload: map[string]interface{}{
//					"key1": "value1",
//				},
//			},
//		},
//		RecurrentTasks: []RecurrentTaskConfig{
//			{
//				Name: "Recurrent Task",
//				Payload: map[string]interface{}{
//					"key1": "value1",
//				},
//				IntervalMinutes: 30,
//			},
//		},
//	}
//	fmt.Println(taskCollection.OneTimeTasks[0].Name) // Output: One-Time Task
type TaskCollection struct {
	OneTimeTasks   []TaskConfig          `json:"one_time_tasks"`  // List of one-time tasks
	RecurrentTasks []RecurrentTaskConfig `json:"recurrent_tasks"` // List of recurrent tasks
}
