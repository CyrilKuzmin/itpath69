package web

import (
	"net/http"
	"strconv"

	"github.com/CyrilKuzmin/itpath69/internal/service"
	"github.com/labstack/echo/v4"
)

const ModulesPerRow = 4

func (w *Web) indexHandler(c echo.Context) error {
	username := w.getUsernameIfAny(c)
	if username == "" {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
	}
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Username": username,
	})
}

func (w *Web) loginPageHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{})
}

func (w *Web) learnHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	// get the list of opened modules and show them
	user, err := w.srv.GetUserByName(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	metas, err := w.srv.ModulesPreview(c.Request().Context(), user, len(user.Modules))
	if err != nil {
		return errInternal(err)
	}
	shiftMetas(metas)
	rowsNum := len(metas) / ModulesPerRow
	if len(metas)%ModulesPerRow != 0 {
		rowsNum++
	}
	rows := make([][]service.ModuleDTO, rowsNum)
	for i := 0; i < rowsNum; i++ {
		row := metas[i*ModulesPerRow : i*ModulesPerRow+ModulesPerRow]
		rows[i] = row
	}
	// render all these structs
	return c.Render(http.StatusOK, "learn.html", map[string]interface{}{
		"User":             user,
		"Username":         user.Username, // for navbar
		"Rows":             rows,
		"ModulesTotal":     user.ModulesTotal,
		"ModulesOpened":    user.ModulesOpened,
		"ModulesCompleted": user.ModulesCompleted,
	})
}

func (w *Web) moduleHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	// get ID from URI
	idParam := c.QueryParam("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return errInternal(err)
	}
	// get user and check if he has permissions for this module
	user, err := w.srv.GetUserByName(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	module, err := w.srv.GetModuleForUser(c.Request().Context(), user, id)
	// render
	return c.Render(http.StatusOK, "module.html", map[string]interface{}{
		"Username":    username,
		"User":        user,
		"Module":      module, // comment form rendering bug
		"IsCompleted": module.IsCompleted,
		"CompletedAt": user.Modules[module.Id].CompletedAt,
		"OpenedAt":    user.Modules[module.Id].CreatedAt,
		"Data":        module.Data,
	})
}

func (w *Web) testingHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	// get ID from URI
	idParam := c.QueryParam("module_id")
	moduleId, err := strconv.Atoi(idParam)
	if err != nil {
		return errInternal(err)
	}
	// Optional param. It's not ampty if we wanna to continue the test (from account page)
	testId := c.QueryParam("test_id")
	test, err := w.srv.GetTestByID(c.Request().Context(), testId, true)
	return c.Render(http.StatusOK, "testing.html", map[string]interface{}{
		"Username": username,
		"Module":   moduleId,
		"Test":     test,
	})
}
