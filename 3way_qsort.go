package algo

import (
	"math/rand"
)

// 1 2 4 6 6 6 7 9 0 2 8 6
//
// - - - i     k   j     ^
func Qsort3way(s []int) {
	l := len(s)
	if l <= 1 {
		return
	}
	n := rand.Intn(l)
	s[n], s[0] = s[0], s[n]
	pivot := s[0]
	i, k, j := 0, 1, 1
	for j < l {
		if s[j] > pivot {
			j++
			continue
		}
		if s[j] < pivot {
			s[i], s[j] = s[j], s[i]
			i++
		}
		s[j], s[k] = s[k], s[j]
		j++
		k++
	}
	Qsort3way(s[:i])
	Qsort3way(s[k:])
}
