package algo

import (
	"fmt"
	"math"
	"strings"
)

type AvlNode struct {
	k int
	h int
	p *AvlNode
	l *AvlNode
	r *AvlNode
}

func (n *AvlNode) preorder() string {
	if n != nil {
		children := fmt.Sprintf("%s %s", n.l.preorder(),
			n.r.preorder())
		children = strings.Trim(children, " ")

		if children == "" {
			return fmt.Sprintf("%d<%d>", n.k, n.h)
		}

		return fmt.Sprintf("%d<%d> {%s}", n.k, n.h, children)
	}

	return ""
}

func (n *AvlNode) isAvl(min, max int) bool {
	if n == nil {
		return true
	}

	if !n.heightmatch() {
		return false
	}

	if n.l != nil && n.l.p != n {
		return false
	}

	if n.r != nil && n.r.p != n {
		return false
	}

	return n.k > min && n.k < max && n.balanced() &&
		n.l.isAvl(min, n.k) && n.r.isAvl(n.k, max)
}

func (n *AvlNode) search(k int) *AvlNode {
	if n == nil {
		return nil
	}
	if k < n.k {
		return n.l.search(k)
	}
	if k > n.k {
		return n.r.search(k)
	}

	return n
}

func (n *AvlNode) height() int {
	if n == nil {
		return 0
	}
	return n.h
}

func (n *AvlNode) insert(k int) *AvlNode {
	if n == nil {
		panic("insert on nil node")
	}

	if k < n.k {
		if n.l == nil {
			n.l = &AvlNode{k: k, p: n, h: 1}
			return n.l
		}
		return n.l.insert(k)
	}

	if k > n.k {
		if n.r == nil {
			n.r = &AvlNode{k: k, p: n, h: 1}
			return n.r
		}
		return n.r.insert(k)
	}

	return nil
}

func (n *AvlNode) balanced() bool {
	delta := n.l.height() - n.r.height()
	return delta >= -1 && delta <= 1
}

func (n *AvlNode) heightmatch() bool {
	if n == nil {
		return true
	}

	h := n.l.height()
	if h < n.r.height() {
		h = n.r.height()
	}
	return n.h == h+1
}

func (n *AvlNode) updateh() {
	leftH, rightH := n.l.height(), n.r.height()
	if leftH > rightH {
		n.h = leftH + 1
	} else {
		n.h = rightH + 1
	}
}

func (n *AvlNode) rightRotate() {
	p := n.p
	n.l.p = n.p
	if p != nil {
		if p.l == n {
			p.l = n.l
		} else {
			p.r = n.l
		}
	}

	n.p = n.l
	n.l = n.l.r
	if n.l != nil {
		n.l.p = n
	}
	n.p.r = n

	n.updateh()
	n.p.updateh()
}

func (n *AvlNode) leftRotate() {
	p := n.p
	n.r.p = n.p
	if p != nil {
		if p.l == n {
			p.l = n.r
		} else {
			p.r = n.r
		}
	}

	n.p = n.r
	n.r = n.r.l
	if n.r != nil {
		n.r.p = n
	}
	n.p.l = n

	n.updateh()
	n.p.updateh()
}

func (n *AvlNode) succ() *AvlNode {
	if n.r != nil {
		s := n.r
		for s.l != nil {
			s = s.l
		}
		return s
	}

	return n.p
}

func (n *AvlNode) rotate() *AvlNode {
	hl, hr := n.l.height(), n.r.height()
	if hl > hr {
		if n.l.l.height() < n.l.r.height() {
			n.l.leftRotate()
		}
		n.rightRotate()
	} else {
		if n.r.l.height() > n.r.r.height() {
			n.r.rightRotate()
		}
		n.leftRotate()
	}

	return n.p
}

func (n *AvlNode) balance() *AvlNode {
	n.updateh()

	if !n.balanced() {
		n = n.rotate()
	}

	if n.p != nil {
		return n.p.balance()
	}

	return n
}

type AvlTree struct {
	root *AvlNode
}

func (t *AvlTree) IsEmpty() bool { return t.root == nil }

func (t *AvlTree) Check() bool { return t.root.isAvl(math.MinInt, math.MaxInt) }

func (t *AvlTree) Visit() string { return t.root.preorder() }

func (t *AvlTree) Search(k int) *AvlNode { return t.root.search(k) }

func (t *AvlTree) Insert(k int) {
	if t.root == nil {
		t.root = &AvlNode{k: k, h: 1}
		return
	}

	node := t.root.insert(k)
	if node == nil {
		return
	}

	if top := node.balance(); top != nil {
		t.root = top
	}
}

func (t *AvlTree) Delete(k int) {
	n := t.root
	for n != nil {
		if n.k == k {
			break
		}

		if n.k > k {
			n = n.l
		} else {
			n = n.r
		}
	}

	if n == nil {
		return
	}

	var succ *AvlNode
	bottom := n
	if n.l == nil && n.r == nil {
		bottom = n.p
	} else if n.l == nil {
		bottom = n.r
		succ = n.r
	} else if n.r == nil {
		bottom = n.l
		succ = n.l
	} else {
		succ = n.succ()
		bottom = succ.r
		if bottom == nil {
			if succ == n.r {
				bottom = succ
			} else {
				bottom = succ.p
			}
		}
		if succ != n.r {
			if succ.r != nil {
				succ.r.p = succ.p
			}
			succ.p.l = succ.r
		}
	}

	if succ != nil {
		succ.p = n.p
		if succ != n.l {
			succ.l = n.l
		}
		if succ != n.r {
			succ.r = n.r
		}
		if succ.l != nil {
			succ.l.p = succ
		}
		if succ.r != nil {
			succ.r.p = succ
		}
		succ.updateh()
	}

	if n.p != nil {
		if n == n.p.l {
			n.p.l = succ
		} else {
			n.p.r = succ
		}
	}

	n.p = nil
	n.l = nil
	n.r = nil
	if bottom == nil {
		t.root = nil
		return
	}

	if top := bottom.balance(); top != nil {
		t.root = top
	}
}
