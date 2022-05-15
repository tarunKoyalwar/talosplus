package alerts

import (
	"strings"
)

// use spycolor.com for decimal codes
var (
	ErrColor  = 12538964
	NormColor = 5955808
)

// FormatMsg : Split Msg if too long
func FormatMsg(dat string) []string {
	arr := []string{}

	counter := 0
	temp := ""

	for _, v := range strings.Fields(dat) {

		if counter > 1800 {

			arr = append(arr, temp)
			temp = ""
			counter = 0

		}

		temp = temp + v + " "
		counter += len(v) + 1

	}

	if temp != "" {
		arr = append(arr, temp)
	}

	return arr

}
