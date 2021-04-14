package menu

import (
	"fmt"
	"os"
	"strconv"

	"sly-beaver/storage"
)

type GuestMenu struct {
	*Menu
	userAction uint8
}

func (m *GuestMenu) ShowFirstLevel() error {
	for {
		fmt.Println("1. Просмотреть номенклатуры")
		fmt.Println("2. Распечатать номенклатуры")

		userAction, err := m.readInput()
		if err != nil {
			return fmt.Errorf("guest action input: %w", err)
		}

		if m.isExitCalled(userAction) {
			os.Exit(0)
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
	switch m.userAction {
	case 1:
	case 2:
		err := m.showReportMenu(s)
		if err != nil {
			return fmt.Errorf("third level guest menu: %w", err)
		}
	}
	return nil
}
