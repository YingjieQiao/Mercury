package main

import (
	"sort"
	"testing"
	"testing/quick"

)

func checkNum(num, expected int, t *testing.T) {
	if num != expected {
		t.Errorf("got %d, expected %d", num, expected)
	}
}

func TestNew(t *testing.T) {
	x := NewConsistent(20)
	if x == nil {
		t.Errorf("expected obj")
	}
	checkNum(x.numOfvNodesPerpNodes, 20, t)
}

func TestAdd(t *testing.T) {
	x := NewConsistent(20)
	x.AddPhysicalNode(123213)
	checkNum(len(x.vNodeID2pNodeID), 20, t)
	checkNum(len(x.sortedHashes), 20, t)
	if sort.IsSorted(x.sortedHashes) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
	x.AddPhysicalNode(5678)
	checkNum(len(x.vNodeID2pNodeID), 40, t)
	checkNum(len(x.sortedHashes), 40, t)
	if sort.IsSorted(x.sortedHashes) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
}

func TestRemove(t *testing.T) {
	x := NewConsistent(20)
	x.AddPhysicalNode(456546)
	x.Remove(456546)
	checkNum(len(x.vNodeID2pNodeID), 0, t)
	checkNum(len(x.sortedHashes), 0, t)
}

func TestRemoveNonExisting(t *testing.T) {
	x := NewConsistent(20)
	x.AddPhysicalNode(456546)
	x.Remove(34657)
	checkNum(len(x.vNodeID2pNodeID), 20, t)
}

func TestGetEmpty(t *testing.T) {
	x := NewConsistent(20)
	_, err := x.Get(76867)
	if err == nil {
		t.Errorf("expected error")
	}
	if err != ErrEmptyCircle {
		t.Errorf("expected empty vNodeID2pNodeID error")
	}
}

func TestGetSingle(t *testing.T) {
	x := NewConsistent(20)
	x.AddPhysicalNode(456546)
	f := func(s uint64) bool {
		y, err := x.Get(s)
		if err != nil {
			t.Logf("error: %d", err)
			return false
		}
		t.Logf("s = %d, y = %d", s, y)
		return y == 456546
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

type gtest struct {
	in  uint64
	out uint64
}

var gmtests = []gtest{
	{4, 456546},
	{5, 1},
	{2, 5},
}



func TestGetMultipleQuick(t *testing.T) {
	x := NewConsistent(20)
	x.AddPhysicalNode(456546)
	x.AddPhysicalNode(5)
	x.AddPhysicalNode(1)
	f := func(s uint64) bool {
		y, err := x.Get(s)
		if err != nil {
			t.Logf("error: %d", err)
			return false
		}
		t.Logf("s = %d, y = %d", s, y)
		return y == 456546 || y == 5 || y == 1
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

var rtestsBefore = []gtest{
	{4, 1},
	{3, 5},
	{2, 5},
}

var rtestsAfter = []gtest{
	{4, 1},
	{3, 1},
	{2, 1},
}

func TestGetMultipleRemove(t *testing.T) {
	x := NewConsistent(20)
	x.numOfvNodesPerpNodes = 1
	x.AddPhysicalNode(456546)
	x.AddPhysicalNode(5)
	x.AddPhysicalNode(1)

	for i, v := range rtestsBefore {
		result, err := x.Get(v.in)
		if err != nil {
			t.Fatal(err)
		}
		if result != v.out {
			t.Errorf("%d. got %d, expected %d before rm", i, result, v.out)
		}
	}
	x.Remove(5)
	for i, v := range rtestsAfter {
		result, err := x.Get(v.in)
		if err != nil {
			t.Fatal(err)
		}
		if result != v.out {
			t.Errorf("%d. got %d, expected %d after rm", i, result, v.out)
		}
	}
}
func TestGetN(t *testing.T) {
	x := NewConsistent(20)
	x.AddPhysicalNode(456546)
	x.AddPhysicalNode(5)
	x.AddPhysicalNode(1)
	members, err := x.GetN(9999999, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(members) != 3 {
		t.Errorf("expected 3 members instead of %d", len(members))
	}
	if members[0] != 1 {
		t.Errorf("wrong members[0]: %d", members[0])
	}
	if members[1] != 5 {
		t.Errorf("wrong members[1]: %d", members[1])
	}
	if members[2] != 456546 {
		t.Errorf("wrong members[2]: %d", members[2])
	}
}


func BenchmarkGet(b *testing.B) {
	x := NewConsistent(20)
	x.AddPhysicalNode(9)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Get(9)
	}
}

func BenchmarkGetLarge(b *testing.B) {
	x := NewConsistent(20)
	for i := 0; i < 10; i++ {
		x.AddPhysicalNode(8 + uint64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.Get(9)
	}
}

func BenchmarkGetN(b *testing.B) {
	x := NewConsistent(20)
	x.AddPhysicalNode(9)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x.GetN(9, 3)
	}
}




// from @edsrzf on github:
func TestAddCollision(t *testing.T) {
	// These two uint64s produce several crc32 collisions after "|i" is
	// appended added by Consistent.eltKey.
	const s1 = 8
	const s2 = 9
	x := NewConsistent(20)
	x.AddPhysicalNode(s1)
	x.AddPhysicalNode(s2)
	elt1, err := x.Get(8)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	y := NewConsistent(20)
	// add elements in opposite order
	y.AddPhysicalNode(s2)
	y.AddPhysicalNode(s1)
	elt2, err := y.Get(s1)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if elt1 != elt2 {
		t.Error(elt1, "and", elt2, "should be equal")
	}
}