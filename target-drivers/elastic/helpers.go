package elastic

import "strings"

func SanitizeJsonString(initial string) string {
	return strings.ReplaceAll(initial, `"`, `\"`)
}
