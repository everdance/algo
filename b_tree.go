package algo

import "math"

const (
	defaultpadkey = math.MaxInt
)

// keys slice has same length as childs for internal nodes
// but last entry is always set to default pad key
type btrnode struct {
	leaf   bool
	keys   []int
	parent *btrnode
	childs []*btrnode
}

func (n *btrnode) level() int {
	if n == nil {
		return 0
	}
	l := 1
	for !n.leaf {
		l++
		n = n.childs[0]
	}
	return l
}

func (n *btrnode) valid(level, order, min, max int, root bool) bool {
	klen, clen := len(n.keys), len(n.childs)
	if n.leaf && (klen > order || clen != 0 || level != 1) {
		return false
	}

	if !n.leaf {
		if root && (clen < 2 || clen > order) {
			return false
		}
		// internal non root node has at least U[order/2] children
		if !root && (clen > order || clen < order/2) {
			return false
		}
		if klen != clen-1 || n.keys[klen-1] != defaultpadkey {
			return false
		}
	}

	i := 0
	for i < klen {
		if i < klen-1 && n.keys[i] > n.keys[i+1] {
			return false
		}
		if n.keys[i] > max || n.keys[i] <= min {
			return false
		}
		i++
	}

	for i, v := range n.childs {
		mx := max
		mi := min
		if i > 0 {
			mi = n.keys[i-1]
		}
		if i < klen-1 {
			mx = n.keys[i]
		}
		if !v.valid(level-1, order, mi, mx, false) {
			return false
		}
	}

	return true
}

func (n *btrnode) index(key int) int {
	len := len(n.keys)
	for i := 0; i < len; i++ {
		if key <= n.keys[i] {
			return i
		}
	}
	return len
}

func (n *btrnode) search(key int) *btrnode {
	index := n.index(key)
	if n.leaf {
		if key == n.keys[index] {
			return n
		}
		return nil
	}

	return n.childs[index].search(key)
}

func (n *btrnode) insert(key int, child *btrnode) *btrnode {
	if n == nil {
		n = &btrnode{}
	}

	i := n.index(key)
	n.keys = append(n.keys, 0)
	if !n.leaf {
		n.childs = append(n.childs, nil)
	}

	len := len(n.keys)
	for j := len; j > i; j-- {
		n.keys[j] = n.keys[j-1]
		if !n.leaf {
			n.childs[j] = n.childs[j-1]
		}
	}

	n.keys[i] = key
	if !n.leaf {
		n.childs[i] = child
		child.parent = n
	}
	return n
}

func (n *btrnode) fixInsert(order int) *btrnode {
	index := len(n.keys) / 2
	if len(n.keys) > order {
		splitted := &btrnode{
			parent: n.parent,
			keys:   n.keys[index+1:],
		}
		if !n.leaf {
			splitted.childs = n.childs[index:]
			for _, c := range splitted.childs {
				c.parent = splitted
			}
		}
		k := n.keys[index]
		n.parent = n.parent.insert(k, splitted)
		return n.parent.fixInsert(order)
	}
	return n
}

func (n *btrnode) remove(key int) *btrnode {
	i := n.index(key)
	rkeys, rchilds := []int{}, []*btrnode{}
	if i < len(n.keys)-1 {
		rkeys = n.keys[i+1:]
		rchilds = n.childs[i+1:]
	}
	n.keys = append(n.keys[:i], rkeys...)
	if !n.leaf {
		n.childs = append(n.childs[:i], rchilds...)
	}

	return n
}

type Btree struct {
	Order int
	root  *btrnode
}

func BTree(order int) *Btree {
	return &Btree{Order: order}
}

func (t *Btree) check() bool {
	if t.root == nil {
		return true
	}
	level := t.root.level()
	return t.root.valid(level, t.Order, math.MinInt, math.MaxInt, true)
}

func (t *Btree) Find(key int) bool {
	n := t.root.search(key)
	if n == nil {
		return false
	}
	for k := range n.keys {
		if k == key {
			return true
		}
	}

	return false
}

func (t *Btree) Insert(key int) {
	if t.root == nil {
		t.root = &btrnode{leaf: true, keys: []int{key}}
		return
	}

	n := t.root
	for !n.leaf {
		n = n.childs[n.index(key)]
	}

	i := n.index(key)
	if n.keys[i] == key {
		return
	}

	n.insert(key, nil)
	top := n.fixInsert(t.Order)
	if top.parent == nil {
		t.root = top
	}
}
