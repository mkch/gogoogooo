package main

import (
	"fmt"
	"net/http"
)

/*
Motivation: Generics pooling where each type has its own pool.
*/

// Pool[T] is a pool of []T.
type Pool[T any] struct {
	_ int // Don't remove. Zero-size is not welcomed here.
}

// Alloc allocates an []T form pool.
func (pool *Pool[T]) Alloc() []T {
	return make([]T, 0) // Just conceptual.
}

// tag[T] is the key of poolMap.
// Zero-size is OK here, because its value will be used.
type tag[T any] struct{}

var poolMap = make(map[any]any)

// SharedAlloc[T] allocates []T from the pool of T.
func SharedAlloc[T any]() []T {
	var pool *Pool[T]
	if anyPool, ok := poolMap[tag[T]{}]; !ok {
		anyPool = new(Pool[T])
		poolMap[tag[T]{}] = anyPool
		pool = anyPool.(*Pool[T])
	} else {
		pool = anyPool.(*Pool[T])
	}

	fmt.Printf("Using pool %T %[1]p\n", pool)

	return pool.Alloc()
}

func main() {
	SharedAlloc[int]()
	SharedAlloc[string]()
	SharedAlloc[int]()
	SharedAlloc[string]()
	SharedAlloc[*http.Request]()
	SharedAlloc[*http.Request]()
	SharedAlloc[func()]()
	SharedAlloc[func()]()

}
