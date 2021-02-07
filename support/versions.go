package support

import (
	"strconv"
	"strings"
)

const LatestSupportedAmphionVersion = "0.1.7"
const MinimumSupportedAmphionVersion = "0.1.7"

func IsAmphionVersionSupported(ver string) bool {
	latestNum := stringVersionToNumber(LatestSupportedAmphionVersion)
	minNum := stringVersionToNumber(MinimumSupportedAmphionVersion)
	verNum := stringVersionToNumber(ver)
	return verNum >= minNum && verNum <= latestNum
}

func stringVersionToNumber(ver string) int {
	numStr := strings.ReplaceAll(ver, ".", "")
	numStr = strings.ReplaceAll(numStr, "rc", "")

	num, err := strconv.Atoi(numStr)
	if err != nil {
		num = 0
	}

	return num
}