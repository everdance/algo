package algo

import (
	"fmt"
	"math"
	"strings"
)

// https://algs4.cs.princeton.edu/33balanced
// implementation uses null node to teriminate all leaf nodes which
// greatly reduce the complexity of dealing with edge cases with leaf node

type lrbNode struct {
	k int
	c Color
	l *lrbNode
	r *lrbNode
	p *lrbNode
}

var nullNode = &lrbNode{
	c: Black,
}

func (n *lrbNode) preorder() string {
	if n == nullNode || n == nil {
		return ""
	}

	children := fmt.Sprintf("%s %s", n.l.preorder(), n.r.preorder())
	children = strings.Trim(children, " ")

	if children == "" {
		return fmt.Sprintf("%v%v", n.k, n.c)
	}

	return fmt.Sprintf("%v%v {%s}", n.k, n.c, children)
}

func (n *lrbNode) search(k int) *lrbNode {
	if n == nil || n == nullNode {
		return nil
	}

	if k < n.k {
		return n.l.search(k)
	} else if k > n.k {
		return n.r.search(k)
	} else {
		return n
	}
}

func (n *lrbNode) insert(k int) *lrbNode {
	if n == nil {
		panic("insert on nil node")
	}

	if k < n.k {
		if n.l == nullNode {
			n.l = &lrbNode{k: k, p: n, c: Red, l: nullNode, r: nullNode}
			return n.l
		}
		return n.l.insert(k)
	}

	if k > n.k {
		if n.r == nullNode {
			n.r = &lrbNode{k: k, p: n, c: Red, l: nullNode, r: nullNode}
			return n.r
		}
		return n.r.insert(k)
	}

	return nil
}

func (n *lrbNode) getchild(d Direction) *lrbNode {
	if d == Right {
		return n.r
	}

	return n.l
}

func (n *lrbNode) rotate(dir Direction) {
	child := n.getchild(!dir)

	if dir == Right {
		n.l = child.r
		n.l.p = n
		child.r = n
	} else {
		n.r = child.l
		n.r.p = n
		child.l = n
	}

	child.p = n.p
	child.c = n.c
	n.p = child
	n.c = Red

	if child.p != nil {
		if child.p.l == n {
			child.p.l = child
		} else {
			child.p.r = child
		}
	}
}

func (n *lrbNode) isRBTree(h, min, max int) bool {
	if n == nullNode {
		return h == 0
	}
	if n.c == Black {
		h--
	}
	if n.p != nil && n == n.p.r && n.c == Red {
		return false
	}
	return n.l.isRBTree(h, min, n.k) && n.r.isRBTree(h, n.k, max)
}

func (n *lrbNode) height() int {
	if n == nullNode {
		return 0
	}
	if n.c == Red {
		return n.l.height()
	}
	return 1 + n.l.height()
}

// find succ in n's right child tree
func (n *lrbNode) succ() *lrbNode {
	if n.r == nullNode {
		panic("right child is empty")
	}
	s := n.r
	for s.l != nullNode {
		s = s.l
	}
	return s
}

type LLRBTree struct {
	root *lrbNode
}

func (t *LLRBTree) height() int {
	return t.root.height()
}

func (t *LLRBTree) fix(n *lrbNode) {
	if n == nil || n.c == Black {
		return
	}

	if n.p == nil {
		t.root = n
		n.c = Black
		return
	}

	if n == n.p.r {
		if n.p.l.c == Black {
			n.p.rotate(Left)
			n = n.l
		} else {
			n.p.l.c = Black
			n.c = Black
			n.p.c = Red
		}
	}

	if n.l.c == Red {
		n.p.rotate(Right)
		n.c = Red
		n.l.c = Black
		n.r.c = Black
		n = n.l
	}

	t.fix(n.p)
}

func (t *LLRBTree) fixDel(n *lrbNode) {
	// from null node as left child of a leaf, right sibling must be black
	// as if right sibling does not exist, then color would be red, we would
	// not enter this func
	for n.p != nil {
		if n == n.p.l {
			s := n.p.r
			c := n.p.c
			n.p.rotate(Left)
			n = s
			if c == Red {
				break
			}
		} else {
			s := n.p.l
			c := n.p.c
			n.p.rotate(Right)
			n = s
			if c == Red {
				n.c = Black
				break
			}
		}
	}
	n.c = Black
	if n.p == nil {
		t.root = n
		if n == nullNode {
			t.root = nil
		}
	}
}

func (t *LLRBTree) Check() bool {
	return t.root.isRBTree(t.height(), math.MinInt, math.MaxInt)
}

func (t *LLRBTree) Visit() string {
	return t.root.preorder()
}

func (t *LLRBTree) IsEmpty() bool {
	return t.root == nil
}

func (t *LLRBTree) Search(k int) bool {
	node := t.root.search(k)
	return node != nil && node.k == k
}

func (t *LLRBTree) Insert(k int) {
	if t.root == nil {
		t.root = &lrbNode{k: k, c: Black, l: nullNode, r: nullNode}
		return
	}

	n := t.root.insert(k)
	if n != nil {
		t.fix(n)
	}
}

func (t *LLRBTree) Delete(k int) {
	n := t.root.search(k)
	if n == nil {
		return
	}

	color := n.c
	succ := n.succ()
	start := succ.l
	// n can only be a leaf node, or has only left red child
	if n.r == nullNode {
		succ = n.l
		start = n.l
	} else { // or both black child
		start = nullNode
		color = succ.c
		succ.c = n.c
		succ.l = n.l
		n.l.p = succ
		if succ != n.r { // succ is not n's direct right child
			succ.r = n.r
			n.r.p = succ
			if succ.p.l == succ {
				succ.p.l = nullNode
			} else {
				succ.p.r = nullNode
			}
			nullNode.p = n.p
		}
	}

	if n.p != nil {
		if n.p.l == n.r {
			n.p.l = succ
		} else {
			n.p.r = succ
		}
	}
	succ.p = n.p // null node parent also be set here
	n.p = nil
	n.l = nil
	n.r = nil

	if color == Red {
		return
	}

	//     n
	//    //
	//   x (R)
	if start != nullNode {
		start.c = Black
		return
	}

	t.fixDel(start)
}
