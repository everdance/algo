package algo

import (
	"fmt"
	"math"
	"strings"
)

// this is an implementation without using null node to teriminate leaf nodes
// the complication of dealing with edge cases with leaf node are huge
// compare this to the llb tree implementation

type Color string
type Direction bool

const (
	Red   Color     = "R"
	Black Color     = "B"
	Left  Direction = true
	Right Direction = false
)

type RBnode struct {
	color Color
	Key   int
	p     *RBnode
	l     *RBnode
	r     *RBnode
}

func (n *RBnode) preorder() string {
	if n != nil {
		children := fmt.Sprintf("%s %s", n.l.preorder(), n.r.preorder())
		children = strings.Trim(children, " ")

		if children == "" {
			return fmt.Sprintf("%v%v", n.Key, n.color)
		}

		return fmt.Sprintf("%v%v {%s}", n.Key, n.color, children)
	}

	return ""
}

func (n *RBnode) isleaf() bool {
	return n != nil && n.l == nil && n.r == nil
}

func (n *RBnode) search(k int) *RBnode {
	if n == nil {
		return nil
	}
	if k < n.Key {
		return n.l.search(k)
	}
	if k > n.Key {
		return n.r.search(k)
	}

	return n
}

func (n *RBnode) insert(k int) *RBnode {
	if n == nil {
		panic("insert on nil node")
	}

	if k < n.Key {
		if n.l == nil {
			n.l = &RBnode{Key: k, p: n, color: Red}
			return n.l
		}
		return n.l.insert(k)
	}

	if k > n.Key {
		if n.r == nil {
			n.r = &RBnode{Key: k, p: n, color: Red}
			return n.r
		}
		return n.r.insert(k)
	}
	// key exists do thing
	return nil
}

func (n *RBnode) rotate(dir Direction, tree *RBTree) {
	child := n.getchild(!dir)
	if child == nil {
		panic("child must exist")
	}

	if dir == Right {
		n.l = child.r
		if n.l != nil {
			n.l.p = n
		}
		child.r = n
	} else {
		n.r = child.l
		if n.r != nil {
			n.r.p = n
		}
		child.l = n
	}

	child.p = n.p
	n.p = child

	if child.p != nil {
		if child.p.l == n {
			child.p.l = child
		} else {
			child.p.r = child
		}
	} else {
		tree.root = child
	}
}

func (n *RBnode) getchild(d Direction) *RBnode {
	if n == nil {
		return nil
	}
	if d == Left {
		return n.l
	}
	return n.r
}

// directly remove n by its child c, c can be nil
func (n *RBnode) replace(d Direction, tree *RBTree) {
	if n == nil {
		panic("node is nil")
	}

	c := n.getchild(d)
	if c != nil {
		c.p = n.p
	}

	if n.p == nil {
		tree.root = c
		return
	}

	if n == n.p.l {
		n.p.l = c
	} else {
		n.p.r = c
	}
}

func (n *RBnode) successor() *RBnode {
	if n == nil {
		return nil
	}

	if n.r != nil {
		min := n.r
		for min.l != nil {
			min = min.l
		}
		return min
	}

	y := n.p
	for y != nil && n == y.r {
		n = y
		y = n.p
	}

	return y
}

func (n *RBnode) transplant(m *RBnode, t *RBTree) {
	if n == nil || m == nil {
		panic("node is nil")
	}

	if n.p != nil {
		if n == n.p.l {
			n.p.l = m
		} else {
			n.p.r = m
		}
	} else {
		t.root = m
	}

	m.p = n.p
	m.l = n.l
	m.r = n.r
	if m.l != nil {
		m.l.p = m
	}
	if m.r != nil {
		m.r.p = m
	}
	// remove n from tree
	n.l = nil
	n.r = nil
	n.p = nil
}

func (n *RBnode) isBST(min, max int) bool {
	if n == nil {
		return true
	}

	if n.Key < min || n.Key > max {
		return false
	}

	return n.l.isBST(min, n.Key) && n.r.isBST(n.Key, max)
}

func (n *RBnode) isRBTree(h int, pColor Color) bool {
	if n == nil {
		return h == 0
	}

	if n.color == Black {
		h--
	} else if pColor == Red {
		return false
	}

	return n.l.isRBTree(h, n.color) && n.r.isRBTree(h, n.color)
}

func (n *RBnode) height() int {
	h := 0
	for n != nil {
		if n.color == Black {
			h++
		}
		n = n.l
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
		t.root = &RBnode{Key: k, color: Black}
		return
	}

	n := t.root.insert(k)
	if n != nil {
		t.fixInsert(n)
	}
	t.root.color = Black
}

func (tree *RBTree) fixInsert(n *RBnode) {
	// n is root node or n's parent is black, then fix is done
	if n.p == nil || n.p.color == Black {
		return
	}

	// because parent exists and has Red color
	// grandparent must exist and its color must be black
	grandpa := n.p.p
	uncle := grandpa.r
	if n.p == grandpa.r {
		uncle = grandpa.l
	}

	if uncle == nil { // n must be a single child of its parent
		grandpa.color = Red
		n.p.color = Black
		if n.p == grandpa.l {
			if n == n.p.r {
				n.p.color = Red
				n.color = Black
				n.p.rotate(Left, tree)
			}
			grandpa.rotate(Right, tree)
		} else {
			if n == n.p.l {
				n.p.color = Red
				n.color = Black
				n.p.rotate(Right, tree)
			}
			grandpa.rotate(Left, tree)
		}
		return
	}

	if uncle.color == Red {
		uncle.color = Black
		n.p.color = Black
		grandpa.color = Red
		tree.fixInsert(grandpa)
		return
	}

	grandpa.color = Red
	n.p.color = Black

	if uncle == grandpa.r {
		if n == n.p.r {
			n.p.color = Red
			n.color = Black
			n.p.rotate(Left, tree)
		}
		grandpa.rotate(Right, tree)
	} else {
		if n == n.p.l {
			n.p.color = Red
			n.color = Black
			n.p.rotate(Right, tree)
		}
		grandpa.rotate(Left, tree)
	}
}

func (t *RBTree) Delete(k int) {
	n := t.root.search(k)
	if n == nil {
		return
	}

	y := n // first bottom node we need to start fix on
	color := n.color
	// the parent node we can hold on to in case the deleted node does not have children
	z := n.p
	if n.l == nil {
		n.replace(Right, t)
		y = n.r
	} else if n.r == nil {
		n.replace(Left, t)
		y = n.l
	} else { // successor must exists in n's right branch
		succ := n.successor()
		color = succ.color
		y = succ.r
		succ.color = n.color

		if succ == n.r {
			n.replace(Right, t)
			succ.l = n.l
			if n.l != nil {
				n.l.p = succ
			}
			z = succ // n is replaced by succ directly
		} else {
			z = succ.p // succ is moved out to replace n
			succ.replace(Right, t)
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

		y = t.fixNil(z)
	}

	t.fixBlack(y)
}

// when y is nil, we have to find a new start node based on deleted node's parent z
// z right now has at either left or right child only, most three layers of children
//
//	  z            z             z
//	 : \          : \           :  \
//	 y  s    OR   y  s    OR    y   s
//	   / \          / \
//	  a   b        a   b
//	 / \ / \
//	i  j k  l
func (t *RBTree) fixNil(z *RBnode) *RBnode {
	// this is impossible as if z is now leaf node
	// then the previously deleted child node must be red node
	//  we should have returned already
	if z.isleaf() || (z.l != nil && z.r != nil) {
		panic("z can only has one child")
	}

	var a, b, s, y *RBnode
	var dir Direction

	if z.r != nil {
		s = z.r
		a = z.r.l
		b = z.r.r
		dir = Left
	} else {
		s = z.l
		a = z.l.r
		b = z.l.l
		dir = Right
	}

	if a != nil && b == nil {
		//   z          a
		//    \   ->   / \
		//     s      z   s
		//    /
		//   a
		s.rotate(!dir, t)
		z.rotate(dir, t)
		a.color = z.color
		z.color = s.color // black
		y = t.root
	} else if a == nil && b == nil {
		//   z
		//    \
		//     s
		s.color = Red // s color is black as the removed child
		y = z
	} else if a == nil && b != nil {
		//   z              s
		//    \            / \
		//     s    ->    z   b
		//      \
		//       b
		z.rotate(dir, t)
		s.color = z.color // s color is black as the removed child
		b.color = Black
		z.color = Black
		y = t.root
	} else if a.isleaf() {
		//   z          s
		//    \   ->   / \
		//     s      z   b
		//    / \      \
		//   a   b      a
		s.color = z.color
		a.color = Red
		b.color = Black
		z.color = Black
		z.rotate(dir, t)
		y = t.root
	} else {
		//   z                          s
		//    \             ->         / \
		//     s   (red)              z   b
		//    / \                      \   \
		//   a   b (black)              a   k
		//  / \ /                      / \
		// i  j k   (red)             i   j
		z.rotate(dir, t)
		s.color = z.color
		z.color = Red
		return t.fixNil(z)
	}

	return y
}

// fix n with extra black carried on it
func (t *RBTree) fixBlack(n *RBnode) {
	for n.p != nil && n.color == Black {
		sibling := n.p.l
		if n == n.p.l {
			sibling = n.p.r
		}

		// as n has double black, sibling must have left and right child
		// otherwise the subtree black height is not equal
		if sibling.isleaf() {
			panic("sibling can't be leaf")
		}

		// rotate to get sibling color black
		if sibling.color == Red {
			sibling.color = n.p.color
			n.p.color = Red
			if sibling == n.p.l {
				n.p.rotate(Right, t)
				sibling = n.p.l
			} else {
				n.p.rotate(Left, t)
				sibling = n.p.r
			}
		}

		if sibling.l.color == Black && sibling.r.color == Black {
			sibling.color = Red
			n = n.p
			continue
		}
		// one child of sibling must be red
		// rotate red to the same child branch direction as sibling to its parent
		if sibling == n.p.l {
			if sibling.l.color == Black {
				sibling.color = Red
				sibling.r.color = Black
				sibling.rotate(Left, t)
			}
			sibling = n.p.l // new rotated sibiling node
			sibling.color = n.p.color
			sibling.l.color = Black
			n.p.color = Black
			n.p.rotate(Right, t)
		} else {
			if sibling.r.color == Black {
				sibling.color = Red
				sibling.l.color = Black
				sibling.rotate(Right, t)
			}
			sibling = n.p.r
			sibling.color = n.p.color
			sibling.r.color = Black
			n.p.color = Black
			n.p.rotate(Left, t)
		}

		n = t.root
	}

	n.color = Black
}
