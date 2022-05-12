package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/CyrilKuzmin/itpath69/internal/domain/comment"
	"github.com/labstack/echo/v4"
)

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
	user, err := w.userService.GetUserByName(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	modulesMeta, err := w.moduleService.ModulesPreview(c.Request().Context(), len(user.Modules))
	if err != nil {
		return errInternal(err)
	}
	// mark completed modules
	for i, row := range modulesMeta {
		for j, m := range row {
			if !user.Modules[m.Id].CompletedAt.IsZero() {
				modulesMeta[i][j].Completed = true
			}
		}
	}
	// render all these structs
	return c.Render(http.StatusOK, "learn.html", map[string]interface{}{
		"User":             user,
		"Username":         user.Username, // for navbar
		"Rows":             modulesMeta,
		"ModulesTotal":     w.moduleService.ModulesTotal(),
		"ModulesOpened":    user.ModulesOpened,
		"ModulesCompleted": user.ModulesCompleted,
	})
}

type ModulePart struct {
	Id       int
	ModuleId int // comment form rendering bug
	Data     template.HTML
	Comments []*comment.CommentDTO
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
	user, err := w.userService.GetUserByName(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	if id > len(user.Modules) {
		return errModuleNotAllowed(id)
	}
	// load module
	module, err := w.moduleService.GetModuleByID(c.Request().Context(), id)
	if err != nil {
		return errInternal(err)
	}
	// list comments for module
	cmts, err := w.commentService.ListCommentsByModule(c.Request().Context(), username, id)
	if err != nil {
		return errInternal(err)
	}
	comments := make(map[int][]*comment.CommentDTO)
	for _, c := range cmts {
		comments[c.PartId] = append(comments[c.PartId], c)
	}
	// need to convert string into template.HTML and add comments
	data := make([]ModulePart, len(module.Data))
	for i, p := range module.Data {
		data[i] = ModulePart{
			Id:       p.Id,
			Data:     template.HTML(p.Data),
			Comments: comments[p.Id],
			ModuleId: module.Id,
		}
	}
	// render
	return c.Render(http.StatusOK, "module.html", map[string]interface{}{
		"Username":    username,
		"User":        user,
		"Module":      module.Meta, // comment form rendering bug
		"Completed":   !user.Modules[module.Meta.Id].CompletedAt.IsZero(),
		"CompletedAt": user.Modules[module.Meta.Id].CompletedAt,
		"OpenedAt":    user.Modules[module.Meta.Id].CreatedAt,
		"Data":        data,
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
	testId := c.QueryParam("test_id")
	test, err := w.testsService.GetTestByID(c.Request().Context(), testId)
	fmt.Println(test)
	return c.Render(http.StatusOK, "testing.html", map[string]interface{}{
		"Username": username,
		"Module":   moduleId,
		"Test":     test,
	})
}
