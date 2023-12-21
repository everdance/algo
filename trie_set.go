package algo

import (
	"fmt"
	"unicode/utf8"
)

// only allow ascii characters 0 - 127
const CharNum = 128

var notSupportedRuneErr = fmt.Errorf("only supports ascii chars 0-127")

type trieNode struct {
	next     [CharNum]*trieNode
	isString bool
}

func (n *trieNode) print() {
	for i, c := range n.next {
		if c != nil {
			fmt.Printf("%c[%v] ", i, c.isString)
			c.print()
			fmt.Println("")
		}
	}
}

func (n *trieNode) get(s []byte, d int) *trieNode {
	if n == nil {
		return nil
	}
	if d == len(s) {
		return n
	}
	return n.next[s[d]].get(s, d+1)
}

func (n *trieNode) add(s []byte, d int) *trieNode {
	if n == nil {
		n = &trieNode{}
	}
	if d == len(s) {
		n.isString = true
	} else {
		n.next[s[d]] = n.next[s[d]].add(s, d+1)
	}
	return n
}

func (n *trieNode) remove(s []byte, d int) *trieNode {
	if n == nil {
		return nil
	}

	if d == len(s) {
		n.isString = false
	} else {
		n.next[s[d]] = n.next[s[d]].remove(s, d+1)
	}

	for _, c := range n.next {
		if c != nil {
			return n
		}
	}

	return nil
}

func (n *trieNode) childs() [][]byte {
	bs := [][]byte{}
	if n.isString {
		bs = append(bs, nil)
	}
	for i := range n.next {
		if n.next[i] != nil {
			for _, b := range n.next[i].childs() {
				bs = append(bs, append([]byte{byte(i)}, b...))
			}
		}
	}
	return bs
}

type TrieSet struct {
	root *trieNode
}

func (s *TrieSet) Print() {
	s.root.print()
}

func (s *TrieSet) IsEmpty() bool {
	return s.root == nil
}

func sanitize(key string) ([]byte, error) {
	b := []byte(key)
	for len(b) > 0 {
		_, size := utf8.DecodeLastRune(b)
		if size > 1 {
			return nil, notSupportedRuneErr
		}
		b = b[:len(b)-size]
	}
	return []byte(key), nil
}

func (s *TrieSet) Contains(key string) bool {
	bs, err := sanitize(key)
	if err != nil {
		return false
	}
	n := s.root.get(bs, 0)
	return n != nil && n.isString
}

func (s *TrieSet) Del(key string) error {
	bs, err := sanitize(key)
	if err != nil {
		return err
	}

	s.root.remove(bs, 0)
	return nil
}

func (s *TrieSet) Put(key string) error {
	bs, err := sanitize(key)
	if err != nil {
		return err
	}

	s.root = s.root.add(bs, 0)
	return nil
}

func (s *TrieSet) KeysWithPrefix(prefix string) []string {
	bs, err := sanitize(prefix)
	if err != nil {
		return nil
	}

	n := s.root.get(bs, 0)
	if n == nil {
		return nil
	}
	results := []string{}
	for _, c := range n.childs() {
		results = append(results, prefix+string(c))
	}

	return results
}
