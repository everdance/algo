package algo_test

import (
	"algo"
	"math/rand"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAlgo_3wayQsort(t *testing.T) {

	i := 0
	s := []int{}
	for i < 10000 {
		s = append(s, i)
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], i
		i++
	}
	algo.Qsort3way(s)
	for i := 0; i < len(s)-1; i++ {
		assert.Assert(t, s[i] <= s[i+1])
	}
}
