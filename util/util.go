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
