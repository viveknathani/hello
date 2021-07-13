package widgets

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// getConnectionData will return a floating-point number indicating the
// download speed, which is mostly what we are concerned with.
func getConnectionData() float64 {
	buffer, err := exec.Command("speedtest", "-f", "json").Output()

	if err != nil {
		log.Fatal(err)
	}

	// Parse the bytes as if they are JSON and store them in the variable payload
	var payload interface{}
	json.Unmarshal(buffer, &payload)

	// Payload is now "extractable".
	// We can transform and store the data as a map.
	data := payload.(map[string]interface{})
	data_download := data["download"].(map[string]interface{})
	bandwidth := data_download["bandwidth"].(float64)

	var result float64 = (bandwidth / float64(1024*1024)) * 8.0
	return result
}

func Run() {
	fmt.Println(getConnectionData())
}
