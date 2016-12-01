package article

import (
  "fmt"
)

// TODO Need some way to push/pull articles.



type Article struct {
  Title string
  Hist Treedoc  // The Treedoc CRDT
  Log []OpLog  // store the local insert/delete operations to be replayed on the host
}


// Create a new empty article with specified title
func NewArticle(title string) {
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
// updates the artivle param to point to the article with the soecified title
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
  a,err := OpenArticle("../articles/", article.Title)
  if err != nil {
    return err
  }

  // The replay of the client's log is executed on the server
  return a.Replay(remoteArticle)
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
  for _,paragraph := range atom {
    fmt.Printf("%s\n", paragraph)
  }
  fmt.Printf("\n\n")
}

func (a *Article) Insert(pos int, atom Atom, site Disambiguator) error {
  n,err := a.Hist.Insert(pos, atom, site)
  if err != nil {
    return err
  }

  // insert into log
  a.Buffer = append(a.Log, Operation{"insert", n})
}

func (a *Article) Delete(pos int, site Disambiguator) error {
  n,err := a.Hist.Delete(pos, site)
  if err != nil {
    return err
  }

  a.Buffer = append(a.Log, Operation{"delete", n})
}

// Replay all commands in the buffer
func (a *Article) Replay(remoteLog OpLog) error {
  var err error
  for _,op := range remoteLog {
    switch op {
    case "insert":
      err = a.Hist.insertNode(op.Node)
    case "delete":
      err = a.Hist.deleteNode(op.Node)
    default:
      return fmt.Errorf("Article::FlushBuffer() - Invalid Command.")
    }
  }
  return err
}
