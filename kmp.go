package algo

// build the brack track table for match string
func bktrack(pat []rune) []int {
	bkt := make([]int, len(pat))
	bkt[0] = 0
	l, i := 0, 1
	for i < len(pat) {
		if pat[i] == pat[l] {
			l++
			bkt[i] = l
			i++
		} else {
			if l > 0 {
				l = bkt[l-1]
			} else {
				bkt[i] = 0
				i++
			}
		}
	}
	return bkt
}

func torune(s string) []rune {
	chars := []rune{}
	for _, v := range s {
		chars = append(chars, v)
	}

	return chars
}

func KMP(s, pattern string) []int {
	str, pat := torune(s), torune(pattern)
	bkt := bktrack(pat)
	pos, i, j := []int{}, 0, 0

	for j < len(str) {
		if str[j] == pat[i] {
			if i == len(pat)-1 {
				pos = append(pos, j-i)
				i = bkt[i]
			} else {
				i++
			}
			j++
		} else if i > 0 {
			i = bkt[i-1]
		} else {
			j++
		}
	}

	return pos
}
