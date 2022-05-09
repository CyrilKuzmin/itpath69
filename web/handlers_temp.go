package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// TEMPORARY handlers
func (w *Web) giveMeModules(c echo.Context) error {
	username := w.getUsernameIfAny(c)
	if username == "admin" {
		return errInternal(fmt.Errorf("fuck"))
	}
	user, err := w.userService.GetUserByName(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	if len(user.Modules) < w.moduleService.ModulesTotal() {
		err := w.userService.OpenNewModules(c.Request().Context(), username)
		if err != nil {
			return errInternal(err)
		}
	}
	return c.String(http.StatusOK, "OK")
}

func (w *Web) completeModule(c echo.Context) error {
	username := w.getUsernameIfAny(c)
	idParam := c.QueryParam("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return errInternal(err)
	}
	err = w.userService.MarkModuleAsCompleted(c.Request().Context(), username, id)
	if err != nil {
		return errInternal(err)
	}
	return c.String(http.StatusOK, "OK")
}
