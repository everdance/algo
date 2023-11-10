package algo

type NodeType bool

const (
	TwoNode   NodeType = false
	ThreeNode NodeType = true
)

type Node23 struct {
	Key    int
	Key2   *int
	Left   *Node23
	Middle *Node23
	Right  *Node23
	Parent *Node23
}

func (n *Node23) ntype() NodeType {
	if n.Key2 != nil {
		return ThreeNode
	}

	return TwoNode
}

func (n *Node23) search(k int) *Node23 {
	if n == nil {
		return nil
	}

	if k < n.Key {
		return n.Left.search(k)
	} else if k > n.Key {
		if n.Key2 != nil {
			if k < *n.Key2 {
				return n.Middle.search(k)
			}
			if k == *n.Key2 {
				return n
			}
		}
		return n.Right.search(k)
	} else {
		return n
	}
}

func (n *Node23) to2node(dir int) {
	if n.ntype() == TwoNode {
		panic("n is already two node")
	}

	if dir == RightD {
		n.Key = *n.Key2
	}

	n.Middle = nil
	n.Key2 = nil
}

func (n *Node23) to3node(k int, child *Node23) {
	if n.ntype() == ThreeNode {
		panic("n is already three node")
	}

	if n.Key > k {
		k1 := n.Key
		n.Key2 = &k1
		n.Key = k
		n.Middle = child
	} else {
		n.Key2 = &k
		n.Middle = n.Right
		n.Right = child
	}
	if child != nil {
		child.Parent = n
	}
}

// addKey middle key from child, with right splitted new child c
// returns the top node in the tree where change stop at
func (n *Node23) addKey(k int, c *Node23) *Node23 {
	if n.ntype() == TwoNode {
		n.to3node(k, c)
		return n
	}

	var splitted *Node23
	left, mid, right := n.Key, k, *n.Key2
	if k < n.Key { // added from left child
		left = k
		mid = n.Key
		splitted = &Node23{Key: right, Left: n.Middle, Right: n.Right}
		n.Key = left
		n.Right = c
	} else if k > *n.Key2 { // added from right child
		mid = *n.Key2
		right = k
		splitted = &Node23{Key: right, Left: n.Right, Right: c}
		n.Right = n.Middle
	} else { // added from middle child
		splitted = &Node23{Key: right, Left: c, Right: n.Right}
		n.Right = n.Middle
	}

	if splitted.Left != nil {
		splitted.Left.Parent = splitted
	}
	if splitted.Right != nil {
		splitted.Right.Parent = splitted
	}
	if n.Right != nil {
		n.Right.Parent = n
	}
	n.Middle = nil
	n.Key2 = nil

	if n.Parent != nil {
		return n.Parent.addKey(mid, splitted)
	}
	// reached tree root
	newRoot := &Node23{Key: mid, Left: n, Right: splitted}
	n.Parent = newRoot
	splitted.Parent = newRoot
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

func (n *Node23) is23tree(h int) bool {
	if n == nil {
		return h == 0
	}
	h--
	ok := n.Left.is23tree(h) && n.Right.is23tree(h)
	if n.ntype() == TwoNode {
		return ok
	}
	return ok && n.Middle.is23tree(h)
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
		n.Key2 = &n.Key
		n.Key = child.Key
		n.Left = child.Left
		n.Middle = child.Right
	} else {
		n.Key2 = &child.Key
		n.Middle = child.Left
		n.Right = child.Right
	}

	if n.Middle != nil {
		n.Middle.Parent = n
		n.Left.Parent = n
	} else {
		n.Left = nil
		n.Right = nil
	}

	child.Parent = nil
	child.Left = nil
	child.Right = nil
}

func (n *Node23) rotate(dir int) {
	var xLeft, xRight *Node23
	var rotateDir int

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

	k := n.Key
	if n.ntype() == ThreeNode {
		k = *n.Key2
	}

	if rotateDir == LeftD {
		if xRight.ntype() == ThreeNode {
			node := &Node23{Key: k, Left: xLeft, Right: xRight.Left, Parent: n}
			if dir == LeftD {
				n.Left = node
			} else {
				n.Middle = node
			}
			node.Left.Parent = node
			node.Right.Parent = node
			n.Key = xRight.Key
			xRight.Left = xRight.Middle
			xRight.to2node(RightD)
		} else {
			xRight.Middle = xRight.Left
			xRight.Left = xLeft
			xLeft.Parent = xRight
			xRight.Key2 = &xRight.Key
			xRight.Key = k
			n.Left = n.Middle
			n.to2node(RightD)
		}
	} else {
		if xLeft.ntype() == ThreeNode {
			node := &Node23{Key: k, Left: xLeft.Right, Right: xRight, Parent: n}
			n.Right = node
			node.Left.Parent = node
			node.Right.Parent = node
			n.Key2 = xLeft.Key2
			xLeft.Right = xLeft.Middle
			xLeft.to2node(LeftD)
		} else {
			xLeft.Middle = xLeft.Right
			xLeft.Right = xRight
			xRight.Parent = xLeft
			xLeft.Key2 = n.Key2
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

	if n.ntype() == ThreeNode || n.getChild(nbDir).ntype() == ThreeNode {
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
		t.root = &Node23{Key: k}
		return
	}

	// find the leaf node to insert
	n := t.root
	for n.Left != nil {
		if k == n.Key || (n.Key2 != nil && k == *n.Key2) {
			return
		}
		if k < n.Key {
			n = n.Left
			continue
		}
		if n.Key2 != nil && k < *n.Key2 {
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
	return t.root.is23tree(t.root.height())
}

func (t *Tree23) Delete(k int) {
	n := t.root
	for n != nil {
		if k == n.Key || (n.Key2 != nil && k == *n.Key2) {
			break
		}
		if k < n.Key {
			n = n.Left
			continue
		}
		if n.Key2 != nil && k < *n.Key2 {
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
		if k == n.Key {
			x = n.Left.biggest()
			xkey = x.Key
			if x.ntype() == ThreeNode {
				xkey = *x.Key2
			}
			n.Key = xkey
		} else {
			x = n.Right.smallest()
			xkey = x.Key
			n.Key2 = &xkey
		}
	}

	if x.ntype() == ThreeNode {
		if xkey == x.Key {
			x.Key = *x.Key2
		}
		x.Key2 = nil
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
