package skiplist

import (
	"fmt"
	"math"
)

var (
	DEFAULT_NODE_NUM = 1000
	DEFAULT_STEP     = 4
)

// Node implement the basic element of skiplist
type Node interface {
	Compare(interface{}) int // compare node, -1 less 0 equal, 1 greater
}

// HeadNode defines struct of head of list
// level 0 link the node
// level 1... represents indexes
// para counts sum of levels
// TODO: operation optimization
// TODO: concurrency supported
// TODO: option supported
// TODO: batch CRUD, indexes rebuild, sorted seq insert
// TODO: complex test case, expand test coverage
type HeadNode struct {
	level int // sum of layer of index, excluding index[0], level 0 means no index created.
	head  *IndexNode
	para  []int
}

// IndexNode defines struct of index
type IndexNode struct {
	index []*IndexNode
	Node  Node
}

func newHeadNode() *HeadNode {
	return &HeadNode{
		level: 0,
		head:  &IndexNode{index: make([]*IndexNode, 0, DEFAULT_NODE_NUM>>DEFAULT_STEP), Node: nil},
		// head: new(IndexNode),
		// para: make([]int, DEFAULT_NODE_NUM/DEFAULT_STEP/DEFAULT_STEP),
		para: make([]int, 1, DEFAULT_NODE_NUM>>DEFAULT_STEP),
		// Node: nil,
	}
}

func newIndexNode(node Node, level int) *IndexNode {
	return &IndexNode{index: make([]*IndexNode, level+1), Node: node}
}

func (hn *HeadNode) String() string {
	return fmt.Sprintf("level: %d/t, para: %#v\n head: %#v\n", hn.level, hn.para, hn.head)
}

// Print prints list of node
func (hn *HeadNode) Print() {
	cur := hn.head.index
	for i := 0; i < hn.para[0]; i++ {
		fmt.Printf("%+v ", cur[0].Node)
		cur = cur[0].index
	}
	fmt.Println()
}

// Length return length of list
func (hn *HeadNode) Length() int {
	return hn.para[0]
}

// IsSorted returns true if list is incr
func (hn *HeadNode) IsSorted() bool {
	cur := hn.head.index
	curN := cur[0].Node
	for i := 1; i < hn.para[0]; i++ {
		if curN.Compare(cur[0].Node) > 0 {
			return false
		}
		cur = cur[0].index
		curN = cur[0].Node
	}
	return true
}

// Insert add node into list
func (hn *HeadNode) Insert(node Node) (err error) {
	i := hn.level + 1 // given that if level upgrade
	cur := hn.head
	update := make([]*IndexNode, i+1) // given that if level upgrade and index[0] or len(para)+1
	for ; i >= 0; i-- {
		// the same level index find
		if len(cur.index) > i { //head node index check
			for cur.index[i] != nil {
				if cur.index[i].Node.Compare(node) > 0 {
					break
				}
				// if cur.index[i].Node.Compare(node) == 0 {
				// 	error or ok
				// }
				cur = cur.index[i]
			}
		}
		// record left of insert node
		update[i] = cur
	}

	sum := float64(hn.para[0] + 1)
	step := float64(DEFAULT_STEP)
	j := hn.level
	var level int
	if sum > math.Pow(step, float64(j+1)) {
		level = j + 1
	} else {
		for ; j > 0; j-- {
			if sum > math.Pow(step, float64(j))*float64(hn.para[j]+1) {
				break
			}
		}
		level = j
	}

	newNode := newIndexNode(node, level)
	if level > hn.level {
		hn.para = append(hn.para, 0)
		hn.level++
	}
	for ; level >= 0; level-- {
		hn.para[level]++
		if level >= len(update[level].index) { // head node will be &{index: [], node: nil}
			update[level].index = append(update[level].index, newNode)
			newNode.index[level] = nil
		} else {
			update[level].index[level], newNode.index[level] = newNode, update[level].index[level]
		}
	}

	return
}

// Search test if a node is in list
func (hn *HeadNode) Search(key interface{}) bool {
	return hn.Get(key) != Node(nil)
}

// Get gets the value of key
func (hn *HeadNode) Get(key interface{}) Node {
	cur := hn.head
	i := hn.level
	for ; i >= 0; i-- {
		for cur.index[i] != nil {
			res := cur.index[i].Node.Compare(key)
			if res == 0 {
				return cur.index[i].Node
			}
			if res > 0 {
				break
			}
			cur = cur.index[i]
		}
	}

	return nil
}

// Delete deletes first node find by key
func (hn *HeadNode) Delete(key interface{}) Node {
	// if !hn.Search(key) {
	// 	return nil
	// }
	flag := false
	cur := hn.head
	i := hn.level
	var level int
	var node Node
	path := make([]*IndexNode, i+1)
	for ; i >= 0; i-- {
		for cur.index[i] != nil {
			res := cur.index[i].Node.Compare(key)
			if res == 0 {
				flag = true
				level = i
				node = cur.index[i].Node
				break
			}
			if res > 0 {
				break
			}
			cur = cur.index[i]
		}
		path[i] = cur
	}
	if !flag {
		return nil
	}

	j := level
	for ; j >= 0; j-- {
		if path[j].index[j].index[j] != nil {
			path[j].index[j] = path[j].index[j].index[j]
		} else {
			path[j].index[j] = nil
		}
		hn.para[j]--
	}

	return node
}
