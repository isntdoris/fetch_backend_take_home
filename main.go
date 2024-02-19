package main

import (
	// "fmt"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/isntdoris/fetch_backend_take_home/api"
)

func main() {
	app := fiber.New()

	receiptGroup := app.Group("/receipts")

	receiptGroup.Post("/process", api.ProcessReceipt)
	receiptGroup.Get("/:id/points", api.GetPoints)

	log.Fatal(app.Listen(":7777"))
}
