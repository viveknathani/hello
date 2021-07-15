package widgets

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// weatherData is a data type that holds temperature and weather condition
// for the environment variable "WEATHER_API_KEY"
type weatherData struct {
	temperature float32
	condition   string
}

// getWeatherData will make an API call to weatherapi.com and
// return the weather information upon success
func getWeatherData() weatherData {

	const server = "http://api.weatherapi.com/v1"
	const path = "/current.json"
	var query string = "?key=" + os.Getenv("WEATHER_API_KEY") + "&q=" + os.Getenv("WEATHER_API_CITY")

	response, err := http.Get(server + path + query)

	if err != nil {
		log.Fatal(err)
	}

	// Close at the end
	defer response.Body.Close()

	// Extract data as bytes
	buffer, bufferError := io.ReadAll(response.Body)

	if bufferError != nil {
		log.Fatal(bufferError)
	}

	// Parse the bytes as if they are JSON and store them in the variable payload
	var payload interface{}
	json.Unmarshal(buffer, &payload)

	// Payload is now "extractable".
	// We can transform and store the data as a map.
	data := payload.(map[string]interface{})

	// This map is further "extractable" if it contains any key that has nested JSON data
	// We can extract the in the following manner because we know the structure of the
	// incoming data.
	dataCurrent := data["current"].(map[string]interface{})
	dataCondition := dataCurrent["condition"].(map[string]interface{})
	temperature := dataCurrent["temp_c"].(float64)

	return weatherData{
		temperature: float32(temperature),
		condition:   dataCondition["text"].(string),
	}
}
