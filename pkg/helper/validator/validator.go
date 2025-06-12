package validator

import (
	"regexp"
)

func IsPhone(str string) bool {
	re := regexp.MustCompile(`^1([38][0-9]|4[579]|5[0-3,5-9]|66|7[0135678]|9[89])\d{8}$`)
	return re.MatchString(str)
}
