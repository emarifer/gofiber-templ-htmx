package handlers

import (
	"time"

	"github.com/a-h/templ"
	"github.com/emarifer/gofiber-templ-htmx/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sujit-baniya/flash"
)

var store *session.Store

const (
	AUTH_KEY       string = "authenticated"
	USER_ID        string = "user_id"
	FROM_PROTECTED string = "from_protected"
	TZONE_KEY      string = "time_zone"
)

func Setup(app *fiber.App) {
	/* Sessions Config */
	store = session.New(session.Config{
		CookieHTTPOnly: true,
		// CookieSecure: true, for https
		Expiration: time.Hour * 1,
	})

	/* Views */
	app.Get("/", flagsMiddleware, HandleViewHome)
	app.Get("/login", flagsMiddleware, HandleViewLogin)
	app.Post("/login", flagsMiddleware, HandleViewLogin)
	app.Get("/register", flagsMiddleware, HandleViewRegister)
	app.Post("/register", flagsMiddleware, HandleViewRegister)

	/* Views protected with session middleware */
	todoApp := app.Group("/todo", AuthMiddleware)
	todoApp.Get("/list", HandleViewList)
	todoApp.Get("/create", HandleViewCreatePage)
	todoApp.Post("/create", HandleViewCreatePage)
	todoApp.Get("/edit/:id", HandleViewEditPage)
	todoApp.Post("/edit/:id", HandleViewEditPage)
	todoApp.Delete("/delete/:id", HandleDeleteTodo)
	todoApp.Post("/logout", HandleLogout)

	/* ↓ Not Found Management - Fallback Page ↓ */
	app.Get("/*", flagsMiddleware, routeNotFoundHandler)
}

func routeNotFoundHandler(c *fiber.Ctx) error {
	fromProtected := c.Locals(FROM_PROTECTED).(bool)

	errorIndex := views.Error404(fromProtected)
	errorPage := views.Home(
		"", fromProtected, true, flash.Get(c), errorIndex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(errorPage))

	c.Status(fiber.StatusNotFound)

	return handler(c)
}

func flagsMiddleware(c *fiber.Ctx) error {
	sess, _ := store.Get(c)
	userId := sess.Get(USER_ID)
	if userId == nil {
		c.Locals(FROM_PROTECTED, false)

		return c.Next()
	}

	c.Locals(FROM_PROTECTED, true)

	return c.Next()
}
