package prettier

import (
	"fmt"
	"strings"
)

const (
	PlaceHolderDollar   = "$"
	PlaceholderQuestion = "?"
)

func Pretty(query string, placeholder string, args ...any) string {
	for i, param := range args {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("'%s'", v) // строки в SQL оборачиваются в кавычки
		case []byte:
			value = fmt.Sprintf("'%s'", string(v)) // преобразование в строку
		default:
			value = fmt.Sprintf("%v", v) // просто подставляем значение
		}

		query = strings.Replace(query, fmt.Sprintf("%s%d", placeholder, i+1), value, 1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.TrimSpace(query)
}
