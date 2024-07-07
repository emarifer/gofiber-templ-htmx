package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/emarifer/gofiber-templ-htmx/models"
	"github.com/emarifer/gofiber-templ-htmx/views"
	"github.com/emarifer/gofiber-templ-htmx/views/auth_views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sujit-baniya/flash"
	"golang.org/x/crypto/bcrypt"
)

/********** Handlers for Auth Views **********/

// Render Home Page
func HandleViewHome(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	hindex := views.HomeIndex(fromProtected)
	home := views.Home("", fromProtected, false, flash.Get(c), hindex)

	handler := adaptor.HTTPHandler(templ.Handler(home))

	return handler(c)
}

// Render Login Page with success/error messages & session management
func HandleViewLogin(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	lindex := auth_views.LoginIndex(fromProtected)
	login := auth_views.Login(
		" | Login", fromProtected, false, flash.Get(c), lindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(login))

	if c.Method() == "POST" {
		// obtaining the time zone from the POST request of the login form
		tzone := ""
		if len(c.GetReqHeaders()["X-Timezone"]) != 0 {
			tzone = c.GetReqHeaders()["X-Timezone"][0]
			// fmt.Println("Tzone:", tzone)
		}

		var (
			user models.User
			err  error
		)
		fm := fiber.Map{
			"type": "error",
		}

		// notice: in production you should not inform the user
		// with detailed messages about login failures
		if user, err = models.CheckEmail(c.FormValue("email")); err != nil {
			// fmt.Println(err)
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
			fm["message"] = "There is no user with that email"

			return flash.WithError(c, fm).Redirect("/login")
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(c.FormValue("password")),
		)
		if err != nil {
			fm["message"] = "Incorrect password"

			return flash.WithError(c, fm).Redirect("/login")
		}

		session, err := store.Get(c)
		if err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/login")
		}

		session.Set(AUTH_KEY, true)
		session.Set(USER_ID, user.ID)
		session.Set(TZONE_KEY, tzone)

		err = session.Save()
		if err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/login")
		}

		fm = fiber.Map{
			"type":    "success",
			"message": "You have successfully logged in!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/todo/list")
	}

	return handler(c)
}

// Render Register Page with success/error messages
func HandleViewRegister(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	rindex := auth_views.RegisterIndex(fromProtected)
	register := auth_views.Register(
		" | Register", fromProtected, false, flash.Get(c), rindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(register))

	if c.Method() == "POST" {
		user := models.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
			Username: c.FormValue("username"),
		}

		err := models.CreateUser(user)
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
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				err = errors.New("the email is already in use")
			}
			fm := fiber.Map{
				"type":    "error",
				"message": fmt.Sprintf("something went wrong: %s", err),
			}

			return flash.WithError(c, fm).Redirect("/register")
		}

		fm := fiber.Map{
			"type":    "success",
			"message": "You have successfully registered!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/login")
	}

	return handler(c)
}

// Authentication Middleware
func AuthMiddleware(c *fiber.Ctx) error {
	fm := fiber.Map{
		"type": "error",
	}

	session, err := store.Get(c)
	if err != nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	if session.Get(AUTH_KEY) == nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	userId := session.Get(USER_ID)
	if userId == nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	user, err := models.GetUserById(fmt.Sprint(userId.(uint64)))
	if err != nil {
		fm["message"] = "You are not authorized"

		return flash.WithError(c, fm).Redirect("/login")
	}

	c.Locals("userId", userId)
	c.Locals("username", user.Username)
	c.Locals(FROM_PROTECTED, true)
	// fromProtected = true

	return c.Next()
}

// Logout Handler
func HandleLogout(c *fiber.Ctx) error {
	fm := fiber.Map{
		"type": "error",
	}

	session, err := store.Get(c)
	if err != nil {
		fm["message"] = "logged out (no session)"

		return flash.WithError(c, fm).Redirect("/login")
	}

	err = session.Destroy()
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("/login")
	}

	fm = fiber.Map{
		"type":    "success",
		"message": "You have successfully logged out!!",
	}

	c.Locals(FROM_PROTECTED, false)
	// fromProtected = false

	return flash.WithSuccess(c, fm).Redirect("/login")
}
