package algo

import (
	"fmt"
	"math"
	"strings"
)

type Node struct {
	Key    int
	Parent *Node
	Left   *Node
	Right  *Node
}

func (n *Node) preOrder() string {
	if n != nil {
		children := fmt.Sprintf("%s %s", n.Left.preOrder(),
			n.Right.preOrder())
		children = strings.Trim(children, " ")

		if children == "" {
			return fmt.Sprintf("%d", n.Key)
		}

		return fmt.Sprintf("%d {%s}", n.Key, children)
	}

	return ""
}

func (n *Node) isBST(min, max int) bool {
	if n == nil {
		return true
	}

	return n.Key > min && n.Key < max &&
		n.Left.isBST(min, n.Key) && n.Right.isBST(n.Key, max)
}

func (n *Node) search(k int) *Node {
	if n == nil {
		return nil
	}
	if k < n.Key {
		return n.Left.search(k)
	}
	if k > n.Key {
		return n.Right.search(k)
	}

	return n
}

func (n *Node) successor() *Node {
	if n == nil {
		return nil
	}
	if n.Right != nil {
		min := n.Right
		for min.Left != nil {
			min = min.Left
		}
		return min
	}

	y := n.Parent
	for y != nil && n == y.Right {
		n = y
		y = n.Parent
	}

	return y
}

func (n *Node) insert(k int) {
	if n == nil {
		panic("insert on nil node")
	}

	if k < n.Key {
		if n.Left == nil {
			n.Left = &Node{Key: k, Parent: n}
		} else {
			n.Left.insert(k)
		}
	}

	if k > n.Key {
		if n.Right == nil {
			n.Right = &Node{Key: k, Parent: n}
		} else {
			n.Right.insert(k)
		}
	}
}

// node has only one child, to remove it
// just move up the child to its position
func (n *Node) replaceByChild(child *Node) {
	if n == nil {
		panic("node is nil")
	}

	if n.Parent != nil {
		if n == n.Parent.Left {
			n.Parent.Left = child
		} else {
			n.Parent.Right = child
		}
	}

	if child != nil {
		child.Parent = n.Parent
	}
}

// use successor node m from the tree to
// replace node n position
func (n *Node) transplant(m *Node) {
	if n == nil {
		panic("node is nil")
	}

	if n.Parent != nil {
		if n == n.Parent.Left {
			n.Parent.Left = m
		} else {
			n.Parent.Right = m
		}
	}

	if m != nil {
		m.Parent = n.Parent
		m.Left = n.Left
		m.Right = n.Right
		if m.Left != nil {
			m.Left.Parent = m
		}
		if m.Right != nil {
			m.Right.Parent = m
		}
	}
}

type BST struct {
	root *Node
}

func (t *BST) IsEmpty() bool { return t.root == nil }

func (t *BST) Check() bool { return t.root.isBST(math.MinInt, math.MaxInt) }

func (t *BST) Visit() string { return t.root.preOrder() }

func (t *BST) Search(k int) *Node { return t.root.search(k) }

func (t *BST) Insert(k int) {
	if t.root == nil {
		t.root = &Node{Key: k}
		return
	}

	t.root.insert(k)
}

func (t *BST) Delete(k int) {
	n := t.root.search(k)
	if n == nil {
		return
	}

	var succ *Node

	if n.Right == nil {
		n.replaceByChild(n.Left)
		succ = n.Left
	} else if n.Left == nil {
		n.replaceByChild(n.Right)
		succ = n.Right
	} else {
		// left child of smallest successor node must be nil
		// otherwise its left child is smaller
		succ = n.successor()
		succ.replaceByChild(succ.Right)
		n.transplant(succ)
	}

	if n == t.root {
		t.root = succ
	}
}
