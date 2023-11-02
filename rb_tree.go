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

func (n *RBnode) isleaf() bool {
	return n != nil && n.Left == nil && n.Right == nil
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
		tree.root = c
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

func (n *RBnode) transplant(m *RBnode, t *RBTree) {
	if n == nil || m == nil {
		panic("node is nil")
	}

	if n.Parent != nil {
		if n == n.Parent.Left {
			n.Parent.Left = m
		} else {
			n.Parent.Right = m
		}
	} else {
		t.root = m
	}

	m.Parent = n.Parent
	m.Left = n.Left
	m.Right = n.Right
	if m.Left != nil {
		m.Left.Parent = m
	}
	if m.Right != nil {
		m.Right.Parent = m
	}
	// remove n from links
	n.Left = nil
	n.Right = nil
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

func (n *RBnode) height() int {
	h := 0
	for n != nil {
		if n.Color == Black {
			h++
		}
		n = n.Left
	}
	return h
}

type RBTree struct {
	root *RBnode
}

func (t *RBTree) IsEmpty() bool { return t.root == nil }

func (t *RBTree) Visit() string { return t.root.preorder() }

func (t *RBTree) Search(k int) *RBnode { return t.root.search(k) }

func (t *RBTree) Height() int {
	return t.root.height()
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
	// the parent node we can hold on to in case the deleted node does not have children
	z := n.Parent
	if n.Left == nil {
		n.removeBy(n.Right, t)
		y = n.Right
	} else if n.Right == nil {
		n.removeBy(n.Left, t)
		y = n.Left
	} else { // successor must exists in n's right branch
		succ := n.successor()
		color = succ.Color
		y = succ.Right
		succ.Color = n.Color

		if succ == n.Right {
			n.removeBy(succ, t)
			succ.Left = n.Left
			if n.Left != nil {
				n.Left.Parent = succ
			}
			z = succ // n is replaced by succ directly
		} else {
			z = succ.Parent // succ is moved out to replace n
			succ.removeBy(succ.Right, t)
			n.transplant(succ, t)
		}
	}

	if color == Red {
		return
	}

	if y == nil {
		if z == nil {
			return
		}

		y = t.fixStartPoint(z)
	}

	t.FixDelete(y)
}

func (t *RBTree) fixStartPoint(z *RBnode) *RBnode {
	// we should start fix from deleted node's child y, but because y is nil,
	// we have to find a new start point based on deleted node's parent z
	// z right now has at either left or right child only, and
	// at most three levels of chidren, the removed child node is black
	//     z       z             z
	//    : \     : \           :  \
	//    y  s    y  s   or     y   s
	//      / \     / \
	//     a   b   a   b
	//    / \ / \
	//   i  j k  l
	var a, b, s, y *RBnode

	if z.Right != nil {
		s = z.Right
		a = z.Right.Left
		b = z.Right.Right
		if a != nil && b == nil {
			//   z          a
			//    \   ->   / \
			//     s      z   s
			//    /
			//   a
			s.rotateRight(t)
			z.rotateLeft(t)
			a.Color = z.Color
			z.Color = s.Color // black
			y = t.root
		} else if a == nil && b == nil {
			//   z
			//    \
			//     s
			s.Color = Red // s color is black as the removed child
			y = z
		} else if a == nil && b != nil {
			//   z              s
			//    \            / \
			//     s    ->    z   b
			//      \
			//       b
			z.rotateLeft(t)
			s.Color = z.Color // s color is black as the removed child
			b.Color = Black
			z.Color = Black
			y = t.root
		} else if a.isleaf() {
			//   z          s
			//    \   ->   / \
			//     s      z   b
			//    / \      \
			//   a   b      a
			s.Color = z.Color
			b.Color = Black
			a.Color = Red
			z.Color = Black
			z.rotateLeft(t)
			y = t.root
		} else {
			//   z                          s
			//    \             ->         / \
			//     s   (red)              z   b
			//    / \                      \   \
			//   a   b (black)              a   k
			//  / \ /                      / \
			// i  j k   (red)             i   j
			z.rotateLeft(t)
			s.Color = z.Color
			z.Color = Red
			return t.fixStartPoint(z)
		}
	} else if z.Left != nil {
		s = z.Left
		a = z.Left.Right
		b = z.Left.Left
		if a != nil && b == nil {
			s.rotateLeft(t)
			z.rotateRight(t)
			a.Color = z.Color
			z.Color = s.Color
			y = t.root
		} else if a == nil && b == nil {
			s.Color = Red
			y = z
		} else if a == nil && b != nil {
			z.rotateRight(t)
			s.Color = z.Color // s color is black as the removed child
			b.Color = Black
			z.Color = Black
			y = t.root
		} else if a.isleaf() {
			s.Color = z.Color
			a.Color = Red
			z.Color = Black
			b.Color = Black
			z.rotateRight(t)
			y = t.root
		} else {
			z.rotateRight(t)
			s.Color = z.Color
			z.Color = Red
			return t.fixStartPoint(z)
		}
	} else {
		// this case is impossible as if z is now leaf node
		// then the previously deleted child node must be red node
		//  we should have returned already
		y = z
	}

	return y
}

func (t *RBTree) FixDelete(n *RBnode) {
	for n.Parent != nil && n.Color == Black {
		sibling := n.Parent.Left
		if n == n.Parent.Left {
			sibling = n.Parent.Right
		}

		// as n has double black, sibling must have left and right child
		// otherwise the subtree black height is not equal
		if sibling.isleaf() {
			panic("sibling can't be leaf")
		}

		// rotate to get sibling color black
		if sibling.Color == Red {
			sibling.Color = n.Parent.Color
			n.Parent.Color = Red
			if sibling == n.Parent.Left {
				n.Parent.rotateRight(t)
				sibling = n.Parent.Left
			} else {
				n.Parent.rotateLeft(t)
				sibling = n.Parent.Right
			}
		}

		if sibling.Left.Color == Black && sibling.Right.Color == Black {
			sibling.Color = Red
			n = n.Parent
			continue
		}
		// one child of sibling must be red
		// rotate red to the same child branch direction as sibling to its parent
		if sibling == n.Parent.Left {
			if sibling.Left.Color == Black {
				sibling.Color = Red
				sibling.Right.Color = Black
				sibling.rotateLeft(t)
			}
			sibling = n.Parent.Left // new rotated sibiling node
			sibling.Color = n.Parent.Color
			sibling.Left.Color = Black
			n.Parent.Color = Black
			n.Parent.rotateRight(t)
		} else {
			if sibling.Right.Color == Black {
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
