package algo_test

import (
	"algo"
	"math/rand"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAlgo_RBTree(t *testing.T) {
	tree := algo.RBTree{}
	assert.Assert(t, tree.IsEmpty())

	tree.Insert(1)
	assert.Assert(t, !tree.IsEmpty())
	assert.Assert(t, tree.Search(1).Key == 1)

	tree.Delete(1)
	assert.Assert(t, tree.IsEmpty())
	assert.Assert(t, tree.Check())

	i := 0
	keys := []int{}
	for i < 10000 {
		keys = append(keys, i)
		j := rand.Intn(i + 1)
		keys[i], keys[j] = keys[j], i
		i++
	}

	for _, k := range keys {
		tree.Insert(k)
		assert.Assert(t, tree.Search(k).Key == k)
		assert.Assert(t, tree.Check(), "insert violates tree")
	}

	for _, k := range keys {
		tree.Delete(k)
		assert.Assert(t, tree.Search(k) == nil)
		assert.Assert(t, tree.Check(), "delete violates tree")
	}

	assert.Assert(t, tree.IsEmpty())
}
