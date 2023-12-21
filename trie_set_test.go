package algo_test

import (
	"algo"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAlgo_TrieSet(t *testing.T) {
	set := algo.TrieSet{}
	assert.Assert(t, set.IsEmpty())

	assert.NilError(t, set.Put("hello"))
	assert.Assert(t, set.Contains("hello"))
	assert.Error(t, set.Put("你好"), "only supports ascii chars 0-127")
	assert.NilError(t, set.Put("bbc"))
	assert.NilError(t, set.Put("mill"))
	assert.NilError(t, set.Put("million"))
	strs := set.KeysWithPrefix("mi")
	assert.DeepEqual(t, strs, []string{"mill", "million"})
}
