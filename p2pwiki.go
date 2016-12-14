package main

import (
  "os"
  "fmt"
  "log"
  "strconv"
  "strings"
  "github.com/nickbradley/p2pwiki/chord"
  "github.com/nickbradley/p2pwiki/article"
  "bitbucket.org/bestchai/dinv/dinvRT"
)

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func unpack(log []byte) {
  var dummy []byte
  dinvRT.Unpack(log, &dummy)
}

func main() {
  srvAddr := os.Args[1]
  mainArg := os.Args[2]
  subArg := os.Args[3:]

  switch mainArg {
  case "server":
    switch subArg[0] {
    case "create":  // p2pwiki 127.0.0.1:2222 server start create
      chord.CreateRing(srvAddr)
    case "join":  // p2pwiki 127.0.0.1:2222 server start join 127.0.0.1:3333
      new_node := chord.NewNode(srvAddr)
      new_node.Join(subArg[1])
    default: log.Fatal("Invalid server start command.")
    }
  case "article":
    subdir := strings.Replace(srvAddr, ":","_",1)
    srvAddr := subArg[0]
    subCmd := subArg[1]
    subArg := subArg[2:]

    disab := subdir

    dinvRT.Initalize("Client")

    cacheDir := "./articles/cache/"+subdir+"/"
    localDir := "./articles/local/"+subdir+"/" //local folder

    if exist,_ := exists(cacheDir); !exist {
      os.MkdirAll(cacheDir, 0777)
    }

    if exist,_ := exists(localDir); !exist {
      os.MkdirAll(localDir, 0777)
    }

    switch subCmd {
    case "pull":  // p2pwiki 127.0.0.1:2222 article pull "<title>"
      title := subArg[0]

      //var a article.Article
      var owner chord.StrLog //string = ""
      err := chord.RPCall(srvAddr, &chord.StrLog{title,dinvRT.Pack(nil)}, &owner, "Node.Find_article") // find the owner of the article

      unpack(owner.Log)

      if err != nil {
        fmt.Println("@@",err)
        break
      }

      println(owner.Str, "should have this article", chord.Hash(owner.Str))

      var art chord.ArtLog

      err = chord.RPCall(owner.Str, &chord.StrLog{title,dinvRT.Pack(nil)}, &art, "Node.Pull")

      if err != nil {
        fmt.Println("#",err)
        break
      }

      unpack(art.Log)

      println("Article "+art.Art.Title+" has been pulled successfully from "+srvAddr)

      art.Art.Save(cacheDir)

    case "insert":  // p2pwiki 127.0.0.1:2222 article insert "<title>" <pos> "<text>"
      title := subArg[0]
      article,err := article.OpenArticle(cacheDir, title)
      if err != nil {
        log.Fatal("You must first pull article.")
      }

      pos,err := strconv.Atoi(subArg[1])
      if err != nil {
        log.Fatal("Invalid position parameter.")
      }

      err = article.Insert(pos, subArg[2], disab)
      if err != nil {
        log.Fatal("Failed to insert paragraph.")
      }
      article.Save(localDir)
      fmt.Println("Insert into",title,"succeed...")
    case "delete":  // p2pwiki 127.0.0.1:2222 article delete "<title>" <pos>
      title := subArg[0]
      article,err := article.OpenArticle(cacheDir, title)
      if err != nil {
        log.Fatal("You must first pull article.")
      }

      pos,err := strconv.Atoi(subArg[1])
      if err != nil {
        log.Fatal("Invalid position parameter.")
      }

      err = article.Delete(pos, srvAddr)
      if err != nil {
        log.Fatal("Failed to delete paragraph.")
      }
      article.Save(cacheDir)
    case "push":  // p2pwiki 127.0.0.1:2222 article push "<title>"
      title := subArg[0]

      a,err := article.OpenArticle(cacheDir, title)
      if err != nil {
        log.Fatal("You must first pull article.")
      }

      a,err = article.OpenArticle(localDir, title)
      if err != nil {
        log.Fatal("No change has been made locally")
      }

      a.Save(cacheDir)
      //break


      var owner chord.StrLog

      err = chord.RPCall(srvAddr, &chord.StrLog{title,dinvRT.Pack(nil)}, &owner, "Node.Find") // find the owner of the article

      unpack(owner.Log)

      if err != nil {
        fmt.Println("@@",err)
        break
      }

      pair := chord.ArtPair{a.Title, a.Log, dinvRT.Pack(nil)}
      var log []byte

      chord.RPCall(owner.Str, &pair, &log, "Node.PushArticle")

      unpack(log)

      println("Article "+a.Title+" has pushed to "+srvAddr)
    case "view":
      title := subArg[0]
      article1,err := article.OpenArticle(cacheDir, title)

      println("Opening article:", cacheDir+title,".json")

      if err != nil {
        fmt.Println(err)
        log.Fatal("You must first pull article.")
      }

      article2,err := article.OpenArticle(localDir, title)

      if err != nil {
        fmt.Println("No Local Change yet...")
        article1.Print()
      } else {
        article2.Print()
      }
    case "discard":
      title := subArg[0]
      _,err := article.OpenArticle(localDir, title)

      println("Opening article:", localDir+title,".json")

      if err != nil {
        fmt.Println(err)
        log.Fatal("No local change made to this article")
      }

      os.Remove(localDir+title+".json")
      fmt.Println("Local change to Article", title, "discard...")
    default:
      log.Fatal("Invalid article command.")
    }
  default: log.Fatal("Invalid command.")
  }

}
