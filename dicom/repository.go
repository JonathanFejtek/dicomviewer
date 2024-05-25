package dicom

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

const defaultFileDest = "/tmp/dicom"

var (
	// ErrFileNotFound error indicating specified file was not found
	ErrFileNotFound = errors.New("file was not found")
)

// FileRepository represents a persistent store for DICOM files
type FileRepository interface {
	GetAll() ([]string, error)
	Get(id string) (*File, error)
	Create(d File) error
}

// localDICOMFileAdapter an implementation of FileRepository that uses local
// temp file storage to store DICOM files
type localDICOMFileAdapter struct{}

// NewLocalFileAdapter construct a local file adapter repository
func NewLocalFileAdapter() FileRepository {
	return &localDICOMFileAdapter{}
}

func (d *localDICOMFileAdapter) GetAll() ([]string, error) {
	entries, err := os.ReadDir(defaultFileDest)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}

	var fileNames []string
	for _, e := range entries {
		fileNames = append(fileNames, e.Name())
	}

	return fileNames, nil
}

// Create create a new DICOM file
func (d *localDICOMFileAdapter) Create(file File) error {
	filename := d.generateFileName(file.ID)

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, file.Raw()); err != nil {
		return err
	}

	if _, err := os.Stat(defaultFileDest); os.IsNotExist(err) {
		os.MkdirAll(defaultFileDest, 0700)
	}

	return os.WriteFile(filename, buffer.Bytes(), 0644)
}

// Get retrieve a DICOM file by id
func (d *localDICOMFileAdapter) Get(id string) (*File, error) {
	filename := d.generateFileName(id)
	file, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	defer file.Close()

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, file); err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	dicomFile := NewFile(
		id,
		fileInfo.Size(),
		// rather than return file, return a bytes reader so that we can
		// safely manage file closing in this scope
		bytes.NewReader(buffer.Bytes()),
	)

	return &dicomFile, nil
}

func (d localDICOMFileAdapter) generateFileName(id string) string {
	return fmt.Sprintf(defaultFileDest+"/%s", id)
}
