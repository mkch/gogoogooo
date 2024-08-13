// Package variance_test demonstrates covariance and contravariance of slice in go.
// Motivation: https://en.wikipedia.org/wiki/Covariance_and_contravariance_(computer_science)#Arrays
package variance_test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type Animal interface {
	Animal()
}

type Cat interface {
	Cat()
	Animal
}

type Dog interface {
	Dog()
	Animal
}

type CatImpl struct{}

func (i CatImpl) Animal() {
	fmt.Println("CatImpl.Animal()")
}

func (i CatImpl) Cat() {
	fmt.Println("CatImpl.Cat()")
}

type DogImpl struct{}

func (i DogImpl) Animal() {
	fmt.Println("DogImpl.Animal()")
}

func (i DogImpl) Dog() {
	fmt.Println("DogImpl.Dog()")
}

func TestCovariant(t *testing.T) {
	var cats = make([]Cat, 0, 2)
	cats = append(cats, CatImpl{})
	var animals []Animal = Variant[Cat, Animal](cats)
	for _, a := range animals {
		a.Animal() // Safe
	}

	var shouldCrash = func() {
		// https://en.wikipedia.org/wiki/Covariance_and_contravariance_(computer_science)
		// the covariant rule is safe for immutable (read-only) arrays.
		_ = append(animals, DogImpl{}) // Not safe: Add a DogImpl to cats.
		cats = cats[:2]
		for _, c := range cats {
			c.Cat() // Crash at the last element.
		}
	}
	_ = shouldCrash
}

func TestContravariant(t *testing.T) {
	var animals = make([]Animal, 0, 2)
	animals = append(animals, DogImpl{})
	var cats = Variant[Animal, Cat](animals)
	cats = append(cats, CatImpl{}) // Safe

	var shouldCrash = func() {
		// https://en.wikipedia.org/wiki/Covariance_and_contravariance_(computer_science)
		// the contravariant rule would be safe for write-only arrays.
		for _, i := range cats { // Not safe: cats may contain DogImpl.
			i.Cat() // Crash at the first element.
		}
	}
	_ = shouldCrash
}

// Variant do covariance if T1 is assignable to T2, do contravariance of T2 is assignable to T1.
// Panics if T1 or T2 is not an interface type or neither is assignable to the other.
func Variant[T1, T2 any](a []T1) []T2 {
	t1 := reflect.TypeOf((*T1)(nil)).Elem()
	if t1.Kind() != reflect.Interface {
		panic(fmt.Errorf("%v is not an interface", t1))
	}
	t2 := reflect.TypeOf((*T2)(nil)).Elem()
	if t2.Kind() != reflect.Interface {
		panic(fmt.Errorf("%v is not an interface", t2))
	}
	if !t1.AssignableTo(t2) && !t2.AssignableTo(t1) {
		panic(fmt.Errorf("incompatible types: %v and %v", t1.String(), t2.String()))
	}

	return *(*[]T2)(unsafe.Pointer(&a))
}
