package elastic

import (
	"strconv"
	"strings"
)

func GetPrimaryKeysAsString(keys []string) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, str := range keys {
		sb.WriteString(strconv.Quote(str))
		if i < len(keys)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")
	return strings.ReplaceAll(sb.String(), `"`, `\"`)
}
