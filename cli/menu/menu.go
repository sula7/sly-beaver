package menu

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

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
}

type Menu struct {
	delim  byte
	reader *bufio.Reader
}

func New(isAdmin bool) Menuer {
	switch {
	case isAdmin && runtime.GOOS == "windows":
		return &AdminMenu{Menu: &Menu{reader: bufio.NewReader(os.Stdin), delim: '\r' + '\n'}}
	case isAdmin && runtime.GOOS != "windows":
		return &AdminMenu{Menu: &Menu{reader: bufio.NewReader(os.Stdin), delim: '\n'}}
	case !isAdmin && runtime.GOOS == "windows":
		return &GuestMenu{Menu: &Menu{reader: bufio.NewReader(os.Stdin), delim: '\r' + '\n'}}
	case !isAdmin && runtime.GOOS != "windows:":
		return &GuestMenu{Menu: &Menu{reader: bufio.NewReader(os.Stdin), delim: '\n'}}
	}
	return nil
}

func (m *Menu) isExitCalled(userInput string) bool {
	if userInput == "`" || userInput == "ё" {
		return true
	}

	return false
}

func (m *Menu) readInput() (string, error) {
	input, err := m.reader.ReadString(m.delim)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}

func (m *Menu) showReportMenu(s storage.Storage) error {
	fmt.Println("Выберите тип отчёта для выгрузки (ё или ` для возврата):")
	fmt.Println("1. Отчёт по созданным за неделю номенклатурам")
	fmt.Println("2. Отчёт по удалённым за неделю номенклатурам")
	fmt.Println("3. Отчёт по текущим номерклатурам")

	var userAction int
	for {
		ua, err := m.readInput()
		if err != nil {
			fmt.Println("admin action input: %w", err)
			continue
		}

		if m.isExitCalled(ua) {
			return nil
		}

		userAction, err = strconv.Atoi(ua)
		if err != nil {
			fmt.Println(inputErrMsg)
			continue
		}

		if userAction > 3 || userAction < 1 {
			fmt.Println(inputErrMsg)
			continue
		}
		break
	}

	switch userAction {
	case 1:
		asserts, err := s.GetLastWeekAllAsserts()
		if err != nil {
			return fmt.Errorf("get all asserts last week: %w", err)
		}

		csvContent := fmt.Sprintln("№;Наименование;Количество;Цена;Срок годности")
		for i := 0; i < len(asserts); i++ {
			csvContent += fmt.Sprintf("%d;%s;%d;%d;%s\n",
				i+1, asserts[i].Name, asserts[i].Amount, asserts[i].Cost, asserts[i].ValidTo)
		}
		csvContent += fmt.Sprintf(";;;;;%s", time.Now().Format("2006-01-02"))

		err = createReportFile(csvContent, allAssertsFilename)
		if err != nil {
			return fmt.Errorf("all asserts: %w", err)
		}

		fmt.Println("Отчёт создан", allAssertsFilename)
		fmt.Println()
	case 2:
		asserts, err := s.GetLastWeekRemovedAsserts()
		if err != nil {
			return fmt.Errorf("get removed asserts last week: %w", err)
		}

		csvContent := fmt.Sprintln("№;Наименование;Количество;Причина")
		for i := 0; i < len(asserts); i++ {
			csvContent += fmt.Sprintf("%d;%s;%d;%s\n",
				i+1, asserts[i].Name, asserts[i].Amount, asserts[i].RemoveReason)
		}
		csvContent += fmt.Sprintf(";;;;%s", time.Now().Format("2006-01-02"))

		err = createReportFile(csvContent, removedAssertsFilename)
		if err != nil {
			return fmt.Errorf("removed asserts: %w", err)
		}

		fmt.Println("Отчёт создан", removedAssertsFilename)
		fmt.Println()
	case 3:
		asserts, err := s.GetCurrentAsserts()
		if err != nil {
			return fmt.Errorf("get current asserts: %w", err)
		}

		csvContent := fmt.Sprintln("№;Наименование;Количество;Цена;Срок годности")
		for i := 0; i < len(asserts); i++ {
			csvContent += fmt.Sprintf("%d;%s;%d;%d;%s\n",
				i+1, asserts[i].Name, asserts[i].Amount, asserts[i].Cost, asserts[i].ValidTo)
		}
		csvContent += fmt.Sprintf(";;;;;%s", time.Now().Format("2006-01-02"))

		err = createReportFile(csvContent, currentAssertsFilename)
		if err != nil {
			return fmt.Errorf("current asserts: %w", err)
		}

		fmt.Println("Отчёт создан", currentAssertsFilename)
		fmt.Println()
	}

	return nil
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
