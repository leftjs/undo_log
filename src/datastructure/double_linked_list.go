package datastructure

type Object interface{}

type Node struct {
	Prev *Node  // 前置节点
	Next *Node  // 后置节点
	Data Object // 当前节点数据
}

func (n *Node) HasNext() bool {
	if n.Next != nil {
		return true
	}
	return false
}

func (n *Node) HasPrev() bool {
	if n.Prev != nil {
		return true
	}

	return false
}

type DoubleLinkedList interface {
	Append(obj Object)              // 追加节点
	InsertPrev(n *Node, obj Object) // 在指定节点之前插入
	InsertNext(n *Node, obj Object) // 在指定节点之后插入
	Remove(n *Node)                 // 移除指定节点
	Search(obj Object) *Node        // 查找指定数据所在节点

	isHead(n *Node) bool // 是否在头部节点
	isTail(n *Node) bool // 是否是尾部节点

}

type LinkedList struct {
	Head *Node // 头部节点
	Tail *Node // 尾部节点
	Size int   // 链表长度
}

func NewLinkedList() *LinkedList {
	return &LinkedList{nil, nil, 0}
}

func (l *LinkedList) Append(obj Object) {

	if l.Size == 0 {
		n := &Node{nil, nil, obj}
		l.Head = n
		l.Tail = n

	} else {
		n := &Node{l.Tail, nil, obj}
		l.Tail.Next = n
		l.Tail = n
	}

	l.Size += 1
}

func (l *LinkedList) InsertPrev(n *Node, obj Object) {
	if n == nil {
		return
	}

	if n.Prev != nil {
		prev := n.Prev
		newN := &Node{prev, n, obj}
		prev.Next = newN
		n.Prev = newN

	} else {
		newN := &Node{nil, n, obj}
		n.Prev = newN

		// update head
		l.Head = newN
	}

	l.Size += 1
}

func (l *LinkedList) InsertNext(n *Node, obj Object) {
	if n == nil {
		return
	}

	if n.Next != nil {
		next := n.Next
		newN := &Node{n, next, obj}
		next.Prev = newN
		n.Next = newN
	} else {
		// append
		newN := &Node{n, nil, obj}
		n.Next = newN

		// update tail
		l.Tail = newN
	}

	l.Size += 1
}

func (l *LinkedList) Remove(n *Node) {
	if n == nil || l.Size == 0 {
		return
	}

	if l.isHead(n) && l.isTail(n) {
		l.Head = nil
		l.Tail = nil
	} else if l.isHead(n) {
		l.Head = n.Next
		n.Next.Prev = nil
	} else if l.isTail(n) {
		l.Tail = n.Prev
		n.Prev.Next = nil
	} else {
		prev := n.Prev
		next := n.Next
		prev.Next = next
		next.Prev = prev
	}

	l.Size -= 1
}

func (l *LinkedList) Search(obj Object) *Node {
	var cur = l.Head
	for {
		if cur.Data == obj {
			return cur
		}
		if cur.HasNext() {
			cur = cur.Next
		} else {
			break
		}
	}
	return nil
}

func (l *LinkedList) isHead(n *Node) bool {
	return n == l.Head
}

func (l *LinkedList) isTail(n *Node) bool {
	return n == l.Tail
}
