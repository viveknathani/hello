// Author: Vivek Nathani
// hello is a command line program which is intended to serve as a
// personal dashboard for myself. This program is not made to be used
// just out of the box for anybody. However, every portion of the source-code
// is customizable, if you know how to write code in Go. The task of this
// program is simple. It makes API calls to different services the author is
// interested in knowing about (local weather, todolist, emails, internet speed).
// If you have any improvements or bug fixes in mind, feel free to create an issue
// on Github and link it with a pull request. Or, you could write an email to me at
// <viveknathani2402@gmail.com>

package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/viveknathani/hello/widgets"
)

// init is tasked with loading up the environment variables
// It should be noted that .env file is not the only environment
// setting that is needed. GMail API requires two additional files
// to be in place which are credentials.json and token.json. The latter
// is created after you log-in from your OAuth consent screen.
func init() {
	err := godotenv.Load(os.ExpandEnv("/home/$USER/hello/.env"))
	if err != nil {
		log.Fatal("Could not load .env file!")
	}
}

func main() {
	test := flag.Bool("test", true, "Do a fresh speedtest.")
	flag.Parse()
	if !(*test) {
		os.Setenv("NO_TEST", "yes")
	}
	widgets.SayHello()
}
