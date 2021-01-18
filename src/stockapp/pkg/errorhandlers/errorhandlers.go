package errorhandlers

import (
	"strings"
)

func ConfirmBody(inputText string) string {
	if strings.Contains(inputText, "\"error\":\"Invalid ApiKey\"") {
		return "Invalid key provided"
	}
	return ""
}
