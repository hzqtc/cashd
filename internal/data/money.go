package data

import (
	"fmt"
	"strings"
)

func FormatMoney(amount float64) string {
	s := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(s, ".")
	integer := parts[0]
	decimal := parts[1]
	return formatInteger(integer) + "." + decimal
}

func FormatMoneyInteger(amount float64) string {
	return formatInteger(fmt.Sprintf("%.0f", amount))
}

// Insert commas into integer on every 3 digits
// 1234 => 1,234; 1234567 => 1,234,567
func formatInteger(integer string) string {
	formattedInteger := ""
	for i, r := range integer {
		if i > 0 && (len(integer)-i)%3 == 0 {
			formattedInteger += ","
		}
		formattedInteger += string(r)
	}
	return formattedInteger

}
