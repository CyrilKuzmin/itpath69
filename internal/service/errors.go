package service

import "fmt"

func errModuleNotAllowed(id int) error {
	return fmt.Errorf("У вас нет доступа к модулю #%v", id)
}

func errPartNotExists(id int) error {
	return fmt.Errorf("Часть с номером %v не существует", id)
}
