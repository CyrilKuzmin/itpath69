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
	testId := c.QueryParam("test_id")
	var test *tests.Test
	if testId != "" {
		test, err = w.testsService.GetTestByID(c.Request().Context(), testId, true)
	} else {
		test, err = w.testsService.GenerateTest(c.Request().Context(), user.Id, moduleId, tests.DefaultQuestionsAmount)
	}
	return c.JSON(http.StatusOK, test)
}

func (w *Web) listTestsHandler(c echo.Context) error {
	return nil
}

type UserResult struct {
	Score    float32 `json:"score"`
	IsPassed bool    `json:"is_passed"`
}

func (w *Web) checkTestHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	userTestData := &tests.Test{}
	err := json.NewDecoder(c.Request().Body).Decode(&userTestData)
	if err != nil {
		return errBadRequest()
	}
	score, err := w.testsService.CheckTest(c.Request().Context(), userTestData.Id, userTestData.Questions)
	if err != nil {
		return errInternal(err)
	}
	var isPassed bool
	if score > tests.DefaultPassThreshold {
		// move it to service
		err = w.userService.MarkModuleAsCompleted(c.Request().Context(), username, userTestData.ModuleId)
		if err != nil {
			return errInternal(err)
		}
		isPassed = true
	}
	err = w.testsService.MarkTestExpired(c.Request().Context(), userTestData.Id)
	if err != nil {
		return errInternal(err)
	}
	return c.JSON(http.StatusOK, UserResult{score, isPassed})
}
