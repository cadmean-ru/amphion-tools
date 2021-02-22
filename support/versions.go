package support

import (
	"strconv"
	"strings"
)

const LatestSupportedAmphionVersion = "0.1.99"
const MinimumSupportedAmphionVersion = "0.1.10"

func IsAmphionVersionSupported(ver string) bool {
	latestNum1, latestNum2, latestNum3 := stringVersionToNumber(LatestSupportedAmphionVersion)
	minNum1, minNum2, minNum3 := stringVersionToNumber(MinimumSupportedAmphionVersion)
	verNum1, verNum2, verNum3 := stringVersionToNumber(ver)
	return verNum1 >= minNum1 && verNum1 <= latestNum1 &&
		verNum2 >= minNum2 && verNum2 <= latestNum2 &&
		verNum3 >= minNum3 && verNum3 <= latestNum3
}

func stringVersionToNumber(ver string) (int, int, int) {
	numStr := strings.ReplaceAll(ver, "v", "")
	numStr = strings.ReplaceAll(numStr, "rc", "")

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