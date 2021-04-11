package menu

import "sly-beaver/storage"

const inputErrMsg = "Неверный ввод"

type Menuer interface {
	ShowFirstLevel()
	ShowSecondLevel(s storage.Storage) error
}

func New(isAdmin bool) Menuer {
	if isAdmin {
		return &AdminMenu{}
	}

	return &GuestMenu{}
}
