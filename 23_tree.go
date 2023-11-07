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

// add the borrow key to child that need one key to fix
func (n *Node23) borrow(k int, p *Node23) *Node23 {
	if n == nil {
		return &Node23{Key: k, Parent: p}
	}
	return nil
}

// try to fix balance within n and its children
func (n *Node23) balance(dir int) (*int, bool) {
	// find a key to borrow from other children

	// or collapse with the other child to make all children height the same
	return nil, false
}

func (n *Node23) fix(dir int) {
	k, ok := n.balance(dir)
	if !ok {
		n.Parent.fix(n.branchDir())
		return
	}

	if k != nil {
		switch dir {
		case LeftD:
			n.Left = n.Left.borrow(*k, n)
		case RightD:
			n.Right = n.Right.borrow(*k, n)
		case MiddleD:
			n.Middle = n.Middle.borrow(*k, n)
		}
	}
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

	parent.fix(dir)
}
