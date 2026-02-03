package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileManager struct {
	uploadDir string
}

func NewFileManager(uploadDir string) *FileManager {
	os.MkdirAll(uploadDir, 0755)
	return &FileManager{
		uploadDir: uploadDir,
	}
}

func (fm *FileManager) SaveUploadedFile(part *multipart.Part, filename string) (string, error) {
	filePath := filepath.Join(fm.uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, part); err != nil {
		_ = os.Remove(filePath)
		return "", err
	}

	return filePath, nil
}

func (fm *FileManager) DeleteFile(filepath string) error {
	return os.Remove(filepath)
}

func HashFile(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
