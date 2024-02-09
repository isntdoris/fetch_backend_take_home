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
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"math"
	"strconv"
	"strings"
	"time"
)

var pointsStore = make(map[string]int)

// Regular expressions for input validation
var (
	// retailerRegex = regexp.MustCompile(`^\S+$`)
	retailerRegex         = regexp.MustCompile(`^[a-zA-Z0-9 &'-]+$`)
	shortDescriptionRegex = regexp.MustCompile(`^[\w\s\-]+$`)
	priceRegex            = regexp.MustCompile(`^\d+\.\d{2}$`)
	dateRegex             = regexp.MustCompile(`^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`)
	timeRegex             = regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)
)

func main() {
	app := fiber.New()

	receiptGroup := app.Group("/receipts")

	receiptGroup.Post("/process", processReceipt)
	receiptGroup.Get("/:id/points", getPoints)

	log.Fatal(app.Listen(":7777"))
}

func processReceipt(c *fiber.Ctx) error {

	// 1: validation
	// 2: create a new ID
	// 3: calculate points
	// 4: store the result
	// 5: return the ID

	// request body parsing
	type ReceiptItem struct {
		Description string `json:"shortDescription"`
		Price       string `json:"price"`
	}
	var receiptRequest struct {
		Retailer     string        `json:"retailer"`
		PurchaseDate string        `json:"purchaseDate"`
		PurchaseTime string        `json:"purchaseTime"`
		Items        []ReceiptItem `json:"items"`
		Total        string        `json:"total"`
	}
	if err := c.BodyParser(&receiptRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Apply regular expression validations
	if !retailerRegex.MatchString(receiptRequest.Retailer) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid retailer format"})
	}
	if !dateRegex.MatchString(receiptRequest.PurchaseDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid purchase date format"})
	}
	if !timeRegex.MatchString(receiptRequest.PurchaseTime) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid purchase time format"})
	}
	for _, item := range receiptRequest.Items {
		if !shortDescriptionRegex.MatchString(item.Description) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item description format"})
		}
		if !priceRegex.MatchString(item.Price) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item price format"})
		}
	}

	// Start point calculation
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	// points += len(regexp.MustCompile(`\w`).FindAllString(receiptRequest.Retailer, -1))
	alphanumericRegex := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(alphanumericRegex.FindAllString(receiptRequest.Retailer, -1))

	// Rule 2 & 3: 50 points for round dollar, 25 points for multiple of 0.25
	total, err := strconv.ParseFloat(receiptRequest.Total, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid total format"})
	}
	if total == math.Floor(total) {
		points += 50 // Round dollar
	}
	if math.Mod(total*100, 25) == 0 {
		points += 25 // Multiple of 0.25
	}

	// Rule 4: 5 points for every two items
	points += (len(receiptRequest.Items) / 2) * 5

	// Rule 5: Points for items with description length multiple of 3
	for _, item := range receiptRequest.Items {
		if len(strings.TrimSpace(item.Description))%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item price format"})
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: 6 points if the purchase date day is odd
	purchaseDate, err := time.Parse("2006-01-02", receiptRequest.PurchaseDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid purchase date format"})
	}
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the purchase time is between 2:00pm and 4:00pm
	purchaseTime, err := time.Parse("15:04", receiptRequest.PurchaseTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid purchase time format"})
	}
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	// ID generation
	receiptId := uuid.New()
	// Simulate points calculation
	// points := 100 // Example fixed points value for demonstration

	pointsStore[receiptId.String()] = points

	// Return the ID and points to the client
	return c.JSON(fiber.Map{"id": receiptId.String(), "points": points})
}

func getPoints(c *fiber.Ctx) error {
	// 1: validation
	// 2: get points from the storage
	// 3: return the result

	// Retrieve receiptId from the URL parameters
	receiptIdParam := c.Params("id")

	// Validate the UUID format of the receiptId
	receiptId, err := uuid.Parse(receiptIdParam)
	if err != nil {
		// If the ID is not a valid UUID, return a bad request error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid receipt ID format"})
	}

	// Check if the receipt exists in the pointsStore
	points, exists := pointsStore[receiptId.String()]
	if !exists {
		// If the receipt ID does not exist in the points store, return a not found error
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Receipt not found"})
	}

	// Return the points associated with the receipt ID
	return c.JSON(fiber.Map{"id": receiptId.String(), "points": points})
}
