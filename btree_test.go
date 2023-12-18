package algo_test

import (
	"algo"
	"fmt"
	"math/rand"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAlgo_BTree(t *testing.T) {
	tree := algo.NewBTree(10)
	assert.Assert(t, tree.IsEmpty())

	i := 0
	keys := []int{}
	for i < 20 {
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

	for _, k := range keys {
		fmt.Println(k, " ->> ")
		tree.Delete(k)
		tree.Print()
		assert.Assert(t, !tree.Search(k))
		assert.Assert(t, tree.Check(), "delete violates tree")
	}

	assert.Assert(t, tree.IsEmpty())
}
