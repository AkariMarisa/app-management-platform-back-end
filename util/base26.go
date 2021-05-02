package base26

const base26 = "abcdefghijklmnopqrstuvwxyz"

func Encode(num int) string {
	numStr := ""
	for num > 0 {
		leftover := num % 26
		numStr = string(base26[leftover]) + numStr
		num = num / 26
	}
	return numStr
}
