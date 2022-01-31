package support

import (
	"strconv"
	"strings"
)

const LatestSupportedAmphionVersion = "0.4.99"
const MinimumSupportedAmphionVersion = "0.3.0"

const ToolsVersion = "0.4.3"

func IsAmphionVersionSupported(ver string) bool {
	vLatest := comparisonNumber(stringVersionToNumber(LatestSupportedAmphionVersion))
	vMin := comparisonNumber(stringVersionToNumber(MinimumSupportedAmphionVersion))
	v := comparisonNumber(stringVersionToNumber(ver))

	return v >= vMin && v <= vLatest
}

func comparisonNumber(v1, v2, v3 int) int {
	return v1*10000 + v2*100 + v3
}

func stringVersionToNumber(ver string) (int, int, int) {
	numStr := strings.Split(ver, "-")[0]
	numStr = strings.ReplaceAll(numStr, "v", "")
	numStr = strings.ReplaceAll(numStr, "rc", "")
	numStr = strings.ReplaceAll(numStr, "preview", "")

	tokens := strings.Split(numStr, ".")

	if len(tokens) < 3 {
		return 0, 0, 0
	}

	num1, err := strconv.Atoi(tokens[0])
	if err != nil {
		num1 = 0
	}

	num2, err := strconv.Atoi(tokens[1])
	if err != nil {
		num2 = 0
	}

	num3, err := strconv.Atoi(tokens[2])
	if err != nil {
		num3 = 0
	}

	return num1, num2, num3
}
