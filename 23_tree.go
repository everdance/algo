package algo

import (
	"fmt"
	"math"
	"strings"
)

type NodeType bool

const (
	TwoNode   NodeType = false
	ThreeNode NodeType = true
)

type Node23 struct {
	Keys   []int
	Left   *Node23
	Middle *Node23
	Right  *Node23
	Parent *Node23
}

func (n *Node23) ntype() NodeType {
	if len(n.Keys) > 1 {
		return ThreeNode
	}

	return TwoNode
}

func (n *Node23) preorder() string {
	if n != nil {
		self := fmt.Sprintf("%d", n.Keys[0])
		children := fmt.Sprintf("%s %s", n.Left.preorder(), n.Right.preorder())

		if n.ntype() == ThreeNode {
			self = fmt.Sprintf("<%d, %d>", n.Keys[0], n.Keys[1])
			children = fmt.Sprintf("%s %s %s", n.Left.preorder(),
				n.Middle.preorder(), n.Right.preorder())
		}
		children = strings.Trim(children, " ")

		if children == "" {
			return self
		}

		return fmt.Sprintf("%s {%s}", self, children)
	}

	return ""
}

func (n *Node23) search(k int) *Node23 {
	if n == nil {
		return nil
	}

	if k < n.Keys[0] {
		return n.Left.search(k)
	} else if k > n.Keys[0] {
		if n.ntype() == ThreeNode {
			if k < n.Keys[1] {
				return n.Middle.search(k)
			}
			if k == n.Keys[1] {
				return n
			}
		}
		return n.Right.search(k)
	} else {
		return n
	}
}

func (n *Node23) to2node(keep int) {
	if n.ntype() == TwoNode {
		panic("n is already two node")
	}

	if keep == RightD {
		n.Keys = n.Keys[1:]
	} else {
		n.Keys = n.Keys[:1]
	}

	n.Middle = nil
}

func (n *Node23) to3node(k int, child *Node23) {
	if n.ntype() == ThreeNode {
		panic("n is already three node")
	}

	if n.Keys[0] > k {
		n.Keys = []int{k, n.Keys[0]}
		n.Middle = child
	} else {
		n.Keys = append(n.Keys, k)
		n.Middle = n.Right
		n.Right = child
	}
	if child != nil {
		child.Parent = n
	}
}

// add middle key from child, with right splitted new child c
// returns the top node in the tree where change stop at
func (n *Node23) addKey(k int, c *Node23) *Node23 {
	if n.ntype() == TwoNode {
		n.to3node(k, c)
		return n
	}

	var splitted *Node23
	left, mid, right := n.Keys[0], k, n.Keys[1]
	if k < n.Keys[0] { // added from left child
		left = k
		mid = n.Keys[0]
		splitted = make23Node(right, n.Middle, n.Right, n.Parent)
		n.Keys = []int{left}
		n.Right = c
	} else if k > n.Keys[1] { // added from right child
		mid = n.Keys[1]
		right = k
		splitted = make23Node(right, n.Right, c, n.Parent)
		n.Right = n.Middle
	} else { // added from middle child
		splitted = make23Node(right, c, n.Right, n.Parent)
		n.Right = n.Middle
	}

	if n.Right != nil {
		n.Right.Parent = n
	}
	n.Middle = nil
	n.Keys = n.Keys[:1]

	if n.Parent != nil {
		return n.Parent.addKey(mid, splitted)
	}
	// reached tree root
	newRoot := make23Node(mid, n, splitted, nil)
	n.Parent = newRoot

	return newRoot
}

func (n *Node23) height() int {
	h := 0
	for n != nil {
		n = n.Left
		h++
	}
	return h
}

func (n *Node23) is23tree(h, low, high int) bool {
	if n == nil {
		return h == 0
	}

	h--

	for _, k := range n.Keys {
		if k < low || k > high {
			return false
		}
	}

	if n.ntype() == TwoNode {
		return n.Left.is23tree(h, low, n.Keys[0]) && n.Right.is23tree(h, n.Keys[0], high)
	}

	return n.Left.is23tree(h, low, n.Keys[0]) &&
		n.Middle.is23tree(h, n.Keys[0], n.Keys[1]) &&
		n.Right.is23tree(h, n.Keys[1], high)
}

func (n *Node23) isleaf() bool {
	return n.Left == nil
}

func (n *Node23) smallest() *Node23 {
	for !n.isleaf() {
		n = n.Left
	}
	return n
}

func (n *Node23) biggest() *Node23 {
	for !n.isleaf() {
		n = n.Right
	}
	return n
}

const (
	LeftD   int = 0
	MiddleD     = 1
	RightD      = 2
)

func (n *Node23) branchDir() int {
	if n == nil || n.Parent == nil {
		return -1
	}

	if n == n.Parent.Left {
		return LeftD
	} else if n == n.Parent.Right {
		return RightD
	} else {
		return MiddleD
	}
}

func (n *Node23) getChild(dir int) *Node23 {
	switch dir {
	case LeftD:
		return n.Left
	case MiddleD:
		return n.Middle
	default:
		return n.Right
	}
}

// collapse 2 node n with the two node child
// so both children has same height
func (n *Node23) collapse(dir int) {
	if n.ntype() == ThreeNode {
		panic("must be two node to collapse")
	}

	child := n.getChild(dir)
	if child.ntype() == ThreeNode {
		panic("child must be two node to collapse")
	}

	if dir == LeftD {
		n.Keys = []int{child.Keys[0], n.Keys[0]}
		n.Left = child.Left
		n.Middle = child.Right
	} else {
		n.Keys = append(n.Keys, child.Keys[0])
		n.Middle = child.Left
		n.Right = child.Right
	}

	if n.Middle != nil {
		n.Middle.Parent = n
		n.Left.Parent = n
		n.Right.Parent = n
	} else {
		n.Left = nil
		n.Right = nil
	}

	child.Parent = nil
	child.Left = nil
	child.Right = nil
}

func make23Node(k int, left, right, parent *Node23) *Node23 {
	n := &Node23{Keys: []int{k}, Left: left, Right: right, Parent: parent}
	if left != nil {
		left.Parent = n
	}
	if right != nil {
		right.Parent = n
	}
	return n
}

func (n *Node23) rotate(dir int) {
	var xLeft, xRight *Node23
	var rotateDir int

	// for left or mid child misses one level
	// we rotate
	switch dir {
	case LeftD:
		rotateDir = LeftD
		xLeft = n.Left
		xRight = n.Right
		if n.ntype() == ThreeNode {
			xRight = n.Middle
		}
	case MiddleD:
		rotateDir = LeftD
		xLeft = n.Middle
		xRight = n.Right
	default:
		rotateDir = RightD
		xRight = n.Right
		xLeft = n.Left
		if n.ntype() == ThreeNode {
			xLeft = n.Middle
		}
	}

	k := n.Keys[0]
	if rotateDir == LeftD {
		if n.ntype() == ThreeNode && dir == MiddleD {
			k = n.Keys[1]
		}
		//        n                   p
		//      /   \					      /   \
		//     hj    p r   ->      n     r
		//          / | \				  / \   / \
		//         o  q  s			 hj  o q   s

		//       n  u                q  u
		//      /  \  \     			  / \   \
		//     hj  qs    w    ->	 n   s    w
		//        / |\	 /\				/ \  /\   /\
		//       p  r t v  x			hj p r t  v x
		if xRight.ntype() == ThreeNode {
			node := make23Node(k, xLeft, xRight.Left, n)
			if dir == LeftD {
				n.Left = node
			} else {
				n.Middle = node
			}
			if k == n.Keys[0] {
				n.Keys[0] = xRight.Keys[0]
			} else {
				n.Keys[1] = xRight.Keys[0]
			}
			xRight.Left = xRight.Middle
			xRight.to2node(RightD)
		} else {
			//       n s                  s
			//      / \  \					     /  \
			//     hj  q   v   ->      n q    v
			//        / \	 /\			    / | \  / \
			//       p   r u w			 hj p  r u  w

			xRight.Middle = xRight.Left
			xRight.Left = xLeft
			if xLeft != nil {
				xLeft.Parent = xRight
			}
			xRight.Keys = []int{k, xRight.Keys[0]}
			if dir == LeftD {
				n.Left = n.Middle
				n.to2node(RightD)
			} else {
				n.to2node(LeftD)
			}
		}
	} else {
		if n.ntype() == ThreeNode {
			k = n.Keys[1]
		}

		if xLeft.ntype() == ThreeNode {
			node := make23Node(k, xLeft.Right, xRight, n)
			n.Right = node
			if n.ntype() == ThreeNode {
				n.Keys[1] = xLeft.Keys[1]
			} else {
				n.Keys[0] = xLeft.Keys[1]
			}
			xLeft.Right = xLeft.Middle
			xLeft.to2node(LeftD)
		} else {
			xLeft.Middle = xLeft.Right
			xLeft.Right = xRight
			if xRight != nil {
				xRight.Parent = xLeft
			}
			xLeft.Keys = append(xLeft.Keys, k)
			n.Right = n.Middle
			n.to2node(LeftD)
		}
	}
}

// child in dir is one level less than n's other children
func (n *Node23) balance(dir int) {
	// reached root
	if n == nil {
		return
	}

	var nbDir int
	if dir == LeftD {
		nbDir = RightD
		if n.ntype() == ThreeNode {
			nbDir = MiddleD
		}
	} else if dir == MiddleD {
		nbDir = RightD
	} else {
		nbDir = LeftD
		if n.ntype() == ThreeNode {
			nbDir = MiddleD
		}
	}

	nbCihld := n.getChild(nbDir)
	if n.ntype() == ThreeNode || (nbCihld != nil && nbCihld.ntype() == ThreeNode) {
		n.rotate(dir)
		return
	}

	n.collapse(nbDir)
	n.Parent.balance(n.branchDir())
}

type Tree23 struct {
	root *Node23
}

func (t *Tree23) Insert(k int) {
	if t.root == nil {
		t.root = &Node23{Keys: []int{k}}
		return
	}

	// find the leaf node to insert
	n := t.root
	for n.Left != nil {
		if k == n.Keys[0] || (len(n.Keys) > 1 && k == n.Keys[1]) {
			return
		}
		if k < n.Keys[0] {
			n = n.Left
			continue
		}
		if len(n.Keys) > 1 && k < n.Keys[1] {
			n = n.Middle
			continue
		}
		n = n.Right
	}

	top := n.addKey(k, nil)
	if top.Parent == nil {
		t.root = top
	}
}

func (t *Tree23) Search(k int) *Node23 {
	return t.root.search(k)
}

func (t *Tree23) Check() bool {
	return t.root.is23tree(t.root.height(), math.MinInt32, math.MaxInt64)
}

func (t *Tree23) IsEmpty() bool {
	return t.root == nil
}

func (t *Tree23) Delete(k int) {
	n := t.root
	for n != nil {
		if k == n.Keys[0] || (len(n.Keys) > 1 && k == n.Keys[1]) {
			break
		}
		if k < n.Keys[0] {
			n = n.Left
			continue
		}
		if len(n.Keys) > 1 && k < n.Keys[1] {
			n = n.Middle
			continue
		}
		n = n.Right
	}

	if n == nil {
		return
	}

	// choose a key from leaf to do balance from bottom
	x, xkey := n, k
	if !n.isleaf() {
		if k == n.Keys[0] {
			x = n.Left.biggest()
			xkey = x.Keys[0]
			if x.ntype() == ThreeNode {
				xkey = x.Keys[1]
			}
			n.Keys[0] = xkey
		} else {
			x = n.Right.smallest()
			xkey = x.Keys[0]
			n.Keys[1] = xkey
		}
	}

	if x.ntype() == ThreeNode {
		if xkey == x.Keys[0] {
			x.Keys = x.Keys[1:]
		} else {
			x.Keys = x.Keys[:1]
		}
		return
	}

	parent := x.Parent
	if parent == nil {
		t.root = nil
		return
	}

	dir := x.branchDir()
	x.Parent = nil
	switch dir {
	case LeftD:
		parent.Left = nil
	case RightD:
		parent.Right = nil
	case MiddleD:
		parent.Middle = nil
	}

	parent.balance(dir)
}

func (t *Tree23) Visit() string {
	return t.root.preorder()
}
