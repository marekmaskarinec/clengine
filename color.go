package cliw

import (
	"strconv"
	"strings"
)

func getColor(in string, bg bool) string {
	/* Replaces old color map with hex codes */
	color := map[string]string{
		"red":     "ff0000",
		"blnk":    "000000",
		"green":   "00ff00",
		"blue":    "0000ff",
		"white":   "ffffff",
		"cyan":    "00ffff",
		"magenta": "ff0000",
	}
	if _, ok := color[in]; ok {
		in = color[in]
	}

	/* Sets color to black if invalid */
	if len(in) < 6 {
		in = "000000"
	}

	/* Splits hexcode to three numbers */
	tmpNums := strings.Split(in, "")
	numsIn := []string{strings.Join(tmpNums[0:2], ""), strings.Join(tmpNums[2:4], ""), strings.Join(tmpNums[4:6], "")}

	/* Converts numbers to base10 */
	nums := []string{}
	var num int64
	for i := range numsIn {
		num, _ = strconv.ParseInt(numsIn[i], 16, 0)
		nums = append(nums, strconv.Itoa(int(num)))
	}
	/* Generates ANSI escape code */
	if !bg {
		return "\033[38;2;" + nums[0] + ";" + nums[1] + ";" + nums[2] + "m"
	} else {
		return "\033[48;2;" + nums[0] + ";" + nums[1] + ";" + nums[2] + "m"
	}
}
