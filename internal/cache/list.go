package lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	size int
	head *ListItem
	tail *ListItem
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	listItem := &ListItem{
		Value: v,
		Next:  l.head,
	}

	l.head = listItem

	if l.size == 0 {
		l.tail = l.head
	} else {
		l.head.Next.Prev = l.head
	}

	l.size++
	return l.head
}

// PushBack добавляем элемент в конец.
func (l *list) PushBack(v interface{}) *ListItem {
	listItem := &ListItem{
		Value: v,
		Prev:  l.tail,
	}

	if l.size == 0 {
		l.head = listItem
	} else {
		l.tail.Next = listItem
	}

	l.tail = listItem
	l.size++
	return l.tail
}

func (l *list) Remove(i *ListItem) {
	switch i {
	case l.head:
		if i.Next != nil {
			i.Next.Prev = nil
		}

		l.head = i.Next
	case l.tail:
		if i.Prev != nil {
			i.Prev.Next = nil
		}
		l.tail = i.Prev
	default:
		prevItem := i.Prev
		nextItem := i.Next

		prevItem.Next = nextItem
		nextItem.Prev = prevItem
	}

	l.size--
	if l.size == 0 {
		l.tail = nil
	}
}

// MoveToFront переместит элемент в начало.
func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}

	prevItem := i.Prev
	nextItem := i.Next

	if prevItem != nil {
		prevItem.Next = nextItem
	}

	if nextItem != nil {
		nextItem.Prev = prevItem
	}

	i.Prev = nil
	i.Next = l.head
	l.head = i
}

func NewList() List {
	return new(list)
}
