package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func errInternal(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Системная ошибка: %v", err.Error()))
}

func errBadRequest() *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, "Ошибка в запросе")
}

func errLoginFailed() *echo.HTTPError {
	return echo.NewHTTPError(http.StatusUnauthorized, "Неверное имя пользователя и/или пароль")
}

func errUserAlreadyExists(username string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Пользователь %v уже существует", username))
}

func errModuleNotAllowed(id int) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("У вас нет доступа к модулю #%v", id))
}

func errCommentTooLong() *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadGateway, "Комментарий слишком длинный")
}
