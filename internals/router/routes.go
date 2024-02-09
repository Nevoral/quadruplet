package router

import (
	"github.com/Nevoral/quadrupot/internals/Robot"
	"github.com/Nevoral/quadrupot/internals/handlers"
	"github.com/gofiber/fiber/v3"
)

func Router(r *Robot.Robot, app *fiber.App) {
	app.Get("/", r.SendGraphPage())
	app.Get("/js/:file", handlers.GetJs)
	app.Get("/assets/:file", handlers.GetAsset)
	app.Get("/defaultConfig", r.SendDefaultConfig())
	app.Post("/", r.SendGraphData())
}

func RouterLeg(l *Robot.Leg, app *fiber.App) {
	app.Get("/", l.SendGraphPageLeg())
	app.Get("/js/:file", handlers.GetJs)
	app.Get("/assets/:file", handlers.GetAsset)
	app.Post("/", l.SendGraphData())
}
