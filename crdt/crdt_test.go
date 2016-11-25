package crdt

import (
  "testing"
  "fmt"
  //"reflect"
)


func TestTreedocInsert(t *testing.T) {
  // tree := &Treedoc{&[]node{node{"c",posId{path{}, "dC"},false}}, nil, nil}
  // tree.Insert("b", 1, "dB")
  // tree.Insert("a", 1, "dA")
  // tree.Insert("d", 4, "dD")
  // tree.Insert("e", 5, "dE")
  // tree.Insert("f", 6, "dF")
  // fmt.Println(tree.Contents())

  mNode := node{"c",posId{path{}, "dC"},false}
  t2 := &Treedoc{}
  t2.miniNodes = append(t2.miniNodes, &mNode)
  //t2 := &Treedoc{&[]node{node{"c",posId{path{}, "dC"},false}}, nil, nil}
  fmt.Println(t2.Contents())
  t2.Insert("e", 10, "dE")
  t2.insertNode(&node{"x",posId{path{true}, "dC"},false})
  fmt.Println(t2.Contents())
  t2.Insert("a", 1, "dA")
  fmt.Println(t2.Contents())
  t2.Insert("d", 3, "dD")
  fmt.Println(t2.traverse())
  fmt.Println(t2.Contents())
  t2.Insert("f", 4, "dF")
  fmt.Println(t2.Contents())
  t2.Insert("b", 2, "dB")
  fmt.Println(t2.Contents())

  t2.Delete(1, "dC")
  fmt.Println(t2.Contents())
}








// Test the dynamic vector clock
func TestEquals(t *testing.T) {
  v := DynamicVectorClock{}
  w := DynamicVectorClock{}

  // check that nil maps are trivially equal
  if !v.Equals(w) {
    t.Error("Expected v to equal w.")
  }

  v.Append("a")
  if v.Equals(w) {
    t.Error("Expected v to not equal w after appending key 'a' to v.")
  }

  w.Append("a")
  if !v.Equals(w) {
    t.Error("Expected v to eqaul w after appending key 'a' to w.")
  }

  v.Increment("a")
  if v.Equals(w) {
    t.Error("Expected v to not equal w after incrementing key 'a' of v.")
  }

  w.Increment("a")
  if !v.Equals(w) {
    t.Error("Expected v to equal w after incrementing key 'a' of w.")
  }

  v.Append("b")
  if v.Equals(w) {
    t.Error("Expected v to not equal w after appending key 'b' to v.")
  }

  w.Append("b")
  if !v.Equals(w) {
    t.Error("Expected v to eqaul w after appending key 'b' to w.")
  }

  v.Increment("b")
  if v.Equals(w) {
    t.Error("Expected v to not equal w after incrementing key 'b' of v.")
  }

  w.Increment("b")
  if !v.Equals(w) {
    t.Error("Expected v to equal w after incrementing key 'b' of w.")
  }
}


func TestBefore(t *testing.T) {
  v := DynamicVectorClock{}
  w := DynamicVectorClock{}

  // v=[{}]; w=[{}]
  if v.Before(w) {
    t.Errorf("Not expecting v < w. Both v and w should be nil.")
  }

  // v=[{a,0}]; w=[{}]
  // w < v
  v.Append("a")
  if v.Before(w) {
    t.Errorf("Not expecting v < w. v should dominate w after appending key 'a'.")
  } else if !w.Before(v) {
    t.Errorf("Expected w < v after appending key 'a' to v.")
  }

  // v=[{a,1}]; w=[{}]
  // w < v
  v.Increment("a")
  if v.Before(w) {
    t.Errorf("Not expecting v < w. v should dominate w after incrementing key 'a'.")
  } else if !w.Before(v) {
    t.Errorf("Expected w < v after incrementing key 'a' to v.")
  }

  // v=[{a,1}]; w=[{b,0}]
  // v <> w
  w.Append("b")
  if v.Before(w) {
    t.Errorf("Not expecting v < w after appending key 'b' to w.")
  } else if w.Before(v) {
    t.Errorf("Not expecting w < v after appending key 'b' to w.")
  }

  // v=[{a,1},{b,0}]; w=[{b,0}]
  // w < v
  v.Append("b")
  if v.Before(w) {
    t.Errorf("Not expecting v < w after appending key 'b' to v.")
  } else if !w.Before(v) {
    t.Errorf("Expecting w < v after appending key 'b' to v.")
  }

  // v=[{a,1},{b,0}]; w=[{b,1}]
  // w <> v
  w.Increment("b")
  if v.Before(w) {
    t.Errorf("Not expecting v < w after incrementing key 'b' of w.")
  } else if w.Before(v) {
    t.Errorf("Not Expecting w < v after incrementing key 'b'of w.")
  }

  // v=[{a,1},{b,0}]; w=[{a,0},{b,1}]
  // w <> v
  w.Append("a")
  if v.Before(w) {
    t.Errorf("Not expecting v < w after appending key 'a' to w.")
  } else if w.Before(v) {
    t.Errorf("Not Expecting w < v after appending key 'a' to w.")
  }

  // v=[{a,1},{b,0}]; w=[{a,1},{b,1}]
  // v < w
  w.Increment("a")
  if !v.Before(w) {
    t.Errorf("Expected v < w after incrementing key 'a' of w.")
  } else if w.Before(v) {
    t.Errorf("Not expecting w < v after incrementing key 'a' of w.")
  }
}
