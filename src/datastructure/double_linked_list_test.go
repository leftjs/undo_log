package datastructure

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLinkedList(t *testing.T) {
	l := NewLinkedList()
	assert.Nil(t, l.Head)
	assert.Nil(t, l.Tail)
	assert.Equal(t, 0, l.Size)
}

func TestLinkedList_Append(t *testing.T) {
	l := NewLinkedList()

	// haha -> hehe -> xixi

	l.Append("haha")
	assert.Equal(t, "haha", l.Head.Data)
	assert.Equal(t, "haha", l.Tail.Data)
	assert.Equal(t, 1, l.Size)

	l.Append("hehe")
	assert.Equal(t, "haha", l.Head.Data)
	assert.Equal(t, "hehe", l.Tail.Data)
	assert.Equal(t, 2, l.Size)

	l.Append("xixi")
	assert.Equal(t, "haha", l.Head.Data)
	assert.Equal(t, "hehe", l.Head.Next.Data)
	assert.Equal(t, l.Head.Next.Data, l.Tail.Prev.Data)
	assert.Equal(t, "xixi", l.Tail.Data)
}

func TestLinkedList_Search(t *testing.T) {
	l := NewLinkedList()

	// haha <-> hehe <-> xixi
	l.Append("haha")
	l.Append("hehe")
	l.Append("xixi")

	assert.Equal(t, "hehe", l.Search("hehe").Data)
}

func TestLinkedList_InsertPrev(t *testing.T) {
	l := NewLinkedList()

	// haha <-> hehe <-> xixi
	l.Append("haha")
	l.Append("hehe")
	l.Append("xixi")

	// insert to head

	// head <-> haha <-> hehe <-> xixi
	l.InsertPrev(l.Search("haha"), "head")

	assert.Equal(t, "head", l.Head.Data)
	assert.Equal(t, "haha", l.Head.Next.Data)
	assert.Equal(t, "head", l.Search("haha").Prev.Data)
	assert.Nil(t, l.Search("head").Prev)
	assert.Equal(t, 4, l.Size)

	// head <-> haha <-> hh <-> hehe <-> xixi
	l.InsertPrev(l.Search("hehe"), "hh")
	assert.Equal(t, "hh", l.Search("haha").Next.Data)
	assert.Equal(t, "hehe", l.Search("hh").Next.Data)
}

func TestLinkedList_InsertNext(t *testing.T) {
	l := NewLinkedList()

	// haha <-> hehe <-> xixi
	l.Append("haha")
	l.Append("hehe")
	l.Append("xixi")

	// insert to tail
	// haha <-> hehe <-> xixi <-> tail
	l.InsertNext(l.Tail, "tail")
	assert.Equal(t, l.Search("tail"), l.Tail)
	assert.Equal(t, 4, l.Size)

	// haha <-> hh <-> hehe <-> xixi <-> tail
	l.InsertNext(l.Head, "hh")
	assert.Equal(t, "hh", l.Head.Next.Data)
	assert.Equal(t, "hehe", l.Head.Next.Next.Data)
	assert.Equal(t, "hh", l.Head.Next.Next.Prev.Data)
}

func TestLinkedList_Remove(t *testing.T) {
	l := NewLinkedList()

	// head <-> haha <-> hehe <-> xixi <-> tail
	l.Append("head")
	l.Append("haha")
	l.Append("hehe")
	l.Append("xixi")
	l.Append("tail")

	// head <-> haha <-> hehe <-> xixi <-> tail
	//          /
	//         /
	//        /
	// head <-> hehe <-> xixi <-> tail
	l.Remove(l.Search("haha"))
	assert.Equal(t, "hehe", l.Head.Next.Data)
	assert.Equal(t, "head", l.Head.Next.Prev.Data)

	// hehe <-> xixi <-> tail
	l.Remove(l.Head)
	assert.Equal(t, "hehe", l.Head.Data)
	assert.Equal(t, "xixi", l.Head.Next.Data)
	assert.Nil(t, l.Head.Prev)

	// hehe <-> xixi
	l.Remove(l.Tail)
	assert.Equal(t, "xixi", l.Tail.Data)
	assert.Equal(t, "hehe", l.Head.Data)
	assert.Equal(t, l.Head, l.Tail.Prev)

}
