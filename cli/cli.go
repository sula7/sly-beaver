package cli

import (
	"database/sql"
	"fmt"

	"sly-beaver/cli/menu"
	"sly-beaver/storage"
)

type CLI struct {
	storage storage.Storage
}

const inputErrMsg = "Неверный ввод, проверьте данные и повторите снова\n"

func NewCLIService(storage storage.Storage) *CLI {
	return &CLI{storage: storage}
}

func (c *CLI) Start() {
	isUserAdmin := c.authUser()

	m := menu.New(isUserAdmin)

	for {
		fmt.Println("Выберите действие:")

		m.ShowFirstLevel()

		m.ShowSecondLevel()
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
			fmt.Println(inputErrMsg)
		}

		fmt.Println("Введите пароль:")
		_, err = fmt.Scanln(&password)
		if err != nil {
			fmt.Println(inputErrMsg)
			continue
		}

		isUserExists, isAdmin, err = c.storage.CheckPassword(login, password)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println("check user in DB:", err)
			continue
		}

		if !isUserExists {
			fmt.Println("Неверный логин и/или пароль\n")
			continue
		}

		fmt.Println()
		break
	}

	return isAdmin
}
