package utils

import (
	"context"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateNickname(s string) bool {
	pattern := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	res := pattern.MatchString(s)
	return res
}
func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}

func GetInterfaceLength(i interface{}) int {
	switch v := i.(type) {
	case []interface{}:
		return len(v)
	case []int:
		return len(v)
	case map[string]interface{}:
		return len(v)
	case string:
		return len(v)
	case interface{}:
		return 1
	default:
		return 0
	}
}
