package algo

import (
	"fmt"
	"strings"
)

type Color string

const (
	Red   Color = "R"
	Black Color = "B"
)

type RBNode struct {
	Key    int
	Color  Color
	Parent *RBNode
	Left   *RBNode
	Right  *RBNode
}

func (n *RBNode) InOrderVisit() string {
	if n != nil {
		children := fmt.Sprintf("%s %s", n.Left.InOrderVisit(),
			n.Right.InOrderVisit())
		children = strings.Trim(children, " ")

		if children == "" {
			return fmt.Sprintf("%v%v", n.Key, n.Color)
		}

		return fmt.Sprintf("%v%v {%s}", n.Key, n.Color, children)
	}

	return ""
}

func (n *RBNode) Search(k int) *RBNode {
	if n == nil {
		return nil
	}
	if k < n.Key {
		return n.Left.Search(k)
	}
	if k > n.Key {
		return n.Right.Search(k)
	}

	return n
}

func (n *RBNode) Insert(k int) *RBNode {
	if n == nil {
		panic("insert on nil node")
	}

	if k < n.Key {
		if n.Left == nil {
			n.Left = &RBNode{Key: k, Parent: n, Color: Red}
			return n.Left
		}
		return n.Left.Insert(k)
	}

	if k > n.Key {
		if n.Right == nil {
			n.Right = &RBNode{Key: k, Parent: n, Color: Red}
			return n.Right
		}
		return n.Right.Insert(k)
	}

	return n
}

func (n *RBNode) RotateRight(tree *RBTree) {
	// left child must exists
	child := n.Left
	if child == nil {
		panic("left child must exist")
	}

	n.Left = child.Right
	if child.Right != nil {
		child.Right.Parent = n
	}

	child.Parent = n.Parent
	n.Parent = child
	child.Right = n

	if child.Parent != nil {
		if child.Parent.Left == n {
			child.Parent.Left = child
		} else {
			child.Parent.Right = child
		}
	} else {
		tree.root = child
	}
}

func (n *RBNode) RotateLeft(tree *RBTree) {
	// right child must exists
	child := n.Right
	if child == nil {
		panic("right child must exist")
	}

	n.Right = child.Left
	if child.Left != nil {
		child.Left.Parent = n
	}

	child.Parent = n.Parent
	n.Parent = child
	child.Left = n

	if child.Parent != nil {
		if child.Parent.Left == n {
			child.Parent.Left = child
		} else {
			child.Parent.Right = child
		}
	} else {
		tree.root = child
	}
}

func FixInsert(tree *RBTree, n *RBNode) {
	// n is root node or n's parent is black, then fix is done
	if n.Parent == nil || n.Parent.Color == Black {
		return
	}

	// grandparent must exist and its color must be black
	grandpa := n.Parent.Parent
	uncle := grandpa.Right
	if n.Parent == grandpa.Right {
		uncle = grandpa.Left
	}

	if uncle == nil {
		grandpa.Color = Red
		n.Color = Red
		n.Parent.Color = Black
		if n.Parent == grandpa.Left {
			grandpa.RotateRight(tree)
		} else {
			grandpa.RotateLeft(tree)
		}
		return
	}

	if uncle.Color == Red {
		uncle.Color = Black
		n.Parent.Color = Black
		grandpa.Color = Red
		FixInsert(tree, grandpa)
		return
	}

	grandpa.Color = Red
	n.Parent.Color = Black

	if uncle == grandpa.Right {
		if n == n.Parent.Right {
			n.Parent.RotateLeft(tree)
		}
		grandpa.RotateRight(tree)
	} else {
		if n == n.Parent.Left {
			n.Parent.RotateRight(tree)
		}
		grandpa.RotateLeft(tree)
	}
}

type RBTree struct {
	root *RBNode
}

func (t *RBTree) IsEmpty() bool { return t.root == nil }

func (t *RBTree) Visit() string { return t.root.InOrderVisit() }

func (t *RBTree) Search(k int) *RBNode { return t.root.Search(k) }

func (t *RBTree) Insert(k int) {
	if t.root == nil {
		t.root = &RBNode{Key: k, Color: Black}
		return
	}

	n := t.root.Insert(k)
	FixInsert(t, n)
	t.root.Color = Black
}

func (t *RBTree) Delete(k int) {
}
