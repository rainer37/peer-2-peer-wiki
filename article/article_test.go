package article

import (
  "testing"
  "fmt"
)


// func TestTreedocInsert(t *testing.T) {
//   // tree := &Treedoc{&[]node{node{"c",posId{path{}, "dC"},false}}, nil, nil}
//   // tree.Insert("b", 1, "dB")
//   // tree.Insert("a", 1, "dA")
//   // tree.Insert("d", 4, "dD")
//   // tree.Insert("e", 5, "dE")
//   // tree.Insert("f", 6, "dF")
//   // fmt.Println(tree.Contents())
//
//   mNode := node{"c",posId{path{}, "dC"},false}
//   t2 := &Treedoc{}
//   t2.miniNodes = append(t2.miniNodes, &mNode)
//   //t2 := &Treedoc{&[]node{node{"c",posId{path{}, "dC"},false}}, nil, nil}
//   fmt.Println(t2.Contents())
//   t2.Insert("e", 10, "dE")
//   t2.insertNode(&node{"x",posId{path{true}, "dC"},false})
//   fmt.Println(t2.Contents())
//   t2.Insert("a", 1, "dA")
//   fmt.Println(t2.Contents())
//   t2.Insert("d", 3, "dD")
//   fmt.Println(t2.traverse())
//   fmt.Println(t2.Contents())
//   t2.Insert("f", 4, "dF")
//   fmt.Println(t2.Contents())
//   t2.Insert("b", 2, "dB")
//   fmt.Println(t2.Contents())
//
//   t2.Delete(1, "dC")
//   fmt.Println(t2.Contents())
// }





// var tree *Treedoc
//
// func TestInsert(t *testing.T) {
//   tree = &Treedoc{}
//   mNode := node{"c",posId{path{}, "dC"},false}
//   tree.miniNodes = append(tree.miniNodes, &mNode)
//
//   fmt.Println(tree.Contents())
//   err := tree.Insert("Beer is delicious.", 1, "127.0.0.1:1234")
//   fmt.Println(tree.Contents())
//   err = tree.Insert("There are many types of beer.", 1, "127.0.0.1:1234")
//   if err != nil {
//     fmt.Println("Something terrible has happened!")
//   }
//   fmt.Println(tree.Contents())
// }




var lb *LocalBuffer

func TestNewArticle(t *testing.T) {
  title := "Beer"
  lb = NewLocalBuffer("Beer", []string{}, "127.0.0.1:1234")
  fmt.Printf("Created new article empty article %s.\n", title)
}

func TestArticleInsert(t *testing.T) {
  lb.Insert(1, "Beer is delicious.")
  lb.Insert(2, "There are many types of beer.")

  lb.Print()

  lb.Insert(3, "Beer has been around for a long time.")

  lb.Print()

  lb.Delete(1)

  lb.Print()

  //
  // fmt.Println(a)
  // //a.Save()
  // a.Delete(1)
  // err := a.Save("../articles/local/")
  // if err != nil {
  //   fmt.Println("Failed to write file", err)
  // }
  // fmt.Println(a)
  // //fmt.Println("JSON is", string(js[:]))
  // b,err := OpenLocalBuffer("../articles/local/", "Article 1")
  // if err != nil {
  //   fmt.Println("Failed to read file", err)
  // }
  // fmt.Println("Article from file\n",b)

}


var sb *SharedBuffer

func TestNewSharedBuffer(t *testing.T) {
  sb = NewSharedBuffer(lb.Title)
  fmt.Println(sb.Hist.isEmpty())
}

func TestReplay(t *testing.T) {
  sb.Replay(lb.Log)
  fmt.Println(sb.Contents())
}
