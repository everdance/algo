package algo

import (
	"fmt"
	"strings"
)

type Node struct {
	Key    int
	Parent *Node
	Left   *Node
	Right  *Node
}

func (n *Node) IsLeaf() bool {
	return n != nil && n.Left == nil && n.Right == nil
}

func (n *Node) InOrderVisit() string {
	if n != nil {
		left := n.Left.InOrderVisit()
		right := n.Right.InOrderVisit()

		return strings.Trim(fmt.Sprintf("%s %d %s", left, n.Key, right), " ")
	}

	return ""
}

func (n *Node) Search(k int) *Node {
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

func (n *Node) Successor() *Node {
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

func (n *Node) Insert(k int) {
	if n == nil {
		panic("insert on nil node")
	}

	if k < n.Key {
		if n.Left == nil {
			n.Left = &Node{Key: k, Parent: n}
		} else {
			n.Left.Insert(k)
		}
	}

	if k > n.Key {
		if n.Right == nil {
			n.Right = &Node{Key: k, Parent: n}
		} else {
			n.Right.Insert(k)
		}
	}
}

func (n *Node) Transplant(m *Node) {
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
	}
}

type BST struct {
	root *Node
}

func (t *BST) IsEmpty() bool { return t.root == nil }

func (t *BST) Visit() string { return t.root.InOrderVisit() }

func (t *BST) Search(k int) *Node { return t.root.Search(k) }

func (t *BST) Insert(k int) {
	if t.root == nil {
		t.root = &Node{Key: k}
		return
	}

	t.root.Insert(k)
}

func (t BST) Delete(k int) {
	n := t.root.Search(k)
	if n == nil {
		return
	}

	if n.Right == nil {
		n.Transplant(n.Left)
		if n == t.root {
			t.root = n.Left
		}
	} else if n.Left == nil {
		n.Transplant(n.Right)
		if n == t.root {
			t.root = n.Right
		}
	} else {
		succ := n.Successor()
		succ.Transplant(succ.Right)
		n.Transplant(succ)
		succ.Left = n.Left
		succ.Right = n.Right
		succ.Left.Parent = succ
		succ.Right.Parent = succ
	}
}
