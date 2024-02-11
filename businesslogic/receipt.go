package businesslogic

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/isntdoris/fetch_backend_take_home/model"
)

// var pointsStore = make(map[string]int)
var pointsStore sync.Map

var (
	ErrorReceiptNotFound = errors.New("receipt not found")
)

func ProcessReceipt(receiptRequest *model.ReceiptRequest) (uuid.UUID, error) {

	// Start point calculation
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	// points += len(regexp.MustCompile(`\w`).FindAllString(receiptRequest.Retailer, -1))
	alphanumericRegex := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(alphanumericRegex.FindAllString(receiptRequest.Retailer, -1))

	// Rule 2 & 3: 50 points for round dollar, 25 points for multiple of 0.25
	total, err := strconv.ParseFloat(receiptRequest.Total, 64)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid total format")
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
				return uuid.UUID{}, errors.New("invalid item price format")
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: 6 points if the purchase date day is odd
	purchaseDate, err := time.Parse("2006-01-02", receiptRequest.PurchaseDate)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid purchase date format")
	}
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the purchase time is between 2:00pm and 4:00pm
	purchaseTime, err := time.Parse("15:04", receiptRequest.PurchaseTime)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid purchase time format")
	}
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	// ID generation
	receiptId := uuid.New()
	// Simulate points calculation
	// points := 100 // Example fixed points value for demonstration

	// pointsStore[receiptId.String()] = points
	pointsStore.Store(receiptId.String(), points)

	return receiptId, nil
}

func GetPoints(receiptId uuid.UUID) (int, error) {

	// points, exists := pointsStore[receiptId.String()]
	// Attempt to load points from the store
	pointsValue, ok := pointsStore.Load(receiptId.String())
	if !ok {
		// If the receipt ID does not exist in the points store, return a not found error
		return 0, ErrorReceiptNotFound
	}

	// If the receipt ID exists, return the points
	points, ok := pointsValue.(int)
	if !ok {
		// If point value type is not correct:
		return 0, errors.New("invalid points type")
	}

	return points, nil
}
