package menu

import "sly-beaver/storage"

const inputErrMsg = "Неверный ввод"

type Menuer interface {
	ShowFirstLevel() error
	ShowSecondLevel(s storage.Storage) error
}

func New(isAdmin bool) Menuer {
	if isAdmin {
		return &AdminMenu{}
	}

	return &GuestMenu{}
}

func isExistCalled(userInput string) bool {
	if userInput == "`" || userInput == "ё" {
		return true
	}

	return false
}
