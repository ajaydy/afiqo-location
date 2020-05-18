package util

import (
	"afiqo-location/api"
	"github.com/google/uuid"
	"math"
	"math/rand"
	"time"
)

func GetGender(gender int) string {
	switch gender {
	case 1:
		return "Male"
	case 2:
		return "Female"

	default:
		return ""
	}
}

func MeterToKilometer(m int) float64 {
	kilometer := float64(m) / 1000
	return math.Round(kilometer*100) / 100
}

func count(number int) int {
	count := 0
	for number != 0 {
		number /= 10
		count += 1
	}
	return count
}

func GetOrderStatus(status int) string {
	switch status {
	case 0:
		return "Open"
	case 1:
		return "Confirmed"
	case 2:
		return "Completed"

	default:
		return "Open"
	}
}

func GetPaymentStatus(status int) string {
	switch status {
	case 0:
		return "Unpaid"
	case 1:
		return "Paid"
	default:
		return "Unpaid"
	}
}

func GetMinDistance(distances []api.Distance) (float64, uuid.UUID) {
	min := distances[0].DistanceValue
	id := distances[0].WarehouseID
	for _, distance := range distances {
		if distance.DistanceValue < min {
			min = distance.DistanceValue
			id = distance.WarehouseID
		}
	}
	return MeterToKilometer(min), id
}

func RandomString(n int) string {

	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func GetShipmentStatus(status int) string {
	switch status {
	case 1:
		return "Order Processing"
	case 2:
		return "Shipped"
	case 3:
		return "Out For Delivery"
	case 4:
		return "Delivered"

	default:
		return "Order Processing"
	}
}
