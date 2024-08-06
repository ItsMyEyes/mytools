package main

import (
	"context"
	"fmt"
	"log"

	drive_google "github.com/ItsMyEyes/my-tools/internal/lib/google/drive"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Initialize the authentication service
	authSrv := drive_google.AuthService{}
	client, err := authSrv.GetGoogleClient(ctx, "./credentials/credentials.json")
	if err != nil {
		panic(err)
	}

	// Create Google Drive Service
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %+v", err)
	}

	r, err := srv.Files.List().Q(fmt.Sprintf("name = '%s'", "2024-08-06@16#test.txt")).Do()
	if err != nil {
		log.Fatalf("Unable to list replies: %+v", err)
	}
	for _, reply := range r.Files {
		fmt.Println(reply.Name)
		// err := srv.Files.Delete(reply.Id).Do()
		// if err != nil {
		// 	log.Fatalf("Unable to delete file: %+v", err)
		// }
	}

	// Upload file to Google Drive
	// filePath := "./test.txt" // replace with your file path
	// f, err := os.Open(filePath)
	// if err != nil {
	// 	log.Fatalf("Unable to open file: %+v", err)
	// }
	// defer f.Close()

	// file := &drive.File{
	// 	Name:    fmt.Sprintf("%s#test.txt", time.Now().Format("2006-01-02@15")), // replace with the name you want to give the file
	// 	Parents: []string{"1X3dhwicYRjmHV0_wtnVJi4vD-bCNfPz1"},
	// }

	// uploadedFile, err := srv.Files.Create(file).Media(f).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to upload file: %+v", err)
	// }

	// fmt.Printf("File '%s' uploaded successfully with ID: %s\n", uploadedFile.Name, uploadedFile.Id)
}
