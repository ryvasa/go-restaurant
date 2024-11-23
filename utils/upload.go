package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UploadFile(file multipart.File, handler *multipart.FileHeader, folder string) (string, error) {
    // Validasi tipe file
    if err := validateFileType(handler.Filename); err != nil {
        return "", err
    }

    // Buat nama file unik
    filename := generateUniqueFilename(handler.Filename)

    // Buat folder jika belum ada
    uploadDir := fmt.Sprintf("uploads/%s", folder)
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        return "", err
    }

    // Buat file baru
    filepath := fmt.Sprintf("%s/%s", uploadDir, filename)
    dst, err := os.Create(filepath)
    if err != nil {
        return "", err
    }
    defer dst.Close()

    // Copy file yang diupload ke file baru
    if _, err = io.Copy(dst, file); err != nil {
        return "", err
    }

    // Return path relatif untuk disimpan di database
    return fmt.Sprintf("%s/%s", folder, filename), nil
}

func validateFileType(filename string) error {
    ext := strings.ToLower(filepath.Ext(filename))
    validTypes := map[string]bool{
        ".jpg":  true,
        ".jpeg": true,
        ".png":  true,
    }

    if !validTypes[ext] {
        return fmt.Errorf("tipe file tidak didukung: %s", ext)
    }
    return nil
}

func generateUniqueFilename(originalFilename string) string {
    ext := filepath.Ext(originalFilename)
    return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}
