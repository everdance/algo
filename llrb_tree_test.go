package algo_test

import (
	"algo"
	"fmt"
	"math/rand"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAlgo_LLRBTree(t *testing.T) {
	tree := algo.LLRBTree{}
	assert.Assert(t, tree.IsEmpty())

	i := 0
	keys := []int{}
	for i < 100 {
		keys = append(keys, i)
		j := rand.Intn(i + 1)
		keys[i], keys[j] = keys[j], i
		i++
	}

	for _, k := range keys {
		tree.Insert(k)
		assert.Assert(t, tree.Search(k), "search failed")
		assert.Assert(t, tree.Check(), "insert violates tree")
	}

	fmt.Println(tree.Visit())

	for _, k := range keys {
		tree.Delete(k)
		fmt.Println(k, " <- ", tree.Visit())
		assert.Assert(t, !tree.Search(k))
		assert.Assert(t, tree.Check(), "delete violates tree")
	}

	assert.Assert(t, tree.IsEmpty())
}
