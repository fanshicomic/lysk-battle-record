package pkg

import "fmt"

func GetValue(input map[string]interface{}, field string) string {
	if val, ok := input[field]; ok {
		return fmt.Sprintf("%v", val)
	}

	return ""
}
