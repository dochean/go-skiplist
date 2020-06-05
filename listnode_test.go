package skiplist

import (
	"testing"
)

var (
	headlist *HeadNode
)

func init() {
	headlist = newHeadNode()

	for i := 0; i <= 10000000; i++ {
		headlist.Insert(newnum(i))
	}
}

type newnum int

func (n newnum) Compare(nn interface{}) int {
	num, ok := nn.(newnum)
	// fmt.Printf("n, newnum: %#v, %#v, %+v\n", n, nn, ok)
	if !ok {
		return -1
	}
	if n > num {
		return 1
	} else if n == num {
		return 0
	}
	return -1
}

func check(flag bool, t *testing.T) {
	if !flag {
		t.Fatalf("list is not sorted\n")
	}
}

func TestListCRUD(t *testing.T) {
	nl := []newnum{
		131, 21, 111, 1221, 231,
	}
	h := newHeadNode()
	for _, v := range nl {
		h.Insert(v)
	}
	check(h.IsSorted(), t)
	h.Insert(newnum(60))
	check(h.IsSorted(), t)

	h.Delete(newnum(21))
	check(h.IsSorted(), t)
	h.Delete(newnum(231))
	check(h.IsSorted(), t)

	v1 := h.Get(newnum(131))
	v2 := h.Get(newnum(1221))
	if v1 == nil || v1.(newnum) != 131 {
		t.Fatal(`wrong value (expected "131")`, v1)
	}
	if v2 == nil || v2.(newnum) != 1221 {
		t.Fatal(`wrong value (expected "1221")`, v2)
	}
}

func BenchmarkIncSet(b *testing.B) {
	b.ReportAllocs()
	list := newHeadNode()

	for i := 0; i < b.N; i++ {
		list.Insert(newnum(i))
	}

	b.SetBytes(int64(b.N))
}

func BenchmarkIncGet(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res := headlist.Search(newnum(i))
		if res == false {
			b.Fatal("failed to Get an element that should exist")
		}
	}

	b.SetBytes(int64(b.N))
}
