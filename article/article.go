package article

import (
  "fmt"
  "strings"
  "io/ioutil"
  "encoding/json"
)

type Article struct {
  Title string
  Hist Treedoc  // The Treedoc CRDT
  Log OpLog  // store the local insert/delete operations to be replayed on the host
}

// A list of insert/delete commands executed locally
type OpLog []Operation

type Operation struct {
  Command string  // "insert" or "delete"
  Path Path
  Value Atom
  Site Disambiguator
}

// Create a new empty article with specified title
func NewArticle(title string) *Article {
  a := Article{}
  a.Title = title

  return &a
}

func OpenArticle(path string, title string) (*Article, error) {
  var a Article
  dat,err := ioutil.ReadFile(path + title + ".json")
  err = json.Unmarshal(dat, &a)

  return &a, err
}


// NOTE @Rain This will be an RPC
// updates the article param to point to the article with the soecified title
func PullArticle(title string, article *Article) error {
  // This executes on the server using the title requested by the client
  a,err := OpenArticle("../articles/", title)
  if err != nil {
    return err
  }

  article = a

  return nil
}


// NOTE @Rain This will be an RPC
// NOTE Don't think we need replayCount but RPC methods require two params
// NOTE We don't need to send the entire article -- just the title and the log.
//      It's easier to just send the entire article though
// This is supposed to execute on the server. The client C sends it's copy of the
// article to the server S. S replays C's log on the (shared) version of the article.
func Push(remoteArticle Article, replayCount *int) error {
  // This should be exectued on the server
  a,err := OpenArticle("../articles/", remoteArticle.Title)
  if err != nil {
    return err
  }

  // The replay of the client's log is executed on the server
  return a.Replay(remoteArticle.Log)
}


// Write the article to the specified path as a JSON file. The file will be named
// with the title of the article.
func (a *Article) Save(path string) error {
  dat,err := json.Marshal(*a)
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(path + a.Title + ".json", dat, 0644)
  return err
}


func (a *Article) Print() {
  atoms := a.Hist.Contents()
  fmt.Printf("\n%s\n", a.Title)
  fmt.Printf("%s\n", strings.Repeat("-", len(a.Title)))
  for _,paragraph := range atoms {
    fmt.Printf("%s\n", paragraph)
  }
  fmt.Printf("\n\n")
}

func (a *Article) Insert(pos int, atom string, site string) error {
  p,err := a.Hist.Insert(pos, Atom(atom), Disambiguator(site))
  if err != nil {
    return err
  }

  // insert into log
  a.Log = append(a.Log, Operation{"insert", p, Atom(atom), Disambiguator(site)})

  return nil
}

func (a *Article) Delete(pos int, site string) error {
  p,err := a.Hist.Delete(pos, Disambiguator(site))
  if err != nil {
    return err
  }

  var empty Atom
  a.Log = append(a.Log, Operation{"delete", p, empty, Disambiguator(site)})

  return nil
}

// Replay all commands in the buffer
func (a *Article) Replay(remoteLog OpLog) error {
  var err error
  for _,op := range remoteLog {
    switch op.Command {
    case "insert":
      err = a.Hist.insertNode(op.Path, &Node{op.Value, op.Site, false, nil, nil})
    case "delete":
      err = a.Hist.deleteNode(op.Path)
    default:
      return fmt.Errorf("Article::FlushBuffer() - Invalid Command.")
    }
  }
  return err
}
