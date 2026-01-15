package utility

/**
 * Your Trie object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Insert(word);
 * param_2 := obj.Search(word);
 * param_3 := obj.StartsWith(prefix);
 */



type TrieNode struct {
	child []*TrieNode
	IsEnd bool
}

func NewTrieNode() *TrieNode {
	t := &TrieNode{}
	t.child = make([]*TrieNode, 26)
	for i := 0; i < 26; i++ {
		t.child[i] = nil
	}
	t.IsEnd = false
	return t
}

type Trie struct {
	Root *TrieNode
}

func Constructor() Trie {
	t := NewTrieNode()
	return Trie{t}
}

func (t *Trie) Insert(word string) {
	tmp := t.Root
	for _, x := range word {
		ind := int(x - 'a')
		if tmp.child[ind] == nil {
			tmp.child[ind] = NewTrieNode()
		}
		tmp = tmp.child[ind]

	}
	tmp.IsEnd = true
}

func (this *Trie) Search(word string) bool {
	tmp := this.Root
	for _, x := range word {
		ind := int(x - 'a')
		if tmp.child[ind] == nil {
			return false
		}
		tmp = tmp.child[ind]
	}
	return tmp.IsEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	tmp := t.Root
	for _, x := range prefix {
		ind := int(x - 'a')
		if tmp.child[ind] == nil {
			return false
		}
		tmp = tmp.child[ind]
	}
	return true
}

func (t *Trie) walkTrie(prefix string) *TrieNode {
	node := t.Root

	for _, c := range prefix {
		c := c
		ok := node.child[int(c-'a')]
		if ok == nil {
			return nil
		}
		node = ok
	}
	return node
}

func (t *Trie) FindCompletion(prefix string) []string {
	node := t.walkTrie(prefix)
	if node == nil {
		return nil
	}

	res := make([]string, 0, len(node.child))
	var traverse func(node *TrieNode, curStr string)
	traverse = func(node *TrieNode, curStr string) {
		if node == nil {
			return
		}
		if node.IsEnd {
			res = append(res, curStr)
		}
		for c, childNode := range node.child {
			if childNode != nil {
				traverse(childNode, curStr+string(rune(c+'a')))
			}
		}
	}
	traverse(node, prefix)
	return res
}

