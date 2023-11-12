package algo_test

import (
	"algo"
	"math/rand"
	"testing"

	"golang.org/x/exp/slices"
	"gotest.tools/v3/assert"
)

func TestAlgo_23Tree(t *testing.T) {
	tree := &algo.Tree23{}

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
		n := tree.Search(k)
		assert.Assert(t, slices.Contains(n.Keys, k), "search failed")
		assert.Assert(t, tree.Check(), "insert violates tree")
	}

	for _, k := range keys {
		tree.Delete(k)
		assert.Assert(t, tree.Check(), "delete violates tree")
	}

	assert.Assert(t, tree.IsEmpty())
}
