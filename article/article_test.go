package article

import (
  "testing"
  "fmt"
)

var a *LocalBuffer

func TestNewArticle(t *testing.T) {
  a = NewLocalBuffer("Article 1", []string{}, "127.0.0.1:1234")
  fmt.Println(a)
}

func TestArticleInsert(t *testing.T) {
  a.Insert(1, "I'm a new paragraph.")
  a.Insert(2, "Oh, I 'member!")
  a.Insert(2, "Do you 'member?")
  fmt.Println(a)
  //a.Save()
  a.Delete(1)
  err := a.Save("../articles/local/")
  if err != nil {
    fmt.Println("Failed to write file", err)
  }
  fmt.Println(a)
  //fmt.Println("JSON is", string(js[:]))
  b,err := OpenLocalBuffer("../articles/local/", "Article 1")
  if err != nil {
    fmt.Println("Failed to read file", err)
  }
  fmt.Println("Article from file\n",b)

}
