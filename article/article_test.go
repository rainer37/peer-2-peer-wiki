package article

import (
//  "testing"
//  "fmt"
  "encoding/gob"
  "bytes"
)

func Clone(a,b interface{}) {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	enc.Encode(a)
	dec.Decode(b)
}


var sharedArticle *Article
var c1Log OpLog
var c2Log OpLog

func createSharedArticle() {
  sharedArticle = NewArticle("chars")
  sharedArticle.Insert(1, "B", "dS")
  sharedArticle.Insert(2, "D", "dS")
  sharedArticle.Log = OpLog{}  // clear the log
}

func ExampleClient1() {
  createSharedArticle()
  localArticle := NewArticle("chars")
  Clone(localArticle,&sharedArticle)

  localArticle.Insert(1, "A", "dC1")
//  localArticle.Insert(3, "C", "dC1")

  c1Log = localArticle.Log

  localArticle.Print()
  // Output:
  //
  // chars
  // -----
  // A
  // B
  // C
  // D
}

func ExampleClient2() {
  localArticle := *sharedArticle

  localArticle.Insert(1, "X", "dC2")

  c2Log = localArticle.Log

  localArticle.Print()
  // Output:
  //
  // chars
  // -----
  // X
  // B
  // D
}

func ExampleMerge() {
  sharedArticle.Replay(c1Log)
  sharedArticle.Replay(c2Log)

  sharedArticle.Print()
  // Output:
  //
  // chars
  // -----
  // A
  // X
  // B
  // C
  // D
}
