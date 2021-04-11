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
	userAction uint8
}

func (m *AdminMenu) ShowFirstLevel() error {
	for {
		var userAction string
		fmt.Println("0. Выход из приложения")
		fmt.Println("1. Создать номенклатуру")
		fmt.Println("2. Удалить номенклатуру")
		fmt.Println("3. Распечатать номенклатуру")

		_, err := fmt.Scanln(&userAction)
		if err != nil {
			fmt.Println("admin action input: %w", err)
			continue
		}

		ua, err := strconv.Atoi(userAction)
		if err != nil {
			fmt.Println(inputErrMsg)
			continue
		}

		if ua == 0 {
			os.Exit(0)
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

		fmt.Println("\nВведите наименование:")
		_, err = fmt.Scanln(&assert.Name)
		if err != nil {
			return fmt.Errorf("create - scan name input: %w", err)
		}
		assert.Name = strings.TrimSpace(assert.Name)

		for {
			var amount string
			fmt.Println("Введите количество:")
			_, err = fmt.Scanln(&amount)
			if err != nil {
				return fmt.Errorf("create - scan amount input: %w", err)
			}

			assert.Amount, err = strconv.ParseInt(strings.TrimSpace(amount), 10, 64)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}
			break
		}

		for {
			var cost string
			fmt.Println("Введите стоимость:")
			_, err = fmt.Scanln(&cost)
			if err != nil {
				return fmt.Errorf("create - scan cost input: %w", err)
			}

			assert.Cost, err = strconv.ParseInt(strings.TrimSpace(cost), 10, 64)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}
			break
		}

		for {
			fmt.Println("Введите срок годности (ГГГГ-ММ-ДД):")
			_, err = fmt.Scanln(&assert.ValidTo)
			if err != nil {
				return fmt.Errorf("create - scan valid to input: %w", err)
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

		fmt.Println("Выберите номенклатуру к удалению (0 для отмены):")

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

			id, err = strconv.ParseInt(rowID, 10, 64)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}

			if id == 0 {
				return nil
			}

			assert.ID = id

			fmt.Println("Введите причину удаления:")
			_, err = fmt.Scanln(&assert.RemoveReason)
			if err != nil {
				return fmt.Errorf("remove - scan reson input: %w", err)
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
	}

	return nil
}
