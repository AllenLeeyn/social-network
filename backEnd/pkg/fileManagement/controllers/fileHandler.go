package controller

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	errorControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/utils"
	"strings"
)

// Allowed file extensions
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// Check if file extension is allowed
func IsAllowedExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedExtensions[ext]
}

func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form with a max memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	// Retrieve all uploaded files
	files := r.MultipartForm.File["fileAttachments"]
	if len(files) == 0 {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	uploadedFiles := make(map[string]string)
	for _, handler := range files {
		file, err := handler.Open()
		if err != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}
		defer file.Close()

		// Call your file upload function
		uploadedFileName, err := FileUpload(file, handler)
		if err != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}

		uploadedFiles[handler.Filename] = uploadedFileName
	}

	utils.ReturnJsonSuccess(w, "Files uploaded successfully", uploadedFiles)
}

func FileUpload(file multipart.File, handler *multipart.FileHeader) (string, error) {
	// Check file extension before proceeding
	if !IsAllowedExtension(handler.Filename) {
		return "", fmt.Errorf("file type not allowed: %s", handler.Filename)
	}

	uploadDir := "./pkg/fileManagement/static/uploads"
	os.MkdirAll(uploadDir, os.ModePerm) // Ensure directory exists

	fileUUID, err := utils.GenerateUuid()
	if err != nil {
		return "", err
	}

	fileExt := filepath.Ext(handler.Filename)
	filePath := filepath.Join(uploadDir, fileUUID+fileExt)
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return "", err
	}
	defer outFile.Close()

	// Copy the uploaded file to the new location
	_, err = io.Copy(outFile, file)
	if err != nil {
		log.Println("Error saving file:", err)
		return "", err
	}

	return fileUUID + fileExt, nil
}
