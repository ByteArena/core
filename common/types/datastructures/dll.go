package datastructures

type DLL struct {
	Head *DLLNode
	Tail *DLLNode
}

type DLLNode struct {
	Next *DLLNode
	Prev *DLLNode
	Val  interface{}
}

func (list *DLL) Append(val interface{}) *DLL {
	node := &DLLNode{}
	node.Val = val
	if list.Tail != nil {
		oldtail := list.Tail
		oldtail.Next = node

		node.Prev = oldtail
	}

	list.Tail = node

	if list.Head == nil {
		list.Head = node
	}

	return list
}

func (list *DLL) Empty() bool {
	return list.Head == nil
}

func (list *DLL) Clear() *DLL {
	list.Head = nil
	list.Tail = nil

	return list
}

func (list *DLL) RemoveVal(val interface{}) *DLL {

	node := list.Head
	found := false

	for node != nil {
		if node.Val == val {
			found = true
			break
		}

		node = node.Next
	}

	if !found {
		return list
	}

	if node.Prev != nil {
		node.Prev.Next = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	}

	if list.Head == node {
		list.Head = node.Next
	}

	if list.Tail == node {
		list.Tail = node.Prev
	}

	return list
}

func (list *DLL) InsertBefore(node *DLLNode, val interface{}) *DLL {

	if node == nil {
		return list
	}

	i := &DLLNode{Val: val}
	i.Next = node
	i.Prev = node.Prev

	//oldprev := node.prev

	if node.Prev != nil {
		node.Prev.Next = i
	}

	node.Prev = i

	if node == list.Head {
		list.Head = i
	}

	return list
}
