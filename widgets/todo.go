package widgets

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// todoData will hold all tasks coming from the Todoist API
type todoData struct {
	tasks []string
}

// getTodoData will make an API call to Todoist and return a struct
// which will contain a slice of tasks
func getTodoData() todoData {

	const server = "https://api.todoist.com/rest/v1"
	const path = "/tasks"
	request, err := http.NewRequest("GET", server+path, nil)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Authorization", "Bearer "+os.Getenv("TODO_API_KEY"))
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	// Close at the end
	defer response.Body.Close()
	buffer, bufferError := io.ReadAll(response.Body)

	if bufferError != nil {
		log.Fatal(bufferError)
	}

	// Parse the bytes as if they are JSON and store them in the variable payload
	var payload interface{}
	json.Unmarshal(buffer, &payload)

	// Payload is now "extractable".
	// We can store the data as an array now.
	data := payload.([]interface{})

	tasks := make([]string, 0)
	for i := 0; i < len(data); i++ {

		// Response data is an array of object, each containing the key, "content".
		ob := data[i].(map[string]interface{})
		content := ob["content"].(string)
		tasks = append(tasks, content)
	}

	return todoData{
		tasks: tasks,
	}
}
