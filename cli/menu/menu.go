package menu

const inputErrMsg = "Неверный ввод"

type Menuer interface {
	ShowFirstLevel()
	ShowSecondLevel()
}

func New(isAdmin bool) Menuer {
	if isAdmin {
		return &AdminMenu{}
	}

	return &GuestMenu{}
}
