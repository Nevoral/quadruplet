package main

import (
	quat "github.com/Nevoral/quadrupot/internals/Quaternions"
	rob "github.com/Nevoral/quadrupot/internals/Robot"
	"github.com/Nevoral/quadrupot/internals/router"
	"github.com/gofiber/fiber/v3"
)

func main() {
	r := rob.NewRobot(&quat.Point3D{350, 0, 0}, &quat.Point3D{0, 0, 0}, &quat.Point3D{-210, 0, 0})

	app := fiber.New()

	router.Router(r, app)

	//soc := pythonAPI.NewBSDSocket()
	//soc.SendMessage("This is the best way to comunicate throught this socket")
	//go soc.OpenSocket("42069")
	//fmt.Println(len([]byte("This is the best way to comunicate throught this socket")))

	//Start the server
	err := app.Listen(":8081")
	if err != nil {
		return
	}
}
