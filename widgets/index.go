// Author: Vivek Nathani
// This is the file that exposes the different components of the
// widgets package from the SayHello function.

package widgets

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	formatIntro   = "Hello, %s!\nCollecting data for you...\n"
	formatWeather = "The temperature is %f degree celsius and the condition: %s."
	formatSpeed   = "Your internet speed is %f Mb/s."
	maxWidth      = 100
	colorReset    = "\033[0m"
	colorRed      = "\033[31m"
	colorGreen    = "\033[32m"
	colorYellow   = "\033[33m"
	colorBlue     = "\033[34m"
	colorPurple   = "\033[35m"
	colorCyan     = "\033[36m"
	colorWhite    = "\033[37m"
	colorPink     = "\033[38;5;13m"
)

// printAnything will literally print anything that goes into
// the first parameter as a string, length number of times.
// If the third parameter, newLine is set to true, there will be
// a newline printed at the end.
func printAnything(ch string, length int, newLine bool) {

	for i := 0; i < length; i++ {
		fmt.Print(ch)
	}

	if newLine {
		fmt.Println()
	}
}

// printColor will print any UNIX-compatible color.
// A list of colors is defined above in the const section.
func printColor(colorName string) {

	os := runtime.GOOS
	if os != "windows" {
		fmt.Print(colorName)
	}
}

// SayHello is the function that exposes all the components
// of the widgets package by printing them to the console
// in a neat way.
func SayHello() {

	// Display introduction
	printColor(colorCyan)
	fmt.Printf(formatIntro, os.Getenv("USER_NAME"))
	printColor(colorReset)

	var emails []emailData
	var todos todoData
	var weather weatherData
	var speed float64

	// Collect data in parallel
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		emails = getEmailData()
	}()
	go func() {
		defer wg.Done()
		todos = getTodoData()
	}()
	go func() {
		defer wg.Done()
		weather = getWeatherData()
	}()
	go func() {
		defer wg.Done()
		if os.Getenv("NO_TEST") == "yes" {
			speed = 40
		} else {
			speed = getConnectionData()
		}
	}()
	wg.Wait()

	// Display weather and internet speed
	printAnything("-", maxWidth, true)
	fmt.Print("|")
	printColor(colorYellow)
	n, _ := fmt.Printf(formatWeather,
		weather.temperature, strings.ToLower(weather.condition))
	printColor(colorReset)
	printAnything(" ", maxWidth-n-2, false)
	fmt.Println("|")
	fmt.Print("|")
	printColor(colorYellow)
	n2, _ := fmt.Printf(formatSpeed, speed)
	printColor(colorReset)
	printAnything(" ", maxWidth-n2-2, false)
	fmt.Println("|")
	printAnything("-", maxWidth, true)

	// Neatness
	fmt.Println()

	// Display the todo list
	printColor(colorCyan)
	fmt.Println("Here's your todo list.")
	printColor(colorReset)
	printAnything("-", maxWidth, true)
	for _, todo := range todos.tasks {

		fmt.Print("|")
		printColor(colorRed)
		n3, _ := fmt.Print(todo)
		printColor(colorReset)
		printAnything(" ", maxWidth-n3-2, false)
		fmt.Println("|")
	}
	printAnything("-", maxWidth, true)

	// Display unread emails if they exist
	if len(emails) != 0 {
		fmt.Println()
		printColor(colorCyan)
		fmt.Println("Here's a summary of your unread emails.")
		printColor(colorReset)
		for _, email := range emails {

			printAnything("-", maxWidth, true)
			fmt.Print("|")
			printColor(colorGreen)
			n3, _ := fmt.Printf("From: %s", email.from)
			printColor(colorReset)
			printAnything(" ", maxWidth-n3-2, false)
			fmt.Println("|")

			fmt.Print("|")
			printColor(colorGreen)
			n4, _ := fmt.Printf("Subject: %s", email.subject)
			printColor(colorReset)
			printAnything(" ", maxWidth-n4-2, false)
			fmt.Println("|")
			printAnything("-", maxWidth, true)
			fmt.Println()
		}
	}
}
