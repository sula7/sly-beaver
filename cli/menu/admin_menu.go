package menu

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"sly-beaver/storage"
)

type AdminMenu struct {
	userAction uint8
}

func (m *AdminMenu) ShowFirstLevel() {
	for {
		var userAction string
		fmt.Println("0. Выход из приложения")
		fmt.Println("1. Создать номенклатуру")
		fmt.Println("2. Удалить номенклатуру")
		fmt.Println("3. Распечатать номенклатуру")

		_, err := fmt.Scanln(&userAction)
		if err != nil {
			fmt.Println(inputErrMsg)
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
}

func (m *AdminMenu) ShowSecondLevel() {
	switch m.userAction {
	case 1:
		assert := storage.Assert{}
		var err error

		for {
			fmt.Println("\nВведите наименование:")
			_, err = fmt.Scanln(&assert.Name)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}
			assert.Name = strings.TrimSpace(assert.Name)
			break
		}

		for {
			var amount string
			fmt.Println("Введите количество:")
			_, err = fmt.Scanln(&amount)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
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
				fmt.Println(inputErrMsg)
				continue
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
				fmt.Println(inputErrMsg)
				continue
			}
			_, err = time.Parse("2006-01-02", assert.ValidTo)
			if err != nil {
				fmt.Println(inputErrMsg)
				continue
			}
			break
		}
	case 2:
	case 3:
	}
}
