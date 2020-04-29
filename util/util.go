package util

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

func GetStatus(status int) string {
	switch status {
	case 0:
		return "Pending"
	case 1:
		return "Successful"
	case 2:
		return "Unsuccessful"

	default:
		return "Pending"
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
