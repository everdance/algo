package algo

import (
	"fmt"
	"strings"
)

// https://algs4.cs.princeton.edu/33balanced

type lrbNode struct {
	k int
	c Color
	l *lrbNode
	r *lrbNode
	p *lrbNode
}

var nullNode = lrbNode{
	c: Black,
}

func (n *lrbNode) preorder() string {
	if n == &nullNode {
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
	if n == nil {
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
		if n.l == &nullNode {
			n.l = &lrbNode{k: k, p: n, c: Red, l: &nullNode, r: &nullNode}
			return n.l
		}
		return n.l.insert(k)
	}

	if k > n.k {
		if n.r == &nullNode {
			n.r = &lrbNode{k: k, p: n, c: Red, l: &nullNode, r: &nullNode}
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

func (n *lrbNode) rotate(dir Direction, tree *LLRBTree) {
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
	} else {
		tree.root = child
		child.c = Black
	}
}

type LLRBTree struct {
	root *lrbNode
}

func (t *LLRBTree) fix(n *lrbNode) {
	if n.c == Black || n == nil {
		return
	}

	if n == n.p.r {
		n.p.rotate(Left, t)
	}

	if n.l.c == Red {
		n.p.rotate(Right, t)
		n.c = Red
		n.l.c = Black
		n.r.c = Black
	}

	t.fix(n)
}

func (t *LLRBTree) Insert(k int) {
	if t.root == nil {
		t.root = &lrbNode{k: k, l: &nullNode, r: &nullNode}
		return
	}

	n := t.root.insert(k)
	if n != nil {
		t.fix(n)
	}
}
