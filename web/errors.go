package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func errInternal(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, "Системная ошибка: %v", err.Error())
}

func errLoginFailed() *echo.HTTPError {
	return echo.NewHTTPError(http.StatusUnauthorized, "Неверное имя пользователя и/или пароль")
}

func errUserAlreadyExists(username string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Пользователь %v уже существует", username))
}
