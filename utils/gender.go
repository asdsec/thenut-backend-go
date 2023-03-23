package utils

const (
	MALE   = "m"
	FEMALE = "f"
)

// IsSupportedGender returns true if the provided gender is supported
func IsSupportedGender(gender string) bool {
	switch gender {
	case MALE, FEMALE:
		return true
	}
	return false
}
