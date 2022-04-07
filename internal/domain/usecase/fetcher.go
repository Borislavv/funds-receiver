package usecase

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	FileExtension = ".csv"
)

type FetcherUseCase struct {
	EtfDbUrl      string
	EtfDbFilesDir string
}

func NewFetcher(etfDbUrl string, etfDbFilesDir string) FetcherUseCase {
	return FetcherUseCase{
		EtfDbUrl:      etfDbUrl,
		EtfDbFilesDir: etfDbFilesDir,
	}
}

// Download - will return a filepath of stored file
func (fetcher FetcherUseCase) Download(fileId int) (string, error) {
	url := fmt.Sprintf("%v/%v", fetcher.EtfDbUrl, fileId)
	filepath := fmt.Sprintf("%v/%v", fetcher.EtfDbFilesDir, fileId)

	// creating a file-container
	file, err := os.Create(filepath + FileExtension)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := fetcher.downloadInFile(url, file); err != nil {
		return "", err
	}

	return filepath + FileExtension, nil
}

func (fetcher FetcherUseCase) downloadInFile(url string, file *os.File) error {
	// get the data
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// write the body to file
	bytes, err := io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	if bytes == 0 {
		return errors.New("unable to download file from remote server")
	}

	return nil
}
