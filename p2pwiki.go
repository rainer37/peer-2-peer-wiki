package p2pwiki

import (
  "os"
  "fmt"
  "strconv"
  "github.com/nickbradley/p2pwiki/article"
)

func main() {
  args := os.Args[1:]

  switch args[0] {
  case "article":
    switch args[1] {
    case "pull":
      title := args[2]

      contents,err := chord.Lookup(title)
      if err != nil {
        // print warning re creating new article
      }
      article := article.NewLocal(title, contents)
      article.Save()

    case "insert":
      title := args[2]
      article,err := article.OpenLocal(title)
      if err != nil {
        fmt.Fatal("You must first pull article.")
      }

      pos,err := strconv.Atoi(args[3])
      if err != nil {
        fmt.Fatal("Invalid position parameter.")
      }
      text := args[4]
      err = article.Insert(pos, text)
      if err != nil {
        fmt.Fatal("Failed to insert paragraph.")
      }
      article.Save()

    case "delete":
      title := args[2]
      article,err := article.OpenLocal(title)
      if err != nil {
        fmt.Fatal("You must first pull article.")
      }

      pos,err := strconv.Atoi(args[3])
      if err != nil {
        fmt.Fatal("Invalid position parameter.")
      }
      text := args[4]
      err = article.Delete(pos)
      if err != nil {
        fmt.Fatal("Failed to delete paragraph.")
      }
      article.Save()
      
    case "push":
      title := args[2]
      article,err := article.OpenLocal(title)
      if err != nil {
        fmt.Fatal("You must first pull article.")
      }
      chord.Send(article.log)
    default:
      fmt.Fatal("Invalid article command.")
    }
  default:
    fmt.Fatal("Invalid command.")
  }
}
