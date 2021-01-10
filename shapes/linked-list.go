package shapes

import (
	"fmt"
)

type LinkedList struct {
	length int64
	head   *Node
	tail   *Node
}

type Node struct {
	val  float64
	next *Node
}

func NewLinkedList() (l *LinkedList) {
	return &LinkedList{}
}

func (l *LinkedList) Length() int64 {
	return l.length
}
func (l *LinkedList) AddFirst(val float64) {
	node := Node{val: val, next: l.head}

	if l.length == 0 {
		l.tail = &node
	}

	l.length = l.length + 1
	l.head = &node
}

func (l *LinkedList) AddLast(val float64) {
	node := Node{val: val, next: nil}

	if l.length == 0 {
		l.length = l.length + 1
		l.head = &node
		l.tail = &node
		return
	}

	l.length = l.length + 1

	l.tail.next = &node
	l.tail = &node
}

func (l *LinkedList) IsEmpty() bool {
	return l.head == nil
}

func (l *LinkedList) RemoveFirst() (result float64) {

	l.length = l.length - 1
	result = l.head.val
	l.head = l.head.next
	return
}

func (l *LinkedList) String() (s string) {

	values := []string{}

	node := l.head

	for node != nil {
		values = append(values, fmt.Sprintf("%v", node.val))
		node = node.next
	}

	s = fmt.Sprintf("%v", values)
	return
}
