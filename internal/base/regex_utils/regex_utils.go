package regex_utils

import (
	"go-dianping/internal/base/regex_patterns"
	"regexp"
)

func IsPhoneInvalid(phone string) bool {
	re := regexp.MustCompile(regex_patterns.PhoneRegex)
	return !re.MatchString(phone)
}
