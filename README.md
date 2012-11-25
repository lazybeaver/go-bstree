go-bstree
=========
Go implementation of a binary search tree

This implementation is optimized for the case where you load a large amount of data into the tree and then repeatedly iterate or lookup values. It is safe for use by concurrent goroutines.


Documentation
=============
    godoc github.com/lazybeaver/go-bstree

Benchmarks
==========
    go test --bench '.*'
    PASS
    BenchmarkTreeInsert      1000000              2212 ns/op
    ok      github.com/lazybeaver/go-bstree 2.636s
