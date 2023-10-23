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

	tree.Delete(3)
	assert.Equal(t, tree.Visit(), "2B {1B 4R {6B {5R 7R}}}")
	tree.Delete(5)
	assert.Equal(t, tree.Visit(), "2B {1B 4R {6B {7R}}}")
	tree.Delete(1)
	assert.Equal(t, tree.Visit(), "2B {4R {6B {7R}}}")
	tree.Delete(6)
	assert.Equal(t, tree.Visit(), "2B {4R {7B}}")
	tree.Delete(4)
	assert.Equal(t, tree.Visit(), "2B {7B}")
	tree.Delete(2)
	assert.Equal(t, tree.Visit(), "7B")
	tree.Delete(7)
	assert.Assert(t, tree.IsEmpty())

	// random order
	// tree.Insert(4)
	// assert.Equal(t, tree.Visit(), "4")
	// tree.Insert(2)
	// assert.Equal(t, tree.Visit(), "4 {2}")
	// tree.Insert(3)
	// assert.Equal(t, tree.Visit(), "4 {2 {3}}")
	// tree.Insert(6)
	// assert.Equal(t, tree.Visit(), "4 {2 {3} 6}")
	// tree.Insert(5)
	// assert.Equal(t, tree.Visit(), "4 {2 {3} 6 {5}}")
	// tree.Insert(1)
	// assert.Equal(t, tree.Visit(), "4 {2 {1 3} 6 {5}}")
	// tree.Insert(7)
	// assert.Equal(t, tree.Visit(), "4 {2 {1 3} 6 {5 7}}")

	// tree.Delete(2)
	// assert.Equal(t, tree.Visit(), "4 {3 {1} 6 {5 7}}")
	// tree.Delete(6)
	// assert.Equal(t, tree.Visit(), "4 {3 {1} 7 {5}}")
	// tree.Delete(1)
	// assert.Equal(t, tree.Visit(), "4 {3 7 {5}}")
	// tree.Delete(7)
	// assert.Equal(t, tree.Visit(), "4 {3 5}")
	// tree.Delete(4)
	// assert.Equal(t, tree.Visit(), "5 {3}")
	// tree.Delete(3)
	// assert.Equal(t, tree.Visit(), "5")
	// tree.Delete(5)
	// assert.Assert(t, tree.IsEmpty())
}
