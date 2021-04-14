package menu

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"sly-beaver/storage"
)

const (
	inputErrMsg            = "Неверный ввод"
	allAssertsFilename     = "all_asserts.csv"
	removedAssertsFilename = "removed_asserts.csv"
	currentAssertsFilename = "current_asserts.csv"
)

type Menuer interface {
	ShowFirstLevel() error
	ShowSecondLevel(s storage.Storage) error
	showThirdLevel(s storage.Storage) error
}

func New(isAdmin bool) Menuer {
	switch {
	case isAdmin && runtime.GOOS == "windows":
		return &AdminMenu{reader: bufio.NewReader(os.Stdin), delim: '\r' + '\n'}
	case isAdmin && runtime.GOOS != "windows":
		return &AdminMenu{reader: bufio.NewReader(os.Stdin), delim: '\n'}
	case !isAdmin && runtime.GOOS == "windows":
		return &GuestMenu{reader: bufio.NewReader(os.Stdin), delim: '\r' + '\n'}
	case !isAdmin && runtime.GOOS != "windows:":
		return &GuestMenu{reader: bufio.NewReader(os.Stdin), delim: '\n'}
	}
	return nil
}

func isExistCalled(userInput string) bool {
	if userInput == "`" || userInput == "ё" {
		return true
	}

	return false
}

func createReportFile(fileContent, filename string) error {
	workDir := filepath.Dir(os.Args[0])
	csvPath := filepath.Join(workDir, filename)
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("create csv file: %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("defer csv file close:", err)
		}
	}()

	_, err = file.WriteString(fileContent)
	if err != nil {
		return fmt.Errorf("write into csv file: %w", err)
	}

	return err
}
