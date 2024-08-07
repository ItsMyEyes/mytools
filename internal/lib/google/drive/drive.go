package drive_google

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type AuthService struct {
}

func (s *AuthService) GetGoogleClient(ctx context.Context, credentialsPath string) (*http.Client, error) {
	// b, err := os.ReadFile(credentialsPath)
	// if err != nil {
	// 	log.Printf("Unable to read credentials file '%s': %+v", credentialsPath, err)
	// 	return nil, err
	// }

	hardCodedConfig := `{"web":{"client_id":"50840867934-e7kg8qrs5d3grdlfo9b6il1n57r2b57c.apps.googleusercontent.com","project_id":"blog-private","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"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","redirect_uris":["http://localhost:8080/"]}}`
	config, err := google.ConfigFromJSON([]byte(hardCodedConfig), drive.DriveFileScope, drive.DriveAppsReadonlyScope)
	if err != nil {
		log.Printf("Unable to parse credentials file to config: %+v", err)
		return nil, err
	}

	client, err := getClient(ctx, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		if tok, err = getTokenFromWeb(ctx, config); err != nil {
			return nil, err
		}
		if err = saveToken(tokFile, tok); err != nil {
			return nil, err
		}
	}
	return config.Client(ctx, tok), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("Unable to open file '%s': %+v", file, err)
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	fmt.Print("Enter the authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Printf("Unable to read authorization code: %+v", err)
		return nil, err
	}

	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		log.Printf("Unable to retrieve token from web: %+v", err)
		return nil, err
	}

	return tok, nil
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("Unable to cache oauth token: %+v", err)
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}
