package api

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/isntdoris/fetch_backend_take_home/businesslogic"
	"github.com/isntdoris/fetch_backend_take_home/model"
)

// Regular expressions for input validation
var (
	// retailerRegex = regexp.MustCompile(`^\S+$`)
	retailerRegex         = regexp.MustCompile(`^[\w\s\-]+$`)
	shortDescriptionRegex = regexp.MustCompile(`^[\w\s\-]+$`)
	priceRegex            = regexp.MustCompile(`^\d+\.\d{2}$`)
	dateRegex             = regexp.MustCompile(`^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`)
	timeRegex             = regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`)
)

func ProcessReceipt(c *fiber.Ctx) error {

	// 1: validation
	// 2: create a new ID
	// 3: calculate points
	// 4: store the result
	// 5: return the ID

	// request body parsing
	receiptRequest := model.ReceiptRequest{}

	if err := c.BodyParser(&receiptRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if receiptRequest.Retailer == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Retailer field is required"})
	}
	if receiptRequest.PurchaseDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "PurchaseDate field is required"})
	}
	if receiptRequest.PurchaseTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "PurchaseTime field is required"})
	}
	if len(receiptRequest.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Items field is required and cannot be empty"})
	}
	if receiptRequest.Total == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Total field is required"})
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

	// Additional check for the presence of items.
	if len(receiptRequest.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Items field is required and cannot be empty"})
	}

	receiptId, err := businesslogic.ProcessReceipt(&receiptRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the ID and points to the client
	return c.JSON(fiber.Map{"id": receiptId.String()})
}

func GetPoints(c *fiber.Ctx) error {
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

	// TODO
	points, err := businesslogic.GetPoints(receiptId)
	if err == businesslogic.ErrorReceiptNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the points associated with the receipt ID
	return c.JSON(fiber.Map{"points": points})
}
