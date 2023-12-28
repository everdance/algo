package algo_test

import (
	"algo"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAlgo_KMP(t *testing.T) {
	assert.DeepEqual(t, algo.KMP("abababab", "d"), []int{})
	assert.DeepEqual(t, algo.KMP("abababab", "ab"), []int{0, 2, 4, 6})
	assert.DeepEqual(t, algo.KMP("hel你好lo,你好", "你好"), []int{3, 8})
}
