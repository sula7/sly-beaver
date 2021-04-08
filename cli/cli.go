package cli

import (
	"database/sql"
	"fmt"

	"sly-beaver/storage"
)

type CLI struct {
	storage storage.Storage
}

func NewCLIService(storage storage.Storage) *CLI {
	return &CLI{storage: storage}
}

func (c *CLI) Start() {
	isUserAdmin := c.authUser()

	for {
		showMainMenu(isUserAdmin)
		var mainMenuInput uint8
		_, err := fmt.Scanln(&mainMenuInput)
		if err != nil {
			fmt.Println("Неверный ввод, проверьте данные и повторите снова\n")
		}
	}

}

func (c *CLI) authUser() bool {
	var login string
	var password string
	var isAdmin bool
	var isUserExists bool

	for {
		fmt.Println("Введите логин:")
		_, err := fmt.Scanln(&login)
		if err != nil {
			fmt.Println("Неверный ввод, проверьте данные и повторите снова\n")
		}

		fmt.Println("Введите пароль:")
		_, err = fmt.Scanln(&password)
		if err != nil {
			fmt.Println("Неверный ввод, проверьте данные и повторите снова\n")
			continue
		}

		isUserExists, isAdmin, err = c.storage.CheckPassword(login, password)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println("check user in DB:", err)
			continue
		}

		if !isUserExists {
			fmt.Println("Неверные логин/пароль, проверьте данные и повторите снова\n")
			continue
		}

		fmt.Println()
		break
	}

	return isAdmin
}

func showMainMenu(isUserAdmin bool) {
	fmt.Println("Выберите действие:")

	switch isUserAdmin {
	case true:
		fmt.Println("1. Создать номенклатуру")
		fmt.Println("2. Удалить номенклатуру")
		fmt.Println("3. Распечатать номенклатуру")
	default:
		fmt.Println("1. Просмотреть номенклатуры")
		fmt.Println("2. Распечатать номенклатуры")
	}
}
