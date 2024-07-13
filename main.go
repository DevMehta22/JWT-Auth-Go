package main

import (
	"fmt"
	"github.com/devmehta22/JWT-Auth/database"
	"github.com/devmehta22/JWT-Auth/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main()  {
	_,err := database.ConnectDB()
	if err != nil {
		fmt.Println(err)
		panic("Error connecting DB")	
	}
	fmt.Println("Connection is Successful")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders:     "Content-Type,Authorization,Accept,Origin,Access-Control-Request-Method,Access-Control-Request-Headers,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Access-Control-Allow-Methods,Access-Control-Expose-Headers,Access-Control-Max-Age,Access-Control-Allow-Credentials",
		AllowCredentials: true,
		AllowOrigins: "http://localhost:8000",
		}))
	routes.SetUpRoutes(app)

	err = app.Listen(":8000")
	if err != nil {
		panic("Error starting server")
		}

}