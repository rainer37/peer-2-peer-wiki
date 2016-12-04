package article

// import (
// //  "testing"
// //  "fmt"
//   "encoding/gob"
//   "bytes"
// )

var sharedArticle *Article
var c1Log OpLog
var c2Log OpLog

func createSharedArticle() {
  sharedArticle = NewArticle("chars")
  sharedArticle.Insert(1, "B", "dS")
  sharedArticle.Insert(2, "D", "dS")
  sharedArticle.Log = OpLog{}  // clear the log

  sharedArticle.Save("../articles/")
}

func ExampleClient1() {
  createSharedArticle()

  localArticle,_ := OpenArticle("../articles/", "chars")

  localArticle.Insert(1, "A", "dC1")
  localArticle.Insert(3, "C", "dC1")

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
  localArticle,_ := OpenArticle("../articles/", "chars")

  localArticle.Insert(1, "X", "dC2")
  localArticle.Delete(3, "dC2")

  c2Log = localArticle.Log

  localArticle.Print()
  // Output:
  //
  // chars
  // -----
  // X
  // B
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
}
