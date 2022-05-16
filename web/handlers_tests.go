package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/CyrilKuzmin/itpath69/internal/service"
	"github.com/labstack/echo/v4"
)

func (w *Web) getTestHandler(c echo.Context) error {
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
	testId := c.QueryParam("test_id")
	var test *service.TestDTO
	if testId != "" {
		test, err = w.srv.GetTestByID(c.Request().Context(), testId, true)
	} else {
		fmt.Println("creating new test from server GET /test")
		test, err = w.srv.CreateNewTest(c.Request().Context(), username, moduleId)
	}
	return c.JSON(http.StatusOK, test)
}

func (w *Web) listTestsHandler(c echo.Context) error {
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	tests, err := w.srv.ListTestsByUsername(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	return c.JSON(http.StatusOK, tests)
}

func (w *Web) checkTestHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	res, err := w.srv.CheckTest(c.Request().Context(), username, c.Request().Body)
	if err != nil {
		return errInternal(err)
	}
	return c.JSON(http.StatusOK, res)
}
