// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	app := fiber.New()

// 	receiptGroup := app.Group("/receipts")

// 	receiptGroup.Post("/process", processReceipt)
// 	receiptGroup.Get("/:id/points", getPoints)

// 	log.Fatal(app.Listen(":7777"))
// }

// func processReceipt(c *fiber.Ctx) error {
// 	fmt.Println("process receipt")

// 	// TODO: validation

// 	// TODO: create a new ID

// 	// TODO: calculate points

// 	// TODO: store the result

// 	// TODO: return the ID

// 	return nil
// }

// func getPoints(c *fiber.Ctx) error {
// 	fmt.Println("get points")
// 	fmt.Println(c.Params("id"))

// 	// TODO: validation

// 	// TODO: get points from the storage

// 	// TODO: return the result

// 	return nil
// }

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
