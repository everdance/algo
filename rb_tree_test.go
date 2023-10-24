package algo_test

import (
	"algo"
	"testing"

	"gotest.tools/v3/assert"
)

func TestAlgo_RBTree(t *testing.T) {
	tree := algo.RBTree{}
	assert.Assert(t, tree.IsEmpty())

	tree.Insert(1)
	assert.Assert(t, tree.IsEmpty() == false)
	assert.Assert(t, tree.Search(1).Key == 1)

	// asc order
	tree.Insert(1)
	assert.Equal(t, tree.Visit(), "1B")
	tree.Insert(2)
	assert.Equal(t, tree.Visit(), "1B {2R}")
	tree.Insert(3)
	assert.Equal(t, tree.Visit(), "2B {1R 3R}")
	tree.Insert(4)
	assert.Equal(t, tree.Visit(), "2B {1B 3B {4R}}")
	tree.Insert(5)
	assert.Equal(t, tree.Visit(), "2B {1B 4B {3R 5R}}")
	tree.Insert(6)
	assert.Equal(t, tree.Visit(), "2B {1B 4R {3B 5B {6R}}}")
	tree.Insert(7)
	assert.Equal(t, tree.Visit(), "2B {1B 4R {3B 6B {5R 7R}}}")
	assert.Assert(t, tree.Check())

	tree.Delete(3)
	assert.Equal(t, tree.Visit(), "2B {1B 4R {6B {5R 7R}}}")
	tree.Delete(5)
	assert.Equal(t, tree.Visit(), "2B {1B 4R {6B {7R}}}")
	tree.Delete(1)
	assert.Equal(t, tree.Visit(), "2B {4R {6B {7R}}}")

	assert.Assert(t, tree.Check())

	tree.Delete(6)
	assert.Equal(t, tree.Visit(), "2B {4R {7B}}")
	tree.Delete(4)
	assert.Equal(t, tree.Visit(), "2B {7B}")
	tree.Delete(2)
	assert.Equal(t, tree.Visit(), "7B")
	tree.Delete(7)

	assert.Assert(t, tree.IsEmpty())
	assert.Assert(t, tree.Check())
}
