package algo_test

import (
	"algo"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAlgo_BST(t *testing.T) {
	tree := algo.BST{}
	assert.Assert(t, tree.IsEmpty())

	tree.Insert(1)
	assert.Assert(t, tree.IsEmpty() == false)
	assert.Assert(t, tree.Search(1).Key == 1)

	tree.Insert(1)
	tree.Insert(2)
	tree.Insert(3)
	tree.Insert(4)
	tree.Insert(5)
	tree.Insert(6)
	tree.Insert(7)

	assert.Equal(t, tree.Visit(), "1 2 3 4 5 6 7")
}
