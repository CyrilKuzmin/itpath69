package web

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

const MaximumCommentLength = 10000

func (w *Web) getCommentsHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	moduleIdValue := c.QueryParam("module_id")
	moduleId, err := strconv.Atoi(moduleIdValue)
	if err != nil {
		return errInternal(err)
	}
	// check if IDs are valid and allowed
	user, err := w.userService.GetUserByName(c.Request().Context(), username)
	if err != nil {
		return errInternal(err)
	}
	if moduleId > user.ModulesOpened {
		return errModuleNotAllowed(moduleId)
	}
	comments, err := w.commentService.ListCommentsByModule(c.Request().Context(), username, moduleId)
	if err != nil {
		return errInternal(err)
	}
	return c.JSON(http.StatusOK, comments)
}

func (w *Web) newCommentHandler(c echo.Context) error {
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	text := c.FormValue("text")
	if len(text) > MaximumCommentLength {
		return errCommentTooLong()
	}
	moduleIdValue := c.FormValue("module_id")
	moduleId, err := strconv.Atoi(moduleIdValue)
	if err != nil {
		return errInternal(err)
	}
	partIdValue := c.FormValue("part_id")
	partId, err := strconv.Atoi(partIdValue)
	if err != nil {
		return errInternal(err)
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
	commentId := c.FormValue("id")
	text := c.FormValue("text")
	// redirect to login page if no session found
	username := w.getUsernameIfAny(c)
	if username == "" {
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
	// update comment
	newComment, err := w.commentService.UpdateComment(c.Request().Context(), username, commentId, text)
	if err != nil {
		return errInternal(err)
	}
	// send OK
	return c.JSON(http.StatusOK, newComment)
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
