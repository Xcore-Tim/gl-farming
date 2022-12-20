package googleDriveAPI

import (
	"context"
	"encoding/json"
	"fmt"
	googleDriveParams "gl-farming/app/constants/google-drive"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
)

// Use Service account
func serviceAccount(secretFile string) *http.Client {
	b, err := os.ReadFile(secretFile)
	if err != nil {
		log.Fatal("error while reading the credential file", err)
	}
	var s = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal(b, &s)
	config := &jwt.Config{
		Email:      s.Email,
		PrivateKey: []byte(s.PrivateKey),
		Scopes: []string{
			drive.DriveScope,
		},
		TokenURL: google.JWTTokenURL,
	}
	client := config.Client(context.Background())
	return client
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, folderId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{folderId},
	}

	file, err := service.Files.Create(f).Media(content).Fields("*").Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func DeleteFile(fileId string) error {

	client := serviceAccount(googleDriveParams.Credentials)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	if err = srv.Files.Delete(fileId).Do(); err != nil {
		return err
	}

	return nil
}

func Upload(fileName, filePath string) (string, string) {
	// Step 1: Open  file
	f, err := os.Open(filePath)

	if err != nil {
		panic(fmt.Sprintf("cannot open file: %v", err))
	}

	defer f.Close()

	// Step 2: Get the Google Drive service
	client := serviceAccount(googleDriveParams.Credentials)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	//give your folder id here in which you want to upload or create new directory
	folderId := googleDriveParams.FolderId

	// Step 4: create the file and upload
	file, err := createFile(srv, fileName, "application/octet-stream", f, folderId)

	if err != nil {
		panic(fmt.Sprintf("Could not create file: %v\n", err))
	}

	return file.Id, file.WebContentLink

}
