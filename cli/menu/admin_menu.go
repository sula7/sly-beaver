package menu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"sly-beaver/storage"
)

type AdminMenu struct {
	userAction uint8
	reader     *bufio.Reader
	delim      byte
}

func (m *AdminMenu) ShowFirstLevel() error {
	fmt.Println("1. Создать номенклатуру")
	fmt.Println("2. Удалить номенклатуру")
	fmt.Println("3. Распечатать номенклатуру")

	for {
		var userAction string

		userAction, err := m.readInput()
		if err != nil {
			fmt.Println("admin action input: %w", err)
			continue
		}

		if isExistCalled(userAction) {
			os.Exit(0)
		}

		ua, err := strconv.Atoi(strings.TrimSpace(userAction))
		if err != nil {
			fmt.Println(inputErrMsg)
			continue
		}

		if ua > 3 || ua < 1 {
			fmt.Println(inputErrMsg)
			continue
		}

		m.userAction = uint8(ua)
		break
	}

	return nil
}

func (m *AdminMenu) ShowSecondLevel(s storage.Storage) error {
	switch m.userAction {
	case 1:
		assert := storage.Assert{}
		var err error

		fmt.Println("\nВведите наименование (ё или ` для отмены):")
		assert.Name, err = m.readInput()
		if err != nil {
			return fmt.Errorf("create - scan name input: %w", err)
		}

		if isExistCalled(assert.Name) {
			return nil
		}

		fmt.Println("Введите количество (ё или ` для отмены):")
		for {
			amount, err := m.readInput()
			if err != nil {
				return fmt.Errorf("create - scan amount input: %w", err)
			}

			if isExistCalled(amount) {
				return nil
			}

			assert.Amount, err = strconv.ParseInt(amount, 10, 64)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}
			break
		}

		fmt.Println("Введите стоимость (ё или ` для отмены):")
		for {
			cost, err := m.readInput()
			if err != nil {
				return fmt.Errorf("create - scan cost input: %w", err)
			}

			if isExistCalled(cost) {
				return nil
			}

			assert.Cost, err = strconv.ParseInt(cost, 10, 64)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}
			break
		}

		fmt.Println("Введите срок годности ГГГГ-ММ-ДД (ё или ` для отмены):")
		for {
			assert.ValidTo, err = m.readInput()
			if err != nil {
				return fmt.Errorf("create - scan valid to input: %w", err)
			}

			if isExistCalled(assert.ValidTo) {
				return nil
			}

			_, err = time.Parse("2006-01-02", assert.ValidTo)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}
			break
		}

		err = s.CreateAssert(&assert)
		if err != nil {
			return fmt.Errorf("create - db exec: %w", err)
		}

		fmt.Println("Запись создана")
		fmt.Println()
	case 2:
		asserts, err := s.GetNotDeletedAsserts()
		if err != nil {
			return fmt.Errorf("получение списка номенклатур: %w", err)
		}

		fmt.Println("Выберите номенклатуру к удалению (ё или ` для отмены):")

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"№", "Дата создания", "Наименование", "Количество"})
		for i := 0; i < len(asserts); i++ {
			t.AppendRow([]interface{}{asserts[i].ID, asserts[i].CreatedAt, asserts[i].Name, asserts[i].Amount},
				table.RowConfig{})
		}
		t.Render()

		var assert storage.Assert

		for {
			var id int64

			rowID, err := m.readInput()
			if err != nil {
				return fmt.Errorf("remove - scan id input: %w", err)
			}

			if isExistCalled(rowID) {
				return nil
			}

			id, err = strconv.ParseInt(rowID, 10, 64)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}

			if id == 0 {
				return nil
			}

			assert.ID = id

			fmt.Println("Введите причину удаления (ё или ` для отмены):")
			assert.RemoveReason, err = m.readInput()
			if err != nil {
				return fmt.Errorf("remove - scan reson input: %w", err)
			}

			if isExistCalled(assert.RemoveReason) {
				return nil
			}

			break
		}

		err = s.AddRemoveReason(&assert)
		if err != nil {
			return fmt.Errorf("remove - db exec: %w", err)
		}

		fmt.Println("Запись удалена")
		fmt.Println()
	case 3:
		err := m.showThirdLevel(s)
		if err != nil {
			return fmt.Errorf("third level menu: %w", err)
		}
	}

	return nil
}

func (m *AdminMenu) showThirdLevel(s storage.Storage) error {
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

		if isExistCalled(ua) {
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

		csvContent := fmt.Sprint("№;Наименование;Количество;Цена;Срок годности\n")
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
	}

	return nil
}

func (m *AdminMenu) readInput() (string, error) {
	input, err := m.reader.ReadString(m.delim)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}
