package menu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"sly-beaver/storage"
)

type GuestMenu struct {
	userAction uint8
	reader     *bufio.Reader
	delim      byte
}

func (m *GuestMenu) ShowFirstLevel() error {
	for {
		fmt.Println("0. Выход из приложения")
		fmt.Println("1. Просмотреть номенклатуры")
		fmt.Println("2. Распечатать номенклатуры")

		userAction, err := m.readInput()
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

func (m *GuestMenu) readInput() (string, error) {
	input, err := m.reader.ReadString(m.delim)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}
