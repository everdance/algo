package algo_test

import (
	"algo"
	"math/rand"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAlgo_BST(t *testing.T) {
	tree := algo.BST{}
	assert.Assert(t, tree.IsEmpty())

	i := 0
	keys := []int{}
	for i < 1000 {
		k := rand.Intn(1000)
		keys = append(keys, k)
		tree.Insert(k)
		assert.Assert(t, tree.Search(k).Key == k)
		assert.Assert(t, tree.Check())
		i++
	}

	for _, k := range keys {
		tree.Delete(k)
		assert.Assert(t, tree.Search(k) == nil)
		assert.Assert(t, tree.Check())
	}

	assert.Assert(t, tree.IsEmpty())
}
