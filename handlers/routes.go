package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var (
	store         *session.Store
	AUTH_KEY      string = "authenticated"
	USER_ID       string = "user_id"
	fromProtected bool   = false
)

func Setup(app *fiber.App) {
	/* Sessions Config */
	store = session.New(session.Config{
		CookieHTTPOnly: true,
		// CookieSecure: true, for https
		Expiration: time.Hour * 1,
	})

	/* Views */
	app.Get("/", HandleViewHome)
	app.Get("/login", HandleViewLogin)
	app.Post("/login", HandleViewLogin)
	app.Get("/register", HandleViewRegister)
	app.Post("/register", HandleViewRegister)

	/* Views protected with session middleware */
	todoApp := app.Group("/todo", AuthMiddleware)
	todoApp.Get("/list", HandleViewList)
	todoApp.Get("/create", HandleViewCreatePage)
	todoApp.Post("/create", HandleViewCreatePage)
	todoApp.Get("/edit/:id", HandleViewEditPage)
	todoApp.Post("/edit/:id", HandleViewEditPage)
	todoApp.Delete("/delete/:id", HandleDeleteTodo)
	todoApp.Post("/logout", HandleLogout)

	/* Page Not Found Management */
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendFile("./views/404.html")
	})
}
