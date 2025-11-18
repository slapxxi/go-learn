package main

import (
	"errors"
	"fmt"
)

type node[T any] struct {
	Data T
	next *node[T]
}

type list[T any] struct {
	start *node[T]
}

func (l *list[T]) add(data T) {
	n := node[T]{
		Data: data,
		next: nil,
	}
	if l.start == nil {
		l.start = &n
		return
	}
	if l.start.next == nil {
		l.start.next = &n
		return
	}
	temp := l.start
	l.start = l.start.next
	l.add(data)
	l.start = temp
}

type TreeLast[T any] []T

func (t TreeLast[T]) replaceLast(element T) (TreeLast[T], error) {
	if len(t) == 0 {
		return t, errors.New("This is empty!")
	}
	t[len(t)-1] = element
	return t, nil
}

type AnotherInt int

type AllInts interface {
	~int
}

type Numeric interface {
	int | int8 | int16 | int32 | int64 | float64
}

func main() {
	myList := list[int]{}
	myList.add(12)
	myList.add(24)
	myList.add(36)
	cur := myList.start
	for {
		fmt.Println("*", cur)
		if cur == nil {
			break
		}
		cur = cur.next
	}
	fmt.Println()

	treeLastStr := TreeLast[string]{"aa", "bb"}
	treeLastStr.replaceLast("cc")
	fmt.Println(treeLastStr)
	treeLastInt := TreeLast[int]{1, 2, 3}
	treeLastInt.replaceLast(4)
	fmt.Println(treeLastInt)
	fmt.Println()

	name := "John"
	PrintSlice([]string{name})
	fmt.Println()

	fmt.Println(Same(1, 3))
	fmt.Println(Same("string", "string"))
	fmt.Println()

	fmt.Println(Add(10, 3))
	fmt.Println(Add(30.3, 3))
	fmt.Println()

	fmt.Println(AddElements([]AnotherInt{1, 2, 3, 4, 5}))
	fmt.Println()

	f1([]int{1, 2, 3})
}

func AddElements[T AllInts](s []T) T {
	sum := T(0)
	for _, v := range s {
		sum += v
	}
	return sum
}

func Add[T Numeric](a, b T) T {
	return a + b
}

func PrintSlice[T any](s []T) {
	for i, v := range s {
		fmt.Println(i, v)
	}
}

func Same[T comparable](a, b T) bool {
	if a == b {
		return true
	}
	return false
}

func f1[S ~[]E, E any](s S) {
}
