// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package list implements a doubly linked list, forked from container/list.
//
// To iterate over a list (where l is a *List):
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}
//
package list

// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value interface{}
}

// Next returns the next list element or nil.
func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}

// Init initializes or clears list l.
func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func New() *List { return new(List).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List) Len() int { return l.len }

// Front returns the first element of list l or nil.
func (l *List) Front() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil.
func (l *List) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero List value.
func (l *List) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insertAfter inserts e after at, increments l.len, and returns e.
//  ensure that e is not an element of any list and at is an element of list l.
func (l *List) insertAfter(e, at *Element) *Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	return e
}

// insertBefore inserts e before at, increments l.len, and returns e.
//  ensure that e is not an element of any list and at is an element of list l.
func (l *List) insertBefore(e, at *Element) *Element {
	p := at.prev
	at.prev = e
	e.next = at
	e.prev = p
	p.next = e
	e.list = l
	l.len++
	return e
}

// remove removes e from its list, decrements l.len, and returns e.
//  ensure that e is an element of list l.
func (l *List) remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// Remove removes e from l if e is an element of list l and returns e.
//  CAUTION: if e is not an element of list l, will crash
func (l *List) Remove(e *Element) *Element {
	if e.list != l {
		panic("list.Remove: Element e is not an element of list l")
	}
	// if e.list == l, l must have been initialized when e was inserted
	// in l or l == nil (e is a zero Element) and l.remove will crash
	return l.remove(e)
}

// PushFront inserts a new element e at the front of list l and returns e.
//  CAUTION: if e is an element of any list, will crash.
//  more information see comment in List.InsertAfter
func (l *List) PushFront(e *Element) *Element {
	if e.list != nil {
		panic("list.PushFront: Element e can not be an element of any list")
	}
	l.lazyInit()
	return l.insertAfter(e, &l.root)
}

// PushBack inserts a new element e at the back of list l and returns e.
//  CAUTION: if e is an element of any list, will crash.
//  more information see comment in List.InsertAfter
func (l *List) PushBack(e *Element) *Element {
	if e.list != nil {
		panic("list.PushBack: Element e can not be an element of any list")
	}
	l.lazyInit()
	return l.insertBefore(e, &l.root)
}

// InsertBefore inserts a new element e before mark and returns e.
//  CAUTION: if e is an element of any list or mark is not an element of list l, will crash.
//  more information see comment in List.InsertAfter
func (l *List) InsertBefore(e, mark *Element) *Element {
	if e.list != nil {
		panic("list.InsertBefore: Element e can not be an element of any list")
	}
	if mark.list != l {
		panic("list.InsertBefore: Element mark is not an element of list l")
	}
	return l.insertBefore(e, mark)
}

// InsertAfter inserts a new element e after mark and returns e.
//  CAUTION: if e is an element of any list or mark is not an element of list l, will crash.
//
//  To get an element that is not an element of any list, you can do like this:
//  e := new(list.Element)
//  e := list.Remove(*Element)
//  e := list.RemoveBack()
//  e := list.RemoveFront()
func (l *List) InsertAfter(e, mark *Element) *Element {
	if e.list != nil {
		panic("list.InsertAfter: Element e can not be an element of any list")
	}
	if mark.list != l {
		panic("list.InsertAfter: Element mark is not an element of list l")
	}
	return l.insertAfter(e, mark)
}

// MoveToFront moves element e to the front of list l and returns e.
//  CAUTION: if e is not an element of list l, will crash
func (l *List) MoveToFront(e *Element) *Element {
	if e.list != l {
		panic("list.MoveToFront: Element e is not an element of list l")
	}
	if l.root.next == e {
		return e
	}
	// see comment in List.Remove about initialization of l
	return l.insertAfter(l.remove(e), &l.root)
}

// MoveToBack moves element e to the back of list l and returns e.
//  CAUTION: if e is not an element of list l, will crash
func (l *List) MoveToBack(e *Element) *Element {
	if e.list != l {
		panic("list.MoveToBack: Element e is not an element of list l")
	}
	if l.root.prev == e {
		return e
	}
	// see comment in List.Remove about initialization of l
	return l.insertBefore(l.remove(e), &l.root)
}

// MoveBefore moves element e to its new position before mark and returns e.
//  CAUTION: if e or mark is not an element of l, will crash
func (l *List) MoveBefore(e, mark *Element) *Element {
	if e.list != l {
		panic("list.MoveBefore: Element e is not an element of list l")
	}
	if mark.list != l {
		panic("list.MoveBefore: Element mark is not an element of list l")
	}
	if e == mark || e == mark.prev {
		return e
	}
	return l.insertBefore(l.remove(e), mark)
}

// MoveAfter moves element e to its new position after mark and returns e.
//  CAUTION: if e or mark is not an element of l, will crash
func (l *List) MoveAfter(e, mark *Element) *Element {
	if e.list != l {
		panic("list.MoveAfter: Element e is not an element of list l")
	}
	if mark.list != l {
		panic("list.MoveAfter: Element mark is not an element of list l")
	}
	if e == mark || e == mark.next {
		return e
	}
	return l.insertAfter(l.remove(e), mark)
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same.
func (l *List) PushBackList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertBefore(&Element{Value: e.Value}, &l.root)
	}
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same.
func (l *List) PushFrontList(other *List) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertAfter(&Element{Value: e.Value}, &l.root)
	}
}

// RemoveFront removes the first element of list l,
// and returns the element or nil if l is empty.
func (l *List) RemoveFront() *Element {
	if l.len == 0 {
		return nil
	}
	return l.remove(l.root.next)
}

// RemoveBack removes the last element of list l,
// and returns the element or nil if l is empty.
func (l *List) RemoveBack() *Element {
	if l.len == 0 {
		return nil
	}
	return l.remove(l.root.prev)
}
