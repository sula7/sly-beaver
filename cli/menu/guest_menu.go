package menu

import (
	"fmt"
	"os"
	"strconv"

	"sly-beaver/storage"
)

type GuestMenu struct {
	userAction uint8
}

func (m *GuestMenu) ShowFirstLevel() error {
	for {
		var userAction string

		fmt.Println("0. Выход из приложения")
		fmt.Println("1. Просмотреть номенклатуры")
		fmt.Println("2. Распечатать номенклатуры")

		_, err := fmt.Scanln(&userAction)
		if err != nil {
			return fmt.Errorf("guest action input: %w", err)
		}

		ua, err := strconv.Atoi(userAction)
		if err != nil {
			fmt.Println(inputErrMsg)
			continue
		}

		if ua == 0 {
			os.Exit(0)
		}

		if ua > 2 || ua < 1 {
			fmt.Println(inputErrMsg)
			continue
		}

		m.userAction = uint8(ua)
		break
	}

	return nil
}

func (m *GuestMenu) ShowSecondLevel(s storage.Storage) error {
	return nil
}

func (m *GuestMenu) showThirdLevel(s storage.Storage) error {
	return nil
}