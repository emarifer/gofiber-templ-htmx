package main

import (
	"log"

	"github.com/emarifer/gofiber-templ-htmx/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		// Setting centralized error hanling.
		ErrorHandler: handlers.CustomErrorHandler,
	})

	app.Static("/", "./assets")

	app.Use(logger.New())

	handlers.Setup(app)

	log.Fatal(app.Listen(":3000"))
}

/* REFERENCES:
https://www.youtube.com/watch?v=Ck919fGGbCw
http://cryto.net/~joepie91/blog/2016/06/13/stop-using-jwt-for-sessions/
http://cryto.net/~joepie91/blog/2016/06/19/stop-using-jwt-for-sessions-part-2-why-your-solution-doesnt-work/

https://github.com/NerdCademyDev/golang/tree/main/23_server_session_auth

Attributes in Templ Components. SEE:
https://templ.guide/syntax-and-usage/attributes/

https://github.com/gofiber/fiber/issues/503
https://docs.gofiber.io/api/ctx/#locals

https://docs.gofiber.io/guide/grouping/
https://github.com/gofiber/fiber/issues/1179
https://docs.gofiber.io/extra/faq/#how-do-i-handle-custom-404-responses
https://docs.gofiber.io/guide/routing/#middleware

https://www.sqlite.org/foreignkeys.html

https://stackoverflow.com/questions/72411062/controlling-indents-in-go-templates

https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/303

https://www.sqlitetutorial.net/sqlite-update/

https://stackoverflow.com/questions/26152088/access-a-map-value-using-a-variable-key-in-a-go-template
*/
