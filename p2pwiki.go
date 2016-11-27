package main

import (
  "os"
  "log"
  "strconv"
  "github.com/nickbradley/p2pwiki/chord"
  "github.com/nickbradley/p2pwiki/article"
)

func main() {
  cacheDir := "../articles/cache/"
  srvAddr := os.Args[1]
  appCmd := os.Args[2]
  subCmd := os.Args[3]
  subArg := os.Args[4:]

  // log.Println(os.Args)
  // log.Println(srvAddr)
  // log.Println(appCmd)
  // log.Println(subCmd)
  // log.Println(subArg)

  switch appCmd {
  case "article":
    switch subCmd {
    case "pull":  // p2pwiki 127.0.0.1:2222 article pull "<title>"
      var contents []string
      title := subArg[0]
      hTitle := chord.Hash(title)

      err := chord.RPCall(srvAddr, hTitle, &contents, "Pull")
      if err != nil {
        // print warning re creating new article
      }
      article := article.NewLocalBuffer(title, contents, srvAddr)
      article.Save(cacheDir)

    case "insert":  // p2pwiki 127.0.0.1:2222 article insert "<title>" <pos> "<text>"
      title := subArg[0]
      article,err := article.OpenLocalBuffer(cacheDir, title)
      if err != nil {
        log.Fatal("You must first pull article.")
      }

      pos,err := strconv.Atoi(subArg[1])
      if err != nil {
        log.Fatal("Invalid position parameter.")
      }
      text := subArg[2]
      err = article.Insert(pos, text)
      if err != nil {
        log.Fatal("Failed to insert paragraph.")
      }
      article.Save(cacheDir)

    case "delete":  // p2pwiki 127.0.0.1:2222 article delete "<title>" <pos>
      title := subArg[0]
      article,err := article.OpenLocalBuffer(cacheDir, title)
      if err != nil {
        log.Fatal("You must first pull article.")
      }

      pos,err := strconv.Atoi(subArg[2])
      if err != nil {
        log.Fatal("Invalid position parameter.")
      }

      err = article.Delete(pos)
      if err != nil {
        log.Fatal("Failed to delete paragraph.")
      }
      article.Save(cacheDir)

    case "push":  // p2pwiki 127.0.0.1:2222 article push "<title>"
      title := subArg[0]

      article,err := article.OpenLocalBuffer(cacheDir, title)
      if err != nil {
        log.Fatal("You must first pull article.")
      }

      // Send the operations log to the server. Retry sending unsuccessful operations
      // until all operations have go through.
      expectCount := len(article.Log.Operations)
      replayCount := 0
      for expectCount > replayCount {
        chord.RPCall(srvAddr, article.Log, &replayCount, "Push")
        err = article.Log.Remove(replayCount)
        if err != nil {
          log.Fatal("Unexpected error encountered while sending changes to server.")
        }
      }

      article.Save(cacheDir)
    default:
      log.Fatal("Invalid article command.")
    }

  case "server":
    switch subCmd {
    case "start":
      switch subArg[0] {
      case "create":  // p2pwiki 127.0.0.1:2222 server start create
        chord.CreateRing()
      case "join":  // p2pwiki 127.0.0.1:2222 server start join 127.0.0.1:3333
        new_node := chord.NewNode(srvAddr)
        new_node.Join(subArg[1])
      default: log.Fatal("Invalid server start command.")
      }
    default:
      log.Fatal("Invalid server command.")
    }
  default:
    log.Fatal("Invalid command.")
  }
}