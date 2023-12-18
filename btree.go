package algo

import (
	"fmt"
	"math"
)

// keys slice has same length as childs for internal nodes
// but last entry is always set to default pad key
type btrnode struct {
	leaf   bool
	keys   []int
	parent *btrnode
	childs []*btrnode
}

func (n *btrnode) print() {
	if n == nil {
		return
	}
	nodetype := "Internal"
	if n.leaf {
		nodetype = "Leaf"
	}
	fmt.Printf("[%s] %v\n", nodetype, n)
	for _, c := range n.childs {
		c.print()
	}
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
		fmt.Println("leaf key and child len range, level", n)
		return false
	}

	if !n.leaf {
		if root && (clen < 2 || clen > order) {
			fmt.Println("root child len range", n)
			return false
		}
		// internal non root node has at least U[order/2] children
		if !root && (clen > order || clen < order/2) {
			fmt.Println("child len range", n)
			return false
		}
		if klen != clen-1 {
			fmt.Println("key child len", n)
			return false
		}
	}

	for i := 0; i < klen; i++ {
		if i < klen-1 && n.keys[i] > n.keys[i+1] {
			fmt.Println("key order invalid", n.keys)
			return false
		}
		if n.keys[i] > max || n.keys[i] <= min {
			fmt.Println("key range invalid, ", max, min, n)
			return false
		}
	}

	for i, v := range n.childs {
		mx := max
		mi := min
		if i > 0 {
			mi = n.keys[i-1]
		}
		if i <= klen-1 {
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
		if index < len(n.keys) && n.keys[index] == key {
			return n
		}
		return nil
	}

	return n.childs[index].search(key)
}

func (n *btrnode) insert(key int, child *btrnode) *btrnode {
	i := n.index(key)
	n.keys = append(n.keys, 0)
	if !n.leaf {
		n.childs = append(n.childs, nil)
	}

	for j := len(n.keys) - 1; j > i; j-- {
		n.keys[j] = n.keys[j-1]
	}
	n.keys[i] = key

	if !n.leaf {
		for j := len(n.childs) - 1; j > i+1; j-- {
			n.childs[j] = n.childs[j-1]
		}
		n.childs[i+1] = child
		child.parent = n
	}

	return n
}

func (n *btrnode) fixInsert(order int) *btrnode {
	index := len(n.keys) / 2
	if len(n.keys) > order {
		splitted := &btrnode{
			parent: n.parent,
			leaf:   n.leaf,
			keys:   append([]int{}, n.keys[index+1:]...),
		}
		if !n.leaf {
			splitted.childs = append([]*btrnode{}, n.childs[index:]...)
			for _, c := range splitted.childs {
				c.parent = splitted
			}
			n.childs = n.childs[:index]
		}
		k := n.keys[index]
		n.keys = n.keys[:index+1]
		// create new root
		if n.parent == nil {
			n.parent = &btrnode{childs: []*btrnode{n, splitted}, keys: []int{k}}
			splitted.parent = n.parent
			return n.parent
		}
		_ = n.parent.insert(k, splitted)
		return n.parent.fixInsert(order)
	}
	return n
}

func (n *btrnode) remove(key int) *btrnode {
	i := n.index(key)
	keys := append([]int{}, n.keys[:i]...)
	for j := i + 1; j < len(n.keys); j++ {
		keys = append(keys, n.keys[j])
	}
	n.keys = keys
	if !n.leaf {
		childs := append([]*btrnode{}, n.childs[:i+1]...)
		for j := i + 2; j < len(n.childs); j++ {
			childs = append(childs, n.childs[j])
		}
		n.childs = childs
	}

	return n
}

func (n *btrnode) delete(key, order int) *btrnode {
	n.remove(key)

	l := len(n.keys)
	if l >= order/2 {
		return n
	}

	p := n.parent
	if p == nil {
		if l == 0 {
			return nil
		}
		return n
	}

	ni := p.index(n.keys[0])
	si := ni + 1 // right sibling index
	if si > len(p.keys) {
		si = ni - 1 // left sibling index
	}
	s := p.childs[si]

	if len(s.keys) > order/2 { // borrow from sibling
		if si > ni {
			n.keys = append(n.keys, s.keys[0])
			p.keys[ni] = s.keys[0]
			s.keys = s.keys[1:]
			if !n.leaf {
				n.childs = append(n.childs, s.childs[0])
				s.childs[0].parent = n
				s.childs = s.childs[1:]
			}
		} else {
			sklen := len(s.keys)
			n.keys = append([]int{s.keys[sklen-1]}, n.keys...)
			s.keys = s.keys[:len(s.keys)-1]
			p.keys[si] = s.keys[sklen-2]
			if !n.leaf {
				child := s.childs[len(s.childs)-1]
				n.childs = append([]*btrnode{child}, n.childs...)
				child.parent = n
				s.childs = s.childs[:len(s.childs)-1]
			}
		}
		return n
	}

	to, from, key := n, s, p.keys[ni]
	if si < ni {
		to, from, key = s, n, p.keys[si]
	}
	// merge with sibiling
	to.keys = append(to.keys, from.keys...)
	to.childs = append(to.childs, from.childs...)
	from.parent = nil
	from.keys = nil
	from.childs = nil

	return p.delete(key, order)
}

type Btree struct {
	order int
	root  *btrnode
}

func NewBTree(order int) *Btree {
	return &Btree{order: order}
}

func (t *Btree) IsEmpty() bool {
	return t.root == nil
}

func (t *Btree) Check() bool {
	if t.root == nil {
		return true
	}
	level := t.root.level()
	return t.root.valid(level, t.order, math.MinInt, math.MaxInt, true)
}

func (t *Btree) Print() {
	t.root.print()
}

func (t *Btree) Search(key int) bool {
	return t.root.search(key) != nil
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
	if i < len(n.keys) && n.keys[i] == key {
		return
	}

	n.insert(key, nil)
	top := n.fixInsert(t.order)
	if top.parent == nil {
		t.root = top
	}
}

func (t *Btree) Delete(key int) {
	if t.root == nil {
		return
	}

	n := t.root
	for !n.leaf {
		n = n.childs[n.index(key)]
	}

	i := n.index(key)
	if n.keys[i] == key {
		top := n.delete(key, t.order)
		if top == nil {
			t.root = nil
		} else if top.parent == nil {
			t.root = top
		}
	}
}
