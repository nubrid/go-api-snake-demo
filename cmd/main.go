package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nubrid/go-api-snake-demo/internal/handlers"
)

func main() {
	// const express = require("express")
	// const app = express()
	// app.use(express.json())
	app := fiber.New()

	app.Get("/new", handlers.CreateNewGame)
	app.Post("/validate", handlers.ValidateMoveSet)

	// try { app.listen(3000) } catch (err) { console.log(err); process.exit(1) }
	log.Fatal(app.Listen(":3000"))
}
