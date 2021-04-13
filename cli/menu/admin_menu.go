package menu

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"sly-beaver/storage"
)

type AdminMenu struct {
	userAction uint8
}

func (m *AdminMenu) ShowFirstLevel() error {
	fmt.Println("1. Создать номенклатуру")
	fmt.Println("2. Удалить номенклатуру")
	fmt.Println("3. Распечатать номенклатуру")

	for {
		var userAction string

		_, err := fmt.Scanln(&userAction)
		if err != nil {
			fmt.Println("admin action input: %w", err)
			continue
		}

		if isExistCalled(userAction) {
			os.Exit(0)
		}

		ua, err := strconv.Atoi(userAction)
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
		_, err = fmt.Scanln(&assert.Name)
		if err != nil {
			return fmt.Errorf("create - scan name input: %w", err)
		}
		if isExistCalled(assert.Name) {
			return nil
		}
		assert.Name = strings.TrimSpace(assert.Name)

		fmt.Println("Введите количество (ё или ` для отмены):")
		for {
			var amount string
			_, err = fmt.Scanln(&amount)
			if err != nil {
				return fmt.Errorf("create - scan amount input: %w", err)
			}
			if isExistCalled(amount) {
				return nil
			}

			assert.Amount, err = strconv.ParseInt(strings.TrimSpace(amount), 10, 64)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}
			break
		}

		fmt.Println("Введите стоимость (ё или ` для отмены):")
		for {
			var cost string
			_, err = fmt.Scanln(&cost)
			if err != nil {
				return fmt.Errorf("create - scan cost input: %w", err)
			}
			if isExistCalled(cost) {
				return nil
			}

			assert.Cost, err = strconv.ParseInt(strings.TrimSpace(cost), 10, 64)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}
			break
		}

		fmt.Println("Введите срок годности ГГГГ-ММ-ДД (ё или ` для отмены):")
		for {
			_, err = fmt.Scanln(&assert.ValidTo)
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
			var rowID string

			_, err = fmt.Scanln(&rowID)
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
			_, err = fmt.Scanln(&assert.RemoveReason)
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
		var ua string

		_, err := fmt.Scanln(&ua)
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
		asserts, err := s.GetAllRowsForCSV()
		if err != nil {
			return fmt.Errorf("get all asserts last week: %w", err)
		}

		csvContent := fmt.Sprint("№;Наименование;Количество;Цена;Срок годности\n")
		for i := 0; i < len(asserts); i++ {
			csvContent += fmt.Sprintf("%d;%s;%d;%d;%s\n",
				i+1, asserts[i].Name, asserts[i].Amount, asserts[i].Cost, asserts[i].ValidTo)
		}
		csvContent += fmt.Sprintf(";;;;;%s", time.Now().Format("2006-01-02"))

		workDir := filepath.Dir(os.Args[0])
		csvPath := filepath.Join(workDir, "all_rows.csv")
		file, err := os.Create(csvPath)
		if err != nil {
			return fmt.Errorf("create csv for all rows csv: %w", err)
		}

		defer func() {
			err := file.Close()
			if err != nil {
				log.Println("defer all rows csv file close:", err)
			}
		}()

		_, err = file.WriteString(csvContent)
		if err != nil {
			return fmt.Errorf("all rows csv file write: %w", err)
		}

		fmt.Println("Отчёт создан", csvPath)
		fmt.Println()
	case 2:
	case 3:
	}

	return nil
}
