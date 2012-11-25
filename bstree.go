/*
  Package bstree implements a standard binary search tree.

  You can use values of any type in the tree as long
  as implementations for bstree.Smaller and bstree.Larger
  exist for that type you want to use.

  This implementation is optimized for the case where you
  load a large amount of data into the tree and then
  repeatedly iterate or lookup values.

  The tree is safe for use by concurrent goroutines.

  Example:
    package main

    import (
      "fmt"
      "github.com/lazybeaver/go-bstree"
    )

    func main() {
      tree := bstree.New(bstree.IntSmaller, bstree.IntLarger)
      tree.Insert(10)
      tree.Insert(20)
      tree.Exists(10)
      fmt.Println(tree)
      min := tree.Minimum()
      max := tree.Maximum()
      depth := tree.Depth()
      fmt.Println(min, max, depth)
      tree.Traverse(bstree.InOrder, func(value interface{}) {
        fmt.Println(value)
      })
    }

  TODO:
    Implement one-time rebalancing of the tree.
    Implement deletion.
*/
package bstree

import (
	"container/list"
	"fmt"
	"sync"
)

// The Smaller and Larger interfaces
type Smaller func(value interface{}, other interface{}) bool
type Larger func(value interface{}, other interface{}) bool

// Int versions of Smaller and Larger
func IntSmaller(value interface{}, other interface{}) bool {
	return int(value.(int)) < int(other.(int))
}

func IntLarger(value interface{}, other interface{}) bool {
	return int(value.(int)) > int(other.(int))
}

// _Node represents a single element in the tree
type _Node struct {
	value interface{}
	left  *_Node
	right *_Node
}

func new_Node(value interface{}) *_Node {
	node := new(_Node)
	node.value = value
	return node
}

func (node *_Node) String() string {
	return fmt.Sprintf("{address: %p | value: %v | left: %p | right: %p}", node, node.value, node.left, node.right)
}

// Tree represents a binary search tree
// You can create a initialized Tree using bstree.New(...)
type Tree struct {
	root    *_Node
	smaller Smaller
	larger  Larger
	size    int
	mutex   sync.RWMutex
}

// New creates an initialized tree
// Time-complexity: O(1)
func New(smaller Smaller, larger Larger) *Tree {
	tree := new(Tree)
	tree.smaller = smaller
	tree.larger = larger
	return tree
}

// Size returns the size of the tree
// Time-complexity: O(1)
func (tree *Tree) Size() int {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	return tree.size
}

// String returns the details of the tree as a string
// Time-complexity: O(1)
func (tree *Tree) String() string {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	return fmt.Sprintf("{root: %p | size: %d}", tree.root, tree.size)
}

// Traversal Algorithms
type Traversal int32
type Visitor func(interface{})

const (
	_ Traversal = iota
	PreOrder
	InOrder
	PostOrder
	LevelOrder
)

// Traverse walks the tree using a specified algorithm and calls visitor on each node.
// Time-complexity: O(size)
func (tree *Tree) Traverse(traversal Traversal, visitor Visitor) {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	switch traversal {
	case PreOrder:
		tree.doPreOrder(tree.root, visitor)
	case InOrder:
		tree.doInOrder(tree.root, visitor)
	case PostOrder:
		tree.doPostOrder(tree.root, visitor)
	case LevelOrder:
		tree.doLevelOrder(visitor)
	}
}

func (tree *Tree) doPreOrder(node *_Node, visitor Visitor) {
	if node == nil {
		return
	}
	visitor(node.value)
	tree.doPreOrder(node.left, visitor)
	tree.doPreOrder(node.right, visitor)
}

func (tree *Tree) doInOrder(node *_Node, visitor Visitor) {
	if node == nil {
		return
	}
	tree.doInOrder(node.left, visitor)
	visitor(node.value)
	tree.doInOrder(node.right, visitor)
}

func (tree *Tree) doPostOrder(node *_Node, visitor Visitor) {
	if node == nil {
		return
	}
	tree.doPostOrder(node.left, visitor)
	tree.doPostOrder(node.right, visitor)
	visitor(node.value)
}

func (tree *Tree) doLevelOrder(visitor Visitor) {
	if tree.root == nil {
		return
	}
	queue := list.New()
	queue.PushBack(tree.root)
	for queue.Len() > 0 {
		element := queue.Front()
		node := (*_Node)(element.Value.(*_Node))
		queue.Remove(element)
		visitor(node.value)
		if node.left != nil {
			queue.PushBack(node.left)
		}
		if node.right != nil {
			queue.PushBack(node.right)
		}
	}
}

// Exists check if a value exists in the tree
// Average case time-complexity: O(depth)
// Worst case time-complexity: O(size)
func (tree *Tree) Exists(value interface{}) bool {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	return tree.doExists(tree.root, value)
}

func (tree *Tree) doExists(node *_Node, value interface{}) bool {
	if node == nil {
		return false
	}
	switch {
	case tree.smaller(value, node.value):
		return tree.doExists(node.left, value)
	case tree.larger(value, node.value):
		return tree.doExists(node.right, value)
	}
	return true
}

// Insert adds value to the tree if it doesn't already exist
// Returns true if the value was inserted, false otherwise.
// Average case time-complexity: O(depth)
// Worst case time-complexity: O(size)
func (tree *Tree) Insert(value interface{}) bool {
	tree.mutex.Lock()
	defer tree.mutex.Unlock()
	if tree.root == nil {
		tree.root = new_Node(value)
		tree.size++
		return true
	}
	if tree.doInsert(tree.root, value) {
		tree.size++
		return true
	}
	return false
}

func (tree *Tree) doInsert(node *_Node, value interface{}) bool {
	if node == nil {
		return false
	}
	switch {
	case tree.smaller(value, node.value):
		if node.left == nil {
			node.left = new_Node(value)
			return true
		} else {
			return tree.doInsert(node.left, value)
		}
	case tree.larger(value, node.value):
		if node.right == nil {
			node.right = new_Node(value)
			return true
		} else {
			return tree.doInsert(node.right, value)
		}
	}
	return false
}

// Minimum returns the smallest value in the tree
// Average case time-complexity: O(depth)
// Worst case time-complexity: O(size)
func (tree *Tree) Minimum() interface{} {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	if tree.root == nil {
		return nil
	}
	node := tree.root
	for node.left != nil {
		node = node.left
	}
	return node.value
}

// Maximum returns the largest value in the tree
// Average case time-complexity: O(depth)
// Worst case time-complexity: O(size)
func (tree *Tree) Maximum() interface{} {
	tree.mutex.RLock()
	defer tree.mutex.RUnlock()
	if tree.root == nil {
		return nil
	}
	node := tree.root
	for node.right != nil {
		node = node.right
	}
	return node.value
}

// Depth returns the depth of the tree
// Time-complexity: O(size)
func (tree *Tree) Depth() int {
	return tree.doDepth(tree.root)
}

func (tree *Tree) doDepth(node *_Node) int {
	if node == nil {
		return 0
	}
	left := tree.doDepth(node.left)
	right := tree.doDepth(node.right)
	var depth int
	if left > right {
		depth = left + 1
	} else {
		depth = right + 1
	}
	return depth
}
