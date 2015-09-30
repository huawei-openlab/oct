package specsConvert

import (
	"regexp"
	"strconv"
)

func SizeStringToByte(val string) int64 {
	re, _ := regexp.Compile("^(\\d+)(g|G|k|K|m|M)[b|B]?$")

	result := re.FindStringSubmatch(val)
	if len(result) != 3 {
		return -1
	}

	baseNum, _ := strconv.ParseInt(result[1], 10, 64)
	switch result[2] {
	case "g":
		fallthrough
	case "G":
		baseNum *= 1024 * 1024 * 1024
	case "m":
		fallthrough
	case "M":
		baseNum *= 1024 * 1024
	case "k":
		fallthrough
	case "K":
		baseNum *= 1024
	}
	return baseNum
}
