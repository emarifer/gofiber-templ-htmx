package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/emarifer/gofiber-templ-htmx/models"
	"github.com/emarifer/gofiber-templ-htmx/views/todo_views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sujit-baniya/flash"
)

/********** Handlers for Todo Views **********/

// Render List Page with success/error messages
func HandleViewList(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	todo := new(models.Todo)
	todo.CreatedBy = c.Locals("userId").(uint64)

	// fm := fiber.Map{"type": "error"}

	todosSlice, err := todo.GetAllTodos()
	if err != nil {
		if strings.Contains(err.Error(), "no such table") ||
			strings.Contains(err.Error(), "database is locked") {
			// "no such table" is the error that SQLite3 produces
			// when some table does not exist, and we have only
			// used it as an example of the errors that can be caught.
			// Here you can add the errors that you are interested
			// in throwing as `500` codes.
			return fiber.NewError(
				fiber.StatusServiceUnavailable,
				"database temporarily out of service",
			)
		}
		// fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		// return flash.WithError(c, fm).Redirect("/todo/list")
	}

	tindex := todo_views.TodoIndex(todosSlice)
	tlist := todo_views.TodoList(
		" | Tasks List",
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		tindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(tlist))

	return handler(c)
}

// Render Create Todo Page with success/error messages
func HandleViewCreatePage(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	if c.Method() == "POST" {
		newTodo := new(models.Todo)
		newTodo.CreatedBy = c.Locals("userId").(uint64)
		newTodo.Title = strings.Trim(c.FormValue("title"), " ")
		newTodo.Description = strings.Trim(c.FormValue("description"), " ")

		fm := fiber.Map{
			"type":    "error",
			"message": "Task title empty!!",
		}
		if newTodo.Title == "" {

			return flash.WithError(c, fm).Redirect("/todo/list")
		}

		if _, err := newTodo.CreateTodo(); err != nil {
			if strings.Contains(err.Error(), "no such table") ||
				strings.Contains(err.Error(), "database is locked") {
				// "no such table" is the error that SQLite3 produces
				// when some table does not exist, and we have only
				// used it as an example of the errors that can be caught.
				// Here you can add the errors that you are interested
				// in throwing as `500` codes.
				return fiber.NewError(
					fiber.StatusServiceUnavailable,
					"database temporarily out of service",
				)
			}
		}

		fm = fiber.Map{
			"type":    "success",
			"message": "Task successfully created!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/todo/list")
	}

	cindex := todo_views.CreateIndex()
	create := todo_views.Create(
		" | Create Todo",
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		cindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(create))

	return handler(c)
}

// Render Edit Todo Page with success/error messages
func HandleViewEditPage(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)
	session, _ := store.Get(c)
	tzone := session.Get(TZONE_KEY).(string)

	idParams, _ := strconv.Atoi(c.Params("id"))
	todoId := uint64(idParams)

	todo := new(models.Todo)
	todo.ID = todoId
	todo.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{"type": "error"}

	recoveredTodo, err := todo.GetNoteById()

	if err != nil {
		if strings.Contains(err.Error(), "no such table") ||
			strings.Contains(err.Error(), "database is locked") {
			// "no such table" is the error that SQLite3 produces
			// when some table does not exist, and we have only
			// used it as an example of the errors that can be caught.
			// Here you can add the errors that you are interested
			// in throwing as `500` codes.
			return fiber.NewError(
				fiber.StatusServiceUnavailable,
				"database temporarily out of service",
			)
		}

		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("/todo/list")
	}

	if c.Method() == "POST" {
		todo.Title = strings.Trim(c.FormValue("title"), " ")
		todo.Description = strings.Trim(c.FormValue("description"), " ")
		if c.FormValue("status") == "on" {
			todo.Status = true
		} else {
			todo.Status = false
		}

		fm = fiber.Map{
			"type":    "error",
			"message": "Task title empty!!",
		}
		if todo.Title == "" {

			return flash.WithError(c, fm).Redirect("/todo/list")
		}

		_, err := todo.UpdateTodo()
		if err != nil {
			if strings.Contains(err.Error(), "no such table") ||
				strings.Contains(err.Error(), "database is locked") {
				// "no such table" is the error that SQLite3 produces
				// when some table does not exist, and we have only
				// used it as an example of the errors that can be caught.
				// Here you can add the errors that you are interested
				// in throwing as `500` codes.
				return fiber.NewError(
					fiber.StatusServiceUnavailable,
					"database temporarily out of service",
				)
			}

			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/todo/list")
		}

		fm = fiber.Map{
			"type":    "success",
			"message": "Task successfully updated!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/todo/list")
	}

	uindex := todo_views.UpdateIndex(recoveredTodo, tzone)
	update := todo_views.Update(
		fmt.Sprintf(" | Edit Todo #%d", recoveredTodo.ID),
		fromProtected,
		false,
		flash.Get(c),
		c.Locals("username").(string),
		uindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(update))

	return handler(c)
}

// Handler Remove Todo
func HandleDeleteTodo(c *fiber.Ctx) error {
	idParams, _ := strconv.Atoi(c.Params("id"))
	todoId := uint64(idParams)

	todo := new(models.Todo)
	todo.ID = todoId
	todo.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{"type": "error"}

	if err := todo.DeleteTodo(); err != nil {
		if strings.Contains(err.Error(), "no such table") ||
			strings.Contains(err.Error(), "database is locked") {
			// "no such table" is the error that SQLite3 produces
			// when some table does not exist, and we have only
			// used it as an example of the errors that can be caught.
			// Here you can add the errors that you are interested
			// in throwing as `500` codes.
			return fiber.NewError(
				fiber.StatusServiceUnavailable,
				"database temporarily out of service",
			)
		}
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect(
			"/todo/list",
			fiber.StatusSeeOther,
		)
	}

	fm = fiber.Map{
		"type":    "success",
		"message": "Task successfully deleted!!",
	}

	return flash.WithSuccess(c, fm).Redirect("/todo/list", fiber.StatusSeeOther)
}
