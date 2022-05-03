package web

import (
	"html/template"
	"net/http"
	"strconv"

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
		c.String(http.StatusForbidden, "")
	}
	// load module
	module, err := w.moduleService.GetModuleByID(c.Request().Context(), id)
	if err != nil {
		return errInternal(err)
	}
	// need to convert string into template.HTML
	data := make([]template.HTML, len(module.Data))
	for i, p := range module.Data {
		data[i] = template.HTML(p.Data)
	}
	// render
	return c.Render(http.StatusOK, "module.html", map[string]interface{}{
		"Username":    username,
		"User":        user,
		"Module":      module.Meta,
		"Completed":   !user.Modules[module.Meta.Id].CompletedAt.IsZero(),
		"CompletedAt": user.Modules[module.Meta.Id].CompletedAt,
		"OpenedAt":    user.Modules[module.Meta.Id].CreatedAt,
		"Data":        data,
	})
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
	err := w.userService.OpenNewModules(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
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
