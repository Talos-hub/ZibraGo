package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Talos-hub/ZibraGo/internal/apperrors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	pathtotoken = "token.json"
)

type GoogleCloudApi struct {
	service *drive.Service
}

// NewGoogleCloudApi returns GoogleCloudApi.
// Where token, it is  credentials token
func NewGoogleCloudApi(token []byte) (*GoogleCloudApi, error) {
	ctx := context.Background()

	config, err := google.ConfigFromJSON(token, drive.DriveFileScope)
	if err != nil {
		return nil, apperrors.NewAppError("error create google cloud api, unable to parse client secret file to config", "NewGoogleCloudApi", apperrors.E_PARSE, err)
	}

	// show a token to user
	config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"

	client, err := getClient(config)
	if err != nil {
		return nil, apperrors.NewAppError("error create google cloud api, cannot get a client", "NewGoogleCloudApi", apperrors.E_CREATE, err)
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, apperrors.NewAppError("error create google cloud api,cannot create drive service", "NewGoogleCloudApi", apperrors.E_FATAL, err)
	}

	return &GoogleCloudApi{
		service: srv,
	}, nil
}

// Auth is authorization function, for chanching a google drive
func (a *GoogleCloudApi) Auth(token []byte) error {
	config, err := google.ConfigFromJSON(token, drive.DriveFileScope)
	if err != nil {
		return apperrors.NewAppError("error create google cloud api, unable to parse client secret file to config", "NewGoogleCloudApi", apperrors.E_PARSE, err)
	}

	// show a token to user
	config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"

	tok, err := getTokenFromWeb(config)
	if err != nil {
		return apperrors.NewAppError("error get token from web", "Auth", apperrors.E_PARSE, err)
	}

	err = saveToken(pathtotoken, tok)
	if err != nil {
		return apperrors.NewAppError("error save a token", "Auth", apperrors.E_CREATE, err)
	}
	return nil
}

// Check check connect to google drive
func (a *GoogleCloudApi) Check() error {
	_, err := a.service.About.Get().Fields("user").Do()
	if err != nil {
		return apperrors.NewAppError("cannot connect to drive", "Check", apperrors.E_API, err)
	}
	return nil
}

// UploadFiles create new file in gogole drive
func (a *GoogleCloudApi) UpLoadFiles(file *os.File) error {
	_, err := file.Seek(0, 0) // Reset file pointer
	if err != nil {
		return apperrors.NewAppError("cannot seek file", "UploadFile", apperrors.E_READ, err)
	}

	filedrive := &drive.File{
		Name: filepath.Base(file.Name()),
	}

	_, err = a.service.Files.Create(filedrive).Media(file).Do()
	if err != nil {
		return apperrors.NewAppError("cannot up load a file", "UpLoadFiles", apperrors.E_API, err)
	}
	return nil
}

// getCien retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := getTokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, fmt.Errorf("error get client: %w", err)
		}
		err := saveToken(tokFile, tok)
		if err != nil {
			return nil, fmt.Errorf("error get client: %w", err)
		}
	}
	return config.Client(context.Background(), tok), nil
}

// getTokenFromWeb is authentication function that return a token for personal drive
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web %v", err)
	}
	return tok, nil
}

// getTokenFromFile reutrns oauth2 token from file
func getTokenFromFile(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %w", err)
		}
		return nil, fmt.Errorf("error open file: %w", err)
	}

	token := &oauth2.Token{}

	err = json.NewDecoder(file).Decode(token)
	if err != nil {
		return nil, fmt.Errorf("error decoding file: %w", err)
	}
	return token, nil
}

// saveToken save a oauth2 token to file
func saveToken(path string, token *oauth2.Token) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error create or open file: %w", err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(token)
	if err != nil {
		return fmt.Errorf("error encoding: %w", err)
	}
	return nil
}
