package bstree

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

// EmptyTree generates an empty tree of integers
func EmptyTree() *Tree {
	return New(IntSmaller, IntLarger)
}

// RandomTree generates a binary search tree using random integers
// You can control the shape of the tree using count and max
// count is the number of nodes in the tree
// max is the maximum value of the random number generator
func RandomTree(count int, max int) *Tree {
	if count > max {
		panic("Cannot generate more random values than max")
	}
	tree := EmptyTree()
	for tree.Size() < count {
		tree.Insert(rand.Intn(max))
	}
	return tree
}

// CompleteTree generates a complete binary search tree
// using sequential integers
func CompleteTree(count int) *Tree {
	tree := EmptyTree()
	doCompleteTree(tree, 1, count)
	return tree
}

func doCompleteTree(tree *Tree, begin int, end int) {
	switch {
	case begin > end:
		return
	case begin <= end:
		mid := (begin + end) / 2
		tree.Insert(mid)
		doCompleteTree(tree, begin, mid-1)
		doCompleteTree(tree, mid+1, end)
	}
}

func TestTree_MinMax(t *testing.T) {
	tree := RandomTree(1000, 1000)
	if actual := tree.Minimum(); 0 != actual {
		t.Errorf("Minimum: {Expected=0 | Actual=%d}", actual)
	}
	if actual := tree.Maximum(); 999 != actual {
		t.Errorf("Maximum: {Expected=999 | Actual=%d}", actual)
	}
}

func TestTree_Depth(t *testing.T) {
	for i := 0; i <= 1024; i++ {
		tree := CompleteTree(i)
		if i != tree.Size() {
			t.Errorf("CompleteTree Size: {Expected: %d | Actual: %d}", i, tree.Size())
		}
		expected := int(math.Ceil(math.Log2(float64(tree.Size() + 1))))
		if expected != tree.Depth() {
			t.Errorf("CompleteTree Depth: {Expected: %d | Actual: %d}", expected, tree.Depth())
		}
	}
}

// Make concurrent goroutines insert different ranges into the tree
func TestTree_InsertParallel(t *testing.T) {
	numroutines := runtime.NumCPU() * 2
	prev := runtime.GOMAXPROCS(numroutines)
	defer runtime.GOMAXPROCS(prev)

	const numinsert = 1000 // per goroutine
	var wg sync.WaitGroup

	tree := EmptyTree()
	for i := 0; i < numroutines; i++ {
		begin := i * numinsert
		end := begin + numinsert - 1
		wg.Add(1)
		go func() {
			// This is NOT a complete tree. We are just using that function handily.
			doCompleteTree(tree, begin, end)
			wg.Done()
		}()
	}

	wg.Wait()
	expected := numroutines * numinsert
	if expected != tree.Size() {
		t.Errorf("Tree Size: {Expected: %d | Actual: %d}", expected, tree.Size())
	}
}

// Do some simple traversal and do blackbox tests
func ExampleTree_Traverse() {
	tree := CompleteTree(15)
	tree.Traverse(PreOrder, func(value interface{}) {
		fmt.Printf("%d,", value)
	})
	fmt.Printf("\n")
	tree.Traverse(InOrder, func(value interface{}) {
		fmt.Printf("%d,", value)
	})
	fmt.Printf("\n")
	tree.Traverse(PostOrder, func(value interface{}) {
		fmt.Printf("%d,", value)
	})
	fmt.Printf("\n")
	tree.Traverse(LevelOrder, func(value interface{}) {
		fmt.Printf("%d,", value)
	})
	fmt.Printf("\n")
	// Output:
	// 8,4,2,1,3,6,5,7,12,10,9,11,14,13,15,
	// 1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,
	// 1,3,2,5,7,6,4,9,11,10,13,15,14,12,8,
	// 8,4,12,2,6,10,14,1,3,5,7,9,11,13,15,
}

// Do some simple inserts and do blackbox tests
func ExampleTree_Insert() {
	tree := New(IntSmaller, IntLarger)
	tree.Insert(20)
	tree.Insert(10)
	tree.Insert(30)
	tree.Insert(30)
	fmt.Println(tree.Exists(20))
	fmt.Println(tree.Exists(50))
	fmt.Println(tree.Exists(30))
	// Output:
	// true
	// false
	// true
}

// Benchmark insert performance on a single core
func BenchmarkTreeInsert(b *testing.B) {
	b.StopTimer()
	tree := EmptyTree()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(rand.Int())
	}
}
