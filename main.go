package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func createFile(filename string) (*os.File, error) {
	// creates and uploads directory if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	// Build the file path and create it
	dst, err := os.Create(filepath.Join("uploads", filename))
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Limit file size to 10MB | 1 << 20 = 1MB
	r.ParseMultipartForm(10 << 20)

	// Retrieve the file from the form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}

	defer file.Close()

	fmt.Fprintf(w, "Uploaded File: %s\n", handler.Filename)
	fmt.Fprintf(w, "File Size: %d\n", handler.Size)
	fmt.Fprintf(w, "MIME Header: %v\n", handler.Header)

	// Save it locally
	dst, err := createFile(handler.Filename)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/upload", fileUploadHandler)

	fmt.Println("Server is running on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
