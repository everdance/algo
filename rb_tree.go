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
	Color  Color
	Key    int
	Parent *RBNode
	Left   *RBNode
	Right  *RBNode
}

func (n *RBNode) preorder() string {
	if n != nil {
		children := fmt.Sprintf("%s %s", n.Left.preorder(),
			n.Right.preorder())
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

func (n *RBNode) DropByChild(child *RBNode) {
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

func (n *RBNode) Successor() *RBNode {
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

func (n *RBNode) ReplaceBy(m *RBNode) {
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

func (tree *RBTree) FixInsert(n *RBNode) {
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
		tree.FixInsert(grandpa)
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

func (t *RBTree) Visit() string { return t.root.preorder() }

func (t *RBTree) Search(k int) *RBNode { return t.root.Search(k) }

func (t *RBTree) Insert(k int) {
	if t.root == nil {
		t.root = &RBNode{Key: k, Color: Black}
		return
	}

	n := t.root.Insert(k)
	t.FixInsert(n)
	t.root.Color = Black
}

func (t *RBTree) FixDelete(n *RBNode) {
	for n.Parent != nil && n.Color == Black {
		sibling := n.Parent.Left
		if n == n.Parent.Left {
			sibling = n.Parent.Right
		}

		if sibling == nil {
			n = n.Parent
			continue
		}

		// rotate to make silbing color be black
		if sibling.Color == Red {
			if sibling == n.Parent.Left && sibling.Right != nil {
				sibling.RotateLeft(t)
				sibling = n.Parent.Left
			} else if sibling == n.Parent.Right && sibling.Left != nil {
				sibling.RotateRight(t)
				sibling = n.Parent.Right
			}
		}

		// as n has double black, sibling must have at least one child
		// otherwise the subtree black height is not equal
		// both children are black if any one exists
		if (sibling.Left == nil || sibling.Left.Color == Black) &&
			(sibling.Right == nil || sibling.Right.Color == Black) {
			sibling.Color = Red
			n = n.Parent
			continue
		}
		// one child of sibling must be red
		// rotate red to the same child branch direction as sibling to its parent
		if sibling == n.Parent.Left {
			if sibling.Left == nil || sibling.Left.Color == Black {
				sibling.Color = Red
				sibling.Right.Color = Black
				sibling.RotateLeft(t)
			}
			sibling = n.Parent.Left // repoint sibling to new rotated node
			sibling.Color = n.Parent.Color
			sibling.Left.Color = Black
			n.Parent.Color = Black
			n.Parent.RotateRight(t)
		} else {
			if sibling.Right == nil || sibling.Right.Color == Black {
				sibling.Color = Red
				sibling.Left.Color = Black
				sibling.RotateRight(t)
			}
			sibling = n.Parent.Right
			sibling.Color = n.Parent.Color
			sibling.Right.Color = Black
			n.Parent.Color = Black
			n.Parent.RotateLeft(t)
		}

		n = t.root
	}

	n.Color = Black
}

func (t *RBTree) Delete(k int) {
	n := t.root.Search(k)
	if n == nil {
		return
	}

	y := n
	color := n.Color
	if n.Left == nil {
		n.DropByChild(n.Right)
		y = n.Right
	} else if n.Right == nil {
		n.DropByChild(n.Left)
		y = n.Left
	} else {
		succ := n.Successor()
		y = succ.Right
		color = succ.Color
		succ.DropByChild(succ.Right)
		n.ReplaceBy(succ)
		succ.Color = n.Color
	}

	if y == nil {
		y = n.Parent
	}

	if color == Red {
		return
	}

	t.FixDelete(y)
}
