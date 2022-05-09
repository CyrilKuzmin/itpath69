package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CyrilKuzmin/itpath69/internal/domain/tests"
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
	user, err := w.userService.GetUserByName(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	test, err := w.testsService.GenerateTest(c.Request().Context(), user.Id, moduleId, tests.DefaultQuestionsAmount)
	return c.JSON(http.StatusOK, test)
}

func (w *Web) listTestsHandler(c echo.Context) error {
	return nil
}

type UserResult struct {
	Result float32 `json:"result"`
}

func (w *Web) checkTestHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	userAnswers := *&tests.Test{}
	err := json.NewDecoder(c.Request().Body).Decode(&userAnswers)
	if err != nil {
		return errBadRequest()
	}
	res, err := w.testsService.CheckTest(c.Request().Context(), userAnswers.Id, userAnswers.Questions)
	if err != nil {
		return errInternal(err)
	}
	return c.JSON(http.StatusOK, UserResult{res})
}
