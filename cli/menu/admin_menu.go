package menu

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"sly-beaver/storage"
)

type AdminMenu struct {
	*Menu
	userAction uint8
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

		if m.isExitCalled(userAction) {
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

		if m.isExitCalled(assert.Name) {
			return nil
		}

		fmt.Println("Введите количество (ё или ` для отмены):")
		for {
			amount, err := m.readInput()
			if err != nil {
				return fmt.Errorf("create - scan amount input: %w", err)
			}

			if m.isExitCalled(amount) {
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

			if m.isExitCalled(cost) {
				return nil
			}

			assert.Cost, err = strconv.ParseInt(cost, 10, 64)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}
			break
		}

		fmt.Println("Введите дату покупки ГГГГ-ММ-ДД (ё или ` для отмены):")
		for {
			assert.BuyDate, err = m.readInput()
			if err != nil {
				return fmt.Errorf("create - scan valid to input: %w", err)
			}

			if m.isExitCalled(assert.BuyDate) {
				return nil
			}

			_, err = time.Parse("2006-01-02", assert.BuyDate)
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

		rows := []table.Row{}
		for i := 0; i < len(asserts); i++ {
			rows = append(rows, table.Row{asserts[i].ID, asserts[i].CreatedAt, asserts[i].Name, asserts[i].Amount})
		}

		m.renderAssertsView(table.Row{"№", "Дата создания", "Наименование", "Количество"}, rows)

		var assert storage.Assert
		for {
			var id int64

			rowID, err := m.readInput()
			if err != nil {
				return fmt.Errorf("remove - scan id input: %w", err)
			}

			if m.isExitCalled(rowID) {
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

			if m.isExitCalled(assert.RemoveReason) {
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
		err := m.showReportMenu(s)
		if err != nil {
			return fmt.Errorf("third level admin menu: %w", err)
		}
	}

	return nil
}
