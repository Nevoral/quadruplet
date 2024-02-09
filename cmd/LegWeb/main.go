package main

import (
	quat "github.com/Nevoral/quadrupot/internals/Quaternions"
	"github.com/Nevoral/quadrupot/internals/Robot"
	"github.com/Nevoral/quadrupot/internals/router"
	"github.com/gofiber/fiber/v3"
)

func main() {

	l := Robot.NewLeg(0, "nevim", &quat.Point3D{0, 0, 0})
	Robot.FrontLeftLeg(l)

	app := fiber.New()

	router.RouterLeg(l, app)

	//Start the server
	err := app.Listen(":8081")
	if err != nil {
		return
	}
}
