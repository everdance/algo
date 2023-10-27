package algo

import (
	"fmt"
	"math"
	"strings"
)

type Color string

const (
	Red   Color = "R"
	Black Color = "B"
)

type RBnode struct {
	Color  Color
	Key    int
	Parent *RBnode
	Left   *RBnode
	Right  *RBnode
}

func (n *RBnode) preorder() string {
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

func (n *RBnode) search(k int) *RBnode {
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

func (n *RBnode) insert(k int) *RBnode {
	if n == nil {
		panic("insert on nil node")
	}

	if k < n.Key {
		if n.Left == nil {
			n.Left = &RBnode{Key: k, Parent: n, Color: Red}
			return n.Left
		}
		return n.Left.insert(k)
	}

	if k > n.Key {
		if n.Right == nil {
			n.Right = &RBnode{Key: k, Parent: n, Color: Red}
			return n.Right
		}
		return n.Right.insert(k)
	}

	return nil
}

func (n *RBnode) rotateRight(tree *RBTree) {
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

func (n *RBnode) rotateLeft(tree *RBTree) {
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

// directly remove n by its only child c
func (n *RBnode) removeBy(c *RBnode, tree *RBTree) {
	if n == nil {
		panic("node is nil")
	}

	if c != nil {
		c.Parent = n.Parent
	}

	if n.Parent == nil {
		return
	}

	if n == n.Parent.Left {
		n.Parent.Left = c
	} else {
		n.Parent.Right = c
	}
}

func (n *RBnode) successor() *RBnode {
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

func (n *RBnode) transplant(m *RBnode) {
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

func (n *RBnode) isBST(min, max int) bool {
	if n == nil {
		return true
	}

	if n.Key < min || n.Key > max {
		return false
	}

	return n.Left.isBST(min, n.Key) && n.Right.isBST(n.Key, max)
}

func (n *RBnode) isRBTree(h int, pColor Color) bool {
	if n == nil {
		return h == 0
	}

	if n.Color == Black {
		h--
	} else if pColor == Red {
		return false
	}

	return n.Left.isRBTree(h, n.Color) && n.Right.isRBTree(h, n.Color)
}

type RBTree struct {
	root *RBnode
}

func (t *RBTree) IsEmpty() bool { return t.root == nil }

func (t *RBTree) Visit() string { return t.root.preorder() }

func (t *RBTree) Search(k int) *RBnode { return t.root.search(k) }

func (t *RBTree) Height() int {
	h := 0
	x := t.root
	for x != nil {
		if x.Color == Black {
			h++
		}
		x = x.Left
	}
	return h
}

func (t *RBTree) Check() bool {
	return t.root.isRBTree(t.Height(), Red) &&
		t.root.isBST(math.MinInt, math.MaxInt)
}

func (t *RBTree) Insert(k int) {
	if t.root == nil {
		t.root = &RBnode{Key: k, Color: Black}
		return
	}

	n := t.root.insert(k)
	if n != nil {
		t.FixInsert(n)
	}
	t.root.Color = Black
}

func (tree *RBTree) FixInsert(n *RBnode) {
	// n is root node or n's parent is black, then fix is done
	if n.Parent == nil || n.Parent.Color == Black {
		return
	}

	// because parent exists and has Red color
	// grandparent must exist and its color must be black
	grandpa := n.Parent.Parent
	uncle := grandpa.Right
	if n.Parent == grandpa.Right {
		uncle = grandpa.Left
	}

	if uncle == nil { // n must be a single child of its parent
		grandpa.Color = Red
		n.Parent.Color = Black
		if n.Parent == grandpa.Left {
			if n == n.Parent.Right {
				n.Parent.Color = Red
				n.Color = Black
				n.Parent.rotateLeft(tree)
			}
			grandpa.rotateRight(tree)
		} else {
			if n == n.Parent.Left {
				n.Parent.Color = Red
				n.Color = Black
				n.Parent.rotateRight(tree)
			}
			grandpa.rotateLeft(tree)
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
			n.Parent.Color = Red
			n.Color = Black
			n.Parent.rotateLeft(tree)
		}
		grandpa.rotateRight(tree)
	} else {
		if n == n.Parent.Left {
			n.Parent.Color = Red
			n.Color = Black
			n.Parent.rotateRight(tree)
		}
		grandpa.rotateLeft(tree)
	}
}

func (t *RBTree) Delete(k int) {
	n := t.root.search(k)
	if n == nil {
		return
	}

	y := n // first bottom node we need to start fix on
	color := n.Color
	succ := n.successor()
	// the parent node we can hold on to in case the deleted node does not have children
	z := n.Parent
	if n.Left == nil {
		n.removeBy(n.Right, t)
		y = n.Right
		succ = y
	} else if n.Right == nil {
		n.removeBy(n.Left, t)
		y = n.Left
		succ = y
	} else { // successor must exists in n's right branch
		color = succ.Color
		y = succ.Right
		succ.Color = n.Color

		if succ == n.Right {
			n.removeBy(succ, t)
			succ.Left = n.Left
			if n.Left != nil {
				n.Left.Parent = succ
			}
			z = succ // n is replaced succ directly
		} else {
			z = succ.Parent // succ is moved out to replace n
			succ.removeBy(succ.Right, t)
			n.transplant(succ)
		}
	}

	if t.root == n {
		t.root = succ
	}

	if y == nil {
		if z == nil {
			return
		}

		y, color = t.fixStartPoint(z, color)
	}

	if color == Red {
		return
	}

	t.FixDelete(y)
}

func (t *RBTree) fixStartPoint(z *RBnode, c Color) (*RBnode, Color) {
	// z has at most left or right child, and only at most one level of grand chidren
	// the deleted child is Black
	//     z             z
	//    : \           :  \
	//    x  s   or     x   s
	//      / \
	//     a   b
	var a, b, s, y *RBnode

	if z.Left != nil {
		s = z.Left
		a = z.Left.Right
		b = z.Left.Left
		if b == nil && a != nil {
			z.Left.rotateLeft(t)
		}
		if b == nil && a != nil {
			z.Right.rotateLeft(t)
			c = a.Color
			a.Color = z.Color
			z.Color = s.Color
			y = a
		} else if a == nil {
			s.Color = z.Color
			z.Color = Red
			y = s
		} else {
			s.Color = z.Color
			a.Color = Red
			z.Color = Black
			b.Color = Black
			y = t.root
		}
		z.rotateRight(t)
	} else if z.Right != nil {
		s = z.Right
		a = z.Right.Left
		b = z.Right.Right
		if b == nil && a != nil {
			//   z          a
			//    \   ->   / \
			//     s      z   s
			//    /
			//   a
			z.Right.rotateRight(t)
			c = a.Color
			a.Color = z.Color
			z.Color = s.Color
			y = a
		} else if a == nil {
			//   z          s
			//    \   ->   /
			//     s      z
			s.Color = z.Color
			z.Color = Red
			y = s
		} else {
			//   z          s
			//    \   ->   / \
			//     s      z   b
			//    / \      \
			//   a   b      a
			s.Color = z.Color
			a.Color = Red
			z.Color = Black
			b.Color = Black
			y = t.root
		}
		z.rotateLeft(t)
	} else {
		y = z
	}

	return y, c
}

func (t *RBTree) FixDelete(n *RBnode) {
	for n.Parent != nil && n.Color == Black {
		sibling := n.Parent.Left
		if n == n.Parent.Left {
			sibling = n.Parent.Right
		}

		// sibling must exists as n is black
		// rotate to make silbing color be black
		if sibling.Color == Red {
			n.Parent.Color = Red
			sibling.Color = Black
			if sibling == n.Parent.Left {
				n.Parent.rotateRight(t)
				sibling = n.Parent.Left
			} else {
				n.Parent.rotateLeft(t)
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
				sibling.rotateLeft(t)
			}
			sibling = n.Parent.Left // repoint sibling to new rotated node
			sibling.Color = n.Parent.Color
			sibling.Left.Color = Black
			n.Parent.Color = Black
			n.Parent.rotateRight(t)
		} else {
			if sibling.Right == nil || sibling.Right.Color == Black {
				sibling.Color = Red
				sibling.Left.Color = Black
				sibling.rotateRight(t)
			}
			sibling = n.Parent.Right
			sibling.Color = n.Parent.Color
			sibling.Right.Color = Black
			n.Parent.Color = Black
			n.Parent.rotateLeft(t)
		}

		n = t.root
	}

	n.Color = Black
}
