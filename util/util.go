package util

import "math/rand"

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

func GetMinDistance(distanceArray []int) int {
	min := distanceArray[0]
	for _, value := range distanceArray {
		if value < min {
			min = value
		}
	}
	return min
}

func RandomString(n int) string {
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
		return "Order Received"
	case 2:
		return "Order Processing"
	case 3:
		return "Shipped"
	case 4:
		return "Out For Delivery"
	case 5:
		return "Delivered"

	default:
		return "Order Received"
	}
}
