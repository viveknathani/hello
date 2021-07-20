package widgets

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// emails will hold two values, from and subject
type emailData struct {
	from    string
	subject string
}

// getEmailData will fetch all unread emails from the Primary section
// of your inbox.
func getEmailData() []emailData {

	// service is an instance of the Gmail API client.
	service := setup()

	// "me" is a specification mandated by the API which indicates we are dealing
	// for the currently authenticated user over here.
	user := "me"
	response, err := service.Users.Messages.List(user).LabelIds("UNREAD", "INBOX").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	messageMeta := response.Messages

	// Memory to store email data
	result := make([]emailData, 0)

	// Run through the metadata of the messages
	for _, meta := range messageMeta {

		// Use the id from metadata of each message to get the whole message
		subResponse, err := service.Users.Messages.Get(user, meta.Id).Do()

		// We are concerned with the headers stored in the payload of the subResponse
		headers := subResponse.Payload.Headers

		if err != nil {
			log.Fatalf("Unable: %v", err)
		}

		var from string = ""
		var subject string = ""

		// We are concenred with the top-level of header that contains
		// From and Subject keys.
		for _, header := range headers {

			if header.Name == "From" {
				from = header.Value
			}
			if header.Name == "Subject" {
				subject = header.Value
			}
		}

		result = append(result, emailData{from: from, subject: subject})
	}
	return result
}

// Below is a lot of boilerplate code that is required to just get
// a list of unread emails from the Gmail API.
// Check: [https://developers.google.com/gmail/api/quickstart/go]

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time. We save the file in the associated user's home directory.
	tokFile := "/home/" + os.Getenv("HOME_USER") + "/bin/hello/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func setup() *gmail.Service {

	ctx := context.TODO()
	// Get the credentials.json file of the associated user.
	b, err := ioutil.ReadFile("/home/" + os.Getenv("HOME_USER") + "/bin/hello/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	return service
}

// End of the boilerplate needed to do a trivial task with Gmail
