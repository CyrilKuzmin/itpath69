package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/CyrilKuzmin/itpath69/internal/domain/comment"
	"github.com/CyrilKuzmin/itpath69/store"
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

func (w *Web) lkHandler(c echo.Context) error {
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
	return c.Render(http.StatusOK, "lk.html", map[string]interface{}{
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
	Comments []*comment.Comment
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
	comments := make(map[int][]*comment.Comment)
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

func (w *Web) newCommentHandler(c echo.Context) error {
	text := c.FormValue("text")
	moduleIdValue := c.FormValue("module_id")
	fmt.Println(text)
	moduleId, err := strconv.Atoi(moduleIdValue)
	if err != nil {
		return errInternal(err)
	}
	partIdValue := c.FormValue("part_id")
	partId, err := strconv.Atoi(partIdValue)
	if err != nil {
		return errInternal(err)
	}
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	// check if IDs are valid and allowed
	user, err := w.userService.GetUserByName(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	if moduleId > user.ModulesOpened {
		return errModuleNotAllowed(moduleId)
	}
	module, err := w.moduleService.GetModuleByID(c.Request().Context(), moduleId)
	if err != nil {
		return errInternal(err)
	}
	if partId > len(module.Data) {
		return errBadRequest()
	}
	// create comment
	comment, err := w.commentService.CreateComment(c.Request().Context(), username, text, moduleId, partId)
	if err != nil {
		return errInternal(err)
	}
	// send OK
	return c.JSON(http.StatusOK, comment)
}

func (w *Web) updateCommentHandler(c echo.Context) error {
	commentId := c.FormValue("comment_id")
	text := c.FormValue("text")
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	fmt.Println(username, commentId, text)
	// update comment
	err := w.commentService.UpdateComment(c.Request().Context(), username, commentId, text)
	if err != nil {
		return errInternal(err)
	}
	// send OK
	return c.String(http.StatusOK, "OK")
}

func (w *Web) deleteCommentHandler(c echo.Context) error {
	commentId := c.QueryParam("id")
	if commentId == "" {
		return errBadRequest()
	}
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	// delete comment
	err := w.commentService.DeleteCommentByID(c.Request().Context(), username, commentId)
	if err != nil {
		return errInternal(err)
	}
	// send OK
	return c.String(http.StatusOK, "OK")
}

func (w *Web) loginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	err := w.userService.CheckUserPassword(c.Request().Context(), username, password)
	if err != nil {
		return errLoginFailed()
	}
	w.setUserSession(c, username)
	return c.String(http.StatusOK, "OK")
}

func (w *Web) logoutHandler(c echo.Context) error {
	sess, err := w.session.Get(c.Request(), "session")
	if err != nil || sess.ID == "" {
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	w.deleteUserSession(c, sess)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func (w *Web) registerHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	_, err := w.userService.CreateUser(c.Request().Context(), username, password)
	if err != nil {
		if store.ErrorIs(err, store.AlreadyExistsErr) {
			return errUserAlreadyExists(username)
		} else {
			return errInternal(err)
		}
	}
	w.setUserSession(c, username)
	return c.String(http.StatusOK, "OK")
}

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
