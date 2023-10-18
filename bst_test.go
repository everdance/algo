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

	// asc order
	tree.Insert(1)
	assert.Equal(t, tree.Visit(), "1")
	tree.Insert(2)
	assert.Equal(t, tree.Visit(), "1 {2}")
	tree.Insert(3)
	assert.Equal(t, tree.Visit(), "1 {2 {3}}")
	tree.Insert(4)
	assert.Equal(t, tree.Visit(), "1 {2 {3 {4}}}")
	tree.Insert(5)
	assert.Equal(t, tree.Visit(), "1 {2 {3 {4 {5}}}}")
	tree.Insert(6)
	assert.Equal(t, tree.Visit(), "1 {2 {3 {4 {5 {6}}}}}")
	tree.Insert(7)
	assert.Equal(t, tree.Visit(), "1 {2 {3 {4 {5 {6 {7}}}}}}")

	tree.Delete(3)
	assert.Equal(t, tree.Visit(), "1 {2 {4 {5 {6 {7}}}}}")
	tree.Delete(5)
	assert.Equal(t, tree.Visit(), "1 {2 {4 {6 {7}}}}")
	tree.Delete(1)
	assert.Equal(t, tree.Visit(), "2 {4 {6 {7}}}")
	tree.Delete(6)
	assert.Equal(t, tree.Visit(), "2 {4 {7}}")
	tree.Delete(4)
	assert.Equal(t, tree.Visit(), "2 {7}")
	tree.Delete(2)
	assert.Equal(t, tree.Visit(), "7")
	tree.Delete(7)
	assert.Assert(t, tree.IsEmpty())

	// random order
	tree.Insert(4)
	assert.Equal(t, tree.Visit(), "4")
	tree.Insert(2)
	assert.Equal(t, tree.Visit(), "4 {2}")
	tree.Insert(3)
	assert.Equal(t, tree.Visit(), "4 {2 {3}}")
	tree.Insert(6)
	assert.Equal(t, tree.Visit(), "4 {2 {3} 6}")
	tree.Insert(5)
	assert.Equal(t, tree.Visit(), "4 {2 {3} 6 {5}}")
	tree.Insert(1)
	assert.Equal(t, tree.Visit(), "4 {2 {1 3} 6 {5}}")
	tree.Insert(7)
	assert.Equal(t, tree.Visit(), "4 {2 {1 3} 6 {5 7}}")

	tree.Delete(2)
	assert.Equal(t, tree.Visit(), "4 {3 {1} 6 {5 7}}")
	tree.Delete(6)
	assert.Equal(t, tree.Visit(), "4 {3 {1} 7 {5}}")
	tree.Delete(1)
	assert.Equal(t, tree.Visit(), "4 {3 7 {5}}")
	tree.Delete(7)
	assert.Equal(t, tree.Visit(), "4 {3 5}")
	tree.Delete(4)
	assert.Equal(t, tree.Visit(), "5 {3}")
	tree.Delete(3)
	assert.Equal(t, tree.Visit(), "5")
	tree.Delete(5)
	assert.Assert(t, tree.IsEmpty())
}
